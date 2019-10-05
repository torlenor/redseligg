package mattermost

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/plugins"
)

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
	config config.MattermostConfig

	ws *websocket.Conn

	sendMessageChan chan events.SendMessage

	token string

	receivers map[plugins.Plugin]chan events.ReceiveMessage

	knownPlugins []plugins.Plugin

	stats stats

	log *logrus.Entry

	lastWsSeqNumber uint32

	MeUser UserObject

	KnownUsers     map[string]User   // key is UserID
	knownUserNames map[string]string // mapping of UserName to UserID
	knownUserIDs   map[string]string // mapping of UserID to UserName

	KnownChannels     map[string]Channel // key is ChannelID
	knownChannelNames map[string]string  // mapping of ChannelName to ChannelID
	knownChannelIDs   map[string]string  // mapping of ChannelID to UserChannelNameName
}

// GetReceiveMessageChannel returns the channel which is used to notify
// about received messages from the bot
func (b *Bot) GetReceiveMessageChannel(plugin plugins.Plugin) chan events.ReceiveMessage {
	b.log.Debugln("Creating receiveChannel for Plugin", plugin.GetName())
	b.receivers[plugin] = make(chan events.ReceiveMessage)
	return b.receivers[plugin]
}

// GetSendMessageChannel returns the channel which is used to
// send messages using the bot. For MattermostBot these messages
// can be normal channel messages, whispers
func (b Bot) GetSendMessageChannel() chan events.SendMessage {
	return b.sendMessageChan
}

func (b *Bot) startMattermostBot(doneChannel chan struct{}) {
	defer close(doneChannel)

	for {
		_, message, err := b.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				b.log.Debugln("Connection closed normally: ", err)
			} else {
				b.log.Errorln("UNHANDELED ERROR: ", err)
			}
			break
		}

		var data map[string]interface{}

		if err := json.Unmarshal(message, &data); err != nil {
			b.log.Errorln("UNHANDELED ERROR: ", err)
			continue
		}

		if event, ok := data["event"]; ok { // Dispatch to event handlers
			switch event {
			case "posted":
				b.handleEventPosted(message)
			default:
				b.log.Warnf("Received unhandeled event %s: %s", event, message)
			}
		} else {
			b.log.Warnf("Received unhandeled message: %s", message)
		}
	}
}

// CreateMattermostBot creates a new instance of a MattermostBot
func CreateMattermostBot(cfg config.MattermostConfig) (*Bot, error) {
	log := logging.Get("MattermostBot")
	log.Printf("MattermostBot is CREATING itself")

	b := Bot{
		config:          cfg,
		log:             log,
		receivers:       make(map[plugins.Plugin]chan events.ReceiveMessage),
		sendMessageChan: make(chan events.SendMessage),

		lastWsSeqNumber: 0,

		KnownUsers:     make(map[string]User),
		knownUserNames: make(map[string]string),
		knownUserIDs:   make(map[string]string),

		KnownChannels:     make(map[string]Channel),
		knownChannelNames: make(map[string]string),
		knownChannelIDs:   make(map[string]string),
	}

	if b.config.UseToken == true {
		b.token = b.config.Token
	} else {
		err := b.login()
		if err != nil {
			b.log.Fatalf("Error logging in: %s", err)
		}
	}

	wsServer := strings.Replace(b.config.Server, "https", "wss", 1)
	ws, err := b.dialGateway(wsServer + "/api/v4/websocket")
	if err != nil {
		b.log.Fatalf(err.Error())
	}
	b.ws = ws

	b.authWs()

	return &b, nil
}

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.sendMessageChan {
		switch sendMsg.Type {
		case events.MESSAGE:
			err := b.sendMessage(sendMsg.ChannelID, sendMsg.Content)
			if err != nil {
				b.log.Errorln("Error sending message:", err)
			}
		case events.WHISPER:
			err := b.sendWhisper(sendMsg.ChannelID, sendMsg.Content)
			if err != nil {
				b.log.Errorln("Error sending whisper:", err)
			}
		default:
		}
	}
}

// Start the Discord Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	b.log.Infoln("MattermostBot is STARTING")
	go b.startMattermostBot(doneChannel)
	go b.startSendChannelReceiver()
	b.log.Infoln("MattermostBot is RUNNING")
}

// Stop the Discord Bot
func (b *Bot) Stop() {
	b.log.Infoln("MattermostBot is SHUTING DOWN")
	b.log.Infof("MattermostBot Stats:\n%s", b.stats.toString())
	err := b.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		b.log.Errorln("write close:", err)
	}

	b.disconnectReceivers()

	b.log.Infoln("MattermostBot is SHUT DOWN")
}

// Status returns the current status of the MattermostBot
func (b *Bot) Status() botinterface.BotStatus {
	status := botinterface.BotStatus{
		Running: true,
		Fail:    false,
		Fatal:   false}
	return status
}

// AddPlugin takes as argument a plugin interface and
// adds it to the MattermostBot by connecting all the required
// channels and starting it
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	plugin.ConnectChannels(b.GetReceiveMessageChannel(plugin), b.GetSendMessageChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
}

func (b *Bot) disconnectReceivers() {
	for plugin, pluginChannel := range b.receivers {
		b.log.Debugln("Disconnecting Plugin", plugin.GetName())
		defer close(pluginChannel)
	}
}

func (b *Bot) addKnownUser(user User) {
	b.log.Debugf("Added new known User: %s (%s)", user.Username, user.ID)
	b.KnownUsers[user.ID] = user
	b.knownUserNames[user.Username] = user.ID
	b.knownUserIDs[user.ID] = user.Username
}

func (b *Bot) addKnownChannel(channel Channel) {
	b.log.Debugf("Added new known Channel: %s (%s)", channel.ID, channel.Name)
	b.KnownChannels[channel.ID] = channel
	b.knownChannelNames[channel.Name] = channel.ID
	b.knownChannelIDs[channel.ID] = channel.Name
}
