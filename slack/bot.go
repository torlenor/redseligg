package slack

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/plugincontainer"
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

type webSocketClient interface {
	Dial(wsURL string) error
	Stop()

	ReadMessage() (int, []byte, error)

	SendMessage(messageType int, data []byte) error
	SendJSONMessage(v interface{}) error
}

// The Bot struct holds parameters related to the bot
type Bot struct {
	config config.SlackConfig
	log    *logrus.Entry

	stats stats

	ws webSocketClient

	channels channelManager
	plugins  plugincontainer.PluginContainer
}

func (b *Bot) startSlackBot(doneChannel chan struct{}) {
	defer close(doneChannel)

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

		if event, ok := data["type"]; ok { // Dispatch to event handlers
			switch event {
			case "hello":
				// nothing to do here, just a greeting from the server
			case "message":
				b.handleEventMessage(message)
			case "desktop_notification":
				b.handleEventDesktopNotification(message)
			case "user_typing":
				b.handleEventUserTyping(message)
			case "channel_created":
				b.handleEventChannelCreated(message)
			case "channel_deleted":
				b.handleEventChannelDeleted(message)
			case "channel_joined":
				b.handleEventChannelJoined(message)
			case "channel_left":
				b.handleEventChannelLeft(message)
			case "member_joined_channel":
				b.handleEventMemberJoinedChannel(message)
			case "group_joined":
				b.handleEventGroupJoined(message)
			default:
				b.log.Warnf("Received unhandled event %s: %s", event, message)
			}
		} else {
			b.log.Warnf("Received unhandled message: %s", message)
		}
	}
}

// CreateSlackBot creates a new instance of a SlackBot
func CreateSlackBot(cfg config.SlackConfig, ws webSocketClient) (*Bot, error) {
	log := logging.Get("SlackBot")
	log.Printf("SlackBot is CREATING itself")

	b := Bot{
		config: cfg,
		log:    log,
		ws:     ws,

		channels: newChannelManager(),
		plugins:  plugincontainer.New(),
	}

	if len(b.config.Token) == 0 {
		return nil, fmt.Errorf("No Slack token defined in config file")
	}

	rtmConnectResponse, err := b.RtmConnect()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Slack servers: %s", err)
	}

	err = b.ws.Dial(rtmConnectResponse.URL)
	if err != nil {
		return nil, err
	}

	err = b.populateChannelList()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &b, nil
}

func (b *Bot) populateChannelList() error {
	conversations, err := b.getConversations()
	if err != nil {
		return err
	}

	for _, channel := range conversations {
		b.channels.addKnownChannel(channel)
	}

	b.log.Infof("Added %d known channels", b.channels.Len())
	return nil
}

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.plugins.SendChannel() {
		switch sendMsg.Type {
		case events.MESSAGE:
			err := b.sendMessage(sendMsg.Ident, sendMsg.Content)
			if err != nil {
				b.log.Errorln("Error sending message:", err)
			}
		case events.WHISPER:
			err := b.sendWhisper(sendMsg.Ident, sendMsg.Content)
			if err != nil {
				b.log.Errorln("Error sending whisper:", err)
			}
		default:
			b.log.Warnf("Bot does not support Send Event %s", sendMsg.Type)
		}
	}
}

// Start the Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	b.log.Infof("SlackBot is STARTING (have %d plugin(s))", b.plugins.Size())
	go b.startSlackBot(doneChannel)
	go b.startSendChannelReceiver()
	b.log.Infoln("SlackBot is RUNNING")
}

// Stop the Bot
func (b *Bot) Stop() {
	b.log.Infoln("SlackBot is SHUTING DOWN")
	b.log.Infof("SlackBot Stats:\n%s", b.stats.toString())

	err := b.ws.SendMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		b.log.Warnln("Error when writing close message to ws:", err)
	}

	b.ws.Stop()

	b.plugins.RemoveAll()
	b.log.Infoln("SlackBot is SHUT DOWN")
}

// Status returns the current status of the SlackBot
func (b *Bot) Status() botinterface.BotStatus {
	status := botinterface.BotStatus{
		Running: true,
		Fail:    false,
		Fatal:   false,
	}
	return status
}

// AddPlugin takes as argument a plugin interface and
// adds it to the SlackBot by connecting all the required
// channels and starting it
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	b.plugins.Add(plugin)
}
