package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"

	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin"
)

var (
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

type stats struct {
	messagesSent int64
	whispersSent int64

	messagesReceived int64
	whispersReceived int64
}

func (s stats) toString() string {
	return fmt.Sprintf("Messages Sent: %d\nMessages Received: %d\nWhispers Sent: %d\nWhispers Received: %d",
		s.messagesSent, s.messagesReceived, s.whispersSent, s.whispersReceived)
}

// The Bot struct holds parameters related to the bot
type Bot struct {
	ws *websocket.Conn

	knownChannels     map[string]channelCreate
	token             string
	ownSnowflakeID    string
	currentSeqNumber  int
	heartBeatSender   *discordHeartBeatSender
	heartBeatStopChan chan struct{}
	seqNumberChan     chan int

	plugins []plugin.Hooks

	guilds        map[string]guildCreate // map[ID]
	guildNameToID map[string]string

	stats stats

	discordOauthConfig *oauth2.Config

	oauth2Handler *oauth2Handler
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

func (b *Bot) startDiscordBot() {
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
			case "PRESENCES_REPLACE":
				b.handlePresencesReplace(data)
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
func CreateDiscordBot(cfg botconfig.DiscordConfig) (*Bot, error) {
	log.Info("DiscordBot is CREATING itself")

	b := Bot{token: cfg.Token}
	url := b.getGateway()
	b.ws = dialGateway(url)

	b.knownChannels = make(map[string]channelCreate)

	b.heartBeatStopChan = make(chan struct{})
	b.seqNumberChan = make(chan int)

	b.guilds = make(map[string]guildCreate)
	b.guildNameToID = make(map[string]string)

	b.discordOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/cb",
		ClientID:     cfg.ID,
		ClientSecret: cfg.Secret,
		Scopes:       []string{"bot"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discordapp.com/api/oauth2/authorize",
			TokenURL: "https://discordapp.com/api/oauth2/token",
		},
	}

	b.oauth2Handler = createOAuth2Handler(*b.discordOauthConfig)

	return &b, nil
}

// Start the Discord Bot
func (b *Bot) Start() {
	log.Infoln("DiscordBot is STARTING")
	go b.startDiscordBot()
	go b.oauth2Handler.startOAuth2Handler()
	log.Infoln("DiscordBot is RUNNING")
}

// Run the Discord Bot (blocking)
func (b *Bot) Run(ctx context.Context) error {
	b.Start()

	<-ctx.Done()

	b.Stop()

	return nil
}

// Stop the Discord Bot
func (b *Bot) Stop() {
	log.Infoln("DiscordBot is SHUTING DOWN")
	log.Infof("DiscordBot Stats:\n%s", b.stats.toString())
	close(b.heartBeatStopChan)
	err := b.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Errorln("write close:", err)
	}

	log.Infoln("DiscordBot is SHUT DOWN")
}

// AddPlugin takes as argument a plugin and
// adds it to the bot providing it with the API
func (b *Bot) AddPlugin(plugin platform.BotPlugin) {
	plugin.SetAPI(b)
	b.plugins = append(b.plugins, plugin)
}

// GetInfo returns information about the Bot
func (b *Bot) GetInfo() platform.BotInfo {
	return platform.BotInfo{
		BotID:    "",
		Platform: "Discord",
		Healthy:  true,
		Plugins:  []platform.PluginInfo{},
	}
}
