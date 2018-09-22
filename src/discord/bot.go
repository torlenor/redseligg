package discord

import (
	"botinterface"
	"encoding/json"
	"io/ioutil"
	"logging"
	"net/http"
	"os"
	"plugins"
	"strings"

	"events"

	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
)

var (
	discordOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/cb",
		ClientID:     os.Getenv("DISCORD_KEY"),
		ClientSecret: os.Getenv("DISCORD_SECRET"),
		Scopes:       []string{"bot"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discordapp.com/api/oauth2/authorize",
			TokenURL: "https://discordapp.com/api/oauth2/token",
		},
	}

	// Some random string, random for each request
	oauthStateString = "random"

	log = logging.Get("DiscordBot")
)

type guild struct {
	snowflakeID string
	name        string
	memberCount int
	channel     []channel
}

// The Bot struct holds parameters related to the bot
type Bot struct {
	ws *websocket.Conn

	sendMessageChan chan events.SendMessage
	commandChan     chan events.Command

	knownChannels     map[string]channelCreate
	token             string
	ownSnowflakeID    string
	currentSeqNumber  int
	heartBeatSender   *discordHeartBeatSender
	heartBeatStopChan chan struct{}
	seqNumberChan     chan int

	receivers map[plugins.Plugin]chan events.ReceiveMessage

	knownPlugins []plugins.Plugin

	guilds        map[string]guild
	guildNameToID map[string]string
}

// GetReceiveMessageChannel returns the channel which is used to notify
// about received messages from the bot
func (b *Bot) GetReceiveMessageChannel(plugin plugins.Plugin) chan events.ReceiveMessage {
	b.receivers[plugin] = make(chan events.ReceiveMessage)
	return b.receivers[plugin]
}

// GetSendMessageChannel returns the channel which is used to
// send messages using the bot. For DiscordBot these messages
// can be normal channel messages, whispers
func (b Bot) GetSendMessageChannel() chan events.SendMessage {
	return b.sendMessageChan
}

// GetCommandChannel gives a channel to control the bot from
// a plugin
func (b Bot) GetCommandChannel() chan events.Command {
	return b.commandChan
}

func (b Bot) apiCall(path string, method string, body string) (r []byte, e error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, "https://discordapp.com/api"+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bot "+b.token)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}

func (b *Bot) startDiscordBot(doneChannel chan struct{}) {
	defer close(doneChannel)
	for {
		_, message, err := b.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Debugln("Connection closed normally: ", err)
			} else {
				log.Errorln("UNHANDELED ERROR: ", err)
			}
			break
		}

		var data map[string]interface{}

		if err := json.Unmarshal(message, &data); err != nil {
			log.Errorln("UNHANDELED ERROR: ", err)
			continue
		}

		if data["op"].(float64) == 10 { // Hello from Discord Gateway
			log.Debugln("Received HELLO from gateway")
			sendIdent(b.token, b.ws)
			heartbeatInterval := int(data["d"].(map[string]interface{})["heartbeat_interval"].(float64))
			b.heartBeatSender = &discordHeartBeatSender{ws: b.ws}
			go heartBeat(heartbeatInterval, b.heartBeatSender, b.heartBeatStopChan, b.seqNumberChan)
			log.Infoln("DiscordBot is READY")
		} else if data["op"].(float64) == 0 { // Dispatch to event handlers
			switch data["t"] {
			case "MESSAGE_CREATE":
				b.handleMessageCreate(data)
			case "READY":
				b.handleReady(data)
			case "GUILD_CREATE":
				b.handleGuildCreate(data)
			case "PRESENCE_UPDATE":
				b.handlePresenceUpdate(data)
			case "PRESENCE_REPLACE":
				b.handlePresenceReplace(data)
			case "TYPING_START":
				b.handleTypingStart(data)
			case "CHANNEL_CREATE":
				b.handleChannelCreate(data)
			case "MESSAGE_REACTION_ADD":
				b.handleMessageReactionAdd(data)
			case "MESSAGE_REACTION_REMOVE":
				b.handleMessageReactionRemove(data)
			case "MESSAGE_DELETE":
				b.handleMessageDelete(data)
			case "MESSAGE_UPDATE":
				b.handleMessageUpdate(data)
			case "CHANNEL_PINS_UPDATE":
				b.handleChannelPinsUpdate(data)
			case "GUILD_MEMBER_UPDATE":
				b.handleGuildMemberUpdate(data)
			default:
				log.Errorln("Unhandeled message:", string(message))
				b.handleUnknown(data)
			}
			b.currentSeqNumber = int(data["s"].(float64))
			b.seqNumberChan <- b.currentSeqNumber
		} else if data["op"].(float64) == 9 { // Invalid Session
			log.Errorln("Invalid Session received. Please try again...")
			return
		} else if data["op"].(float64) == 11 { // Heartbeat ACK
			log.Debugln("Heartbeat ACKed from Gateway")
		} else { // opcode which is not handled, yet
			log.Errorf("data: %s", data)
		}
	}
}

