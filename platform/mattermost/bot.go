package mattermost

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/torlenor/redseligg/botconfig"
	"github.com/torlenor/redseligg/commanddispatcher"
	"github.com/torlenor/redseligg/storage"

	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
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
	platform.BotImpl

	config botconfig.MattermostConfig

	ws *websocket.Conn

	token string

	plugins []plugin.Hooks

	stats stats

	log *logrus.Entry

	lastWsSeqNumber uint32

	MeUser UserObject

	KnownUsers     map[string]userData // key is UserID
	knownUserNames map[string]string   // mapping of UserName to UserID
	knownUserIDs   map[string]string   // mapping of UserID to UserName

	KnownChannels     map[string]channelData // key is ChannelID
	knownChannelNames map[string]string      // mapping of ChannelName to ChannelID
	knownChannelIDs   map[string]string      // mapping of ChannelID to UserChannelNameName
}

func (b *Bot) startMattermostBot() {
	for {
		_, message, err := b.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				b.log.Debugln("Connection closed normally: ", err)
			} else {
				b.log.Errorln("UNHANDLED ERROR: ", err)
			}
			break
		}

		var data map[string]interface{}

		if err := json.Unmarshal(message, &data); err != nil {
			b.log.Errorln("UNHANDLED ERROR: ", err)
			continue
		}

		if event, ok := data["event"]; ok { // Dispatch to event handlers
			switch event {
			case "posted":
				b.handleEventPosted(message)
			default:
				b.log.Warnf("Received unhandled event %s: %s", event, message)
			}
		} else {
			b.log.Warnf("Received unhandled message: %s", message)
		}
	}
}

// CreateMattermostBot creates a new instance of a MattermostBot
func CreateMattermostBot(cfg botconfig.MattermostConfig, storage storage.Storage, commandDispatcher *commanddispatcher.CommandDispatcher) (*Bot, error) {
	log := logging.Get("MattermostBot")
	log.Printf("MattermostBot is CREATING itself")

	b := Bot{
		BotImpl: platform.BotImpl{
			ProvidedFeatures: map[string]bool{
				platform.FeatureMessagePost: true,
			},
			Dispatcher: commandDispatcher,
			Storage:    storage,
		},

		config: cfg,
		log:    log,

		lastWsSeqNumber: 0,

		KnownUsers:     make(map[string]userData),
		knownUserNames: make(map[string]string),
		knownUserIDs:   make(map[string]string),

		KnownChannels:     make(map[string]channelData),
		knownChannelNames: make(map[string]string),
		knownChannelIDs:   make(map[string]string),
	}

	err := b.login()
	if err != nil {
		return nil, fmt.Errorf("Error logging in: %s", err)
	}

	wsServer := strings.Replace(b.config.Server, "https", "wss", 1)
	ws, err := b.dialGateway(wsServer + "/api/v4/websocket")
	if err != nil {
		return nil, err
	}
	b.ws = ws

	b.authWs()

	return &b, nil
}

// Start the Mattermost Bot
func (b *Bot) Start() {
	b.log.Infoln("MattermostBot is STARTING")
	go b.startMattermostBot()
	for _, plugin := range b.plugins {
		plugin.OnRun()
	}
	b.log.Infoln("MattermostBot is RUNNING")
}

// Run the Mattermost Bot (blocking)
func (b *Bot) Run(ctx context.Context) error {
	b.Start()

	<-ctx.Done()

	for _, plugin := range b.plugins {
		plugin.OnStop()
	}

	b.Stop()

	return nil
}

// Stop the Mattermost Bot
func (b *Bot) Stop() {
	b.log.Infoln("MattermostBot is SHUTING DOWN")
	b.log.Infof("MattermostBot Stats:\n%s", b.stats.toString())
	err := b.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		b.log.Errorln("write close:", err)
	}

	b.log.Infoln("MattermostBot is SHUT DOWN")
}

// AddPlugin takes as argument a plugin and
// adds it to the bot providing it with the API
func (b *Bot) AddPlugin(plugin platform.BotPlugin) {
	err := plugin.SetAPI(b)
	if err != nil {
		b.log.Errorf("Could not add plugin %s: %s", plugin.PluginType(), err)
	} else {
		b.plugins = append(b.plugins, plugin)
	}
}

func (b *Bot) addKnownUser(user userData) {
	b.log.Debugf("Added new known User: %s (%s)", user.Username, user.ID)
	b.KnownUsers[user.ID] = user
	b.knownUserNames[user.Username] = user.ID
	b.knownUserIDs[user.ID] = user.Username
}

func (b *Bot) addKnownChannel(channel channelData) {
	b.log.Debugf("Added new known Channel: %s (%s)", channel.ID, channel.Name)
	b.KnownChannels[channel.ID] = channel
	b.knownChannelNames[channel.Name] = channel.ID
	b.knownChannelIDs[channel.ID] = channel.Name
}

// GetInfo returns information about the Bot
func (b *Bot) GetInfo() platform.BotInfo {
	return platform.BotInfo{
		BotID:    "",
		Platform: "Mattermost",
		Healthy:  true,
		Plugins:  []platform.PluginInfo{},
	}
}
