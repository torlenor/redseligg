package discord

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/torlenor/AbyleBotter/events"
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
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	ws                 *websocket.Conn
	receiveMessageChan chan events.ReceiveMessage
	knownChannels      map[string]channelCreate
	token              string
	ownSnowflakeID     string
	currentSeqNumber   int
}

// GetReceiveMessageChannel returns the channel which is used to notify
// about received messages from the bot. For DiscordBot these messages
// can be normal channel messages, whispers
func (b Bot) GetReceiveMessageChannel() chan events.ReceiveMessage {
	return b.receiveMessageChan
}

func (b Bot) apiCall(path string, method string, body string) (r []byte, e error) {
	client := &http.Client{}

	req, _ := http.NewRequest(method, "https://discordapp.com/api"+path, strings.NewReader(body))

	req.Header.Add("Authorization", "Bot "+b.token)
	req.Header.Add("Content-Type", "application/json")

	response, _ := client.Do(req)

	return ioutil.ReadAll(response.Body)
}

func (b *Bot) startDiscordBot(doneChannel chan struct{}, ws *websocket.Conn) {
	defer close(doneChannel)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("DiscordBot: Connection closed normally: ", err)
			} else {
				log.Println("DiscordBot: UNHANDELED ERROR: ", err)
			}
			break
		}

		var data map[string]interface{}

		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("DiscordBot: UNHANDELED ERROR: ", err)
			continue
		}

		if data["op"].(float64) == 10 { // Hello from Discord Gateway
			log.Println("DiscordBot: Received HELLO from gateway")
			sendIdent(b.token, ws)
			heartbeatInterval := int(data["d"].(map[string]interface{})["heartbeat_interval"].(float64))
			go b.heartBeat(heartbeatInterval, ws) // Start sending heartbeats
			log.Println("DiscordBot: DiscordBot is READY")
		} else if data["op"].(float64) == 0 { // Dispatch to event handlers
			switch data["t"] {
			case "MESSAGE_CREATE":
				b.handleMessageCreate(data)
			case "READY":
				b.handleReady(data)
			case "GUILD_CREATE":
				handleGuildCreate(data)
			case "PRESENCE_UPDATE":
				handlePresenceUpdate(data)
			case "PRESENCE_REPLACE":
				log.Println(string(message))
				handlePresenceReplace(data)
			case "TYPING_START":
				handleTypingStart(data)
			case "CHANNEL_CREATE":
				b.handleChannelCreate(data)
			case "MESSAGE_REACTION_ADD":
				handleMessageReactionAdd(data)
			case "MESSAGE_REACTION_REMOVE":
				handleMessageReactionRemove(data)
			case "MESSAGE_DELETE":
				handleMessageDelete(data)
			case "MESSAGE_UPDATE":
				handleMessageUpdate(data)
			case "CHANNEL_PINS_UPDATE":
				handleCHannelPinsUpdate(data)
			case "GUILD_MEMBER_UPDATE":
				handleGuildMemberUpdate(data)
			default:
				log.Println(string(message))
				handleUnknown(data)
			}
			b.currentSeqNumber = int(data["s"].(float64))
		} else if data["op"].(float64) == 9 { // Invalid Session
			log.Printf("DiscordBot: Invalid Session received. Please try again...")
			return
		} else if data["op"].(float64) == 11 { // Heartbeat ACK
			log.Printf("DiscordBot: Heartbeat ACKed from Gateway")
		} else { // opcode which is not handled, yet
			log.Printf("data: %s", data)
		}
	}
}

// CreateDiscordBot creates a new instance of a DiscordBot
func CreateDiscordBot(token string) *Bot {
	log.Printf("DiscordBot: DiscordBot is CREATING itself using TOKEN = %s", token)
	b := Bot{token: token}
	url := b.getGateway()
	b.ws = dialGateway(url)

	b.receiveMessageChan = make(chan events.ReceiveMessage)
	b.knownChannels = make(map[string]channelCreate)

	return &b
}

// Start the Discord Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	log.Println("DiscordBot: DiscordBot is STARTING")
	go b.startDiscordBot(doneChannel, b.ws)
}

// Stop the Discord Bot
func (b Bot) Stop() {
	log.Println("DiscordBot: DiscordBot is SHUTING DOWN")
	err := b.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
	}
	defer close(b.receiveMessageChan)
}