// CreateDiscordBot creates a new instance of a DiscordBot
func CreateDiscordBot(token string) (*Bot, error) {
	log.Printf("DiscordBot is CREATING itself using TOKEN = %s", token)
	b := Bot{token: token}
	url := b.getGateway()
	b.ws = dialGateway(url)

	b.sendMessageChan = make(chan events.SendMessage)
	b.commandChan = make(chan events.Command)

	b.knownChannels = make(map[string]channelCreate)

	b.heartBeatStopChan = make(chan struct{})
	b.seqNumberChan = make(chan int)

	b.guilds = make(map[string]guild)
	b.guildNameToID = make(map[string]string)

	b.receivers = make(map[plugins.Plugin]chan events.ReceiveMessage)

	return &b, nil
}

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.sendMessageChan {
		switch sendMsg.Type {
		case events.MESSAGE:
			err := b.sendMessage(sendMsg.Ident, sendMsg.Content)
			if err != nil {
				log.Errorln("Error sending message:", err)
			}
		case events.WHISPER:
			err := b.sendWhisper(sendMsg.Ident, sendMsg.Content)
			if err != nil {
				log.Errorln("Error sending whisper:", err)
			}
		default:
		}
	}
}

func (b *Bot) startCommandChannelReceiver() {
	for cmd := range b.commandChan {
		switch cmd.Command {
		case string("DemoCommand"):
			log.Infoln("Received DemoCommand with server name" + cmd.Payload)
		default:
			log.Errorln("Received unhandeled command" + cmd.Command)
		}
	}
}

// Start the Discord Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	log.Infoln("DiscordBot is STARTING")
	go b.startDiscordBot(doneChannel)
	go b.startSendChannelReceiver()
	go b.startCommandChannelReceiver()
	log.Infoln("DiscordBot is RUNNING")
}

// Stop the Discord Bot
func (b *Bot) Stop() {
	log.Infoln("DiscordBot is SHUTING DOWN")
	close(b.heartBeatStopChan)
	err := b.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Errorln("write close:", err)
	}

	b.disconnectReceivers()

	log.Infoln("DiscordBot is SHUT DOWN")
}

// Status returns the current status of the DiscordBot
func (b *Bot) Status() botinterface.BotStatus {
	status := botinterface.BotStatus{
		Running: true,
		Fail:    false,
		Fatal:   false}
	return status
}

// AddPlugin takes as argument a plugin interface and
// adds it to the DiscordBot by connecting all the required
// channels and starting it
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	plugin.ConnectChannels(b.GetReceiveMessageChannel(plugin), b.GetSendMessageChannel(), b.GetCommandChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
}

func (b *Bot) disconnectReceivers() {
	for plugin, pluginChannel := range b.receivers {
		log.Debugln("Disconnecting Plugin", plugin.GetName())
		defer close(pluginChannel)
	}
}
