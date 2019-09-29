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
	config config.SlackConfig
	log    *logrus.Entry

	stats stats

	ws *websocket.Conn

	sendMessageChan chan events.SendMessage
	commandChan     chan events.Command

	rtmConnectResponse RtmConnectResponse

	receivers map[plugins.Plugin]chan events.ReceiveMessage

	knownPlugins []plugins.Plugin

	channels channelManager
}

// GetReceiveMessageChannel returns the channel which is used to notify
// about received messages from the bot
func (b *Bot) GetReceiveMessageChannel(plugin plugins.Plugin) <-chan events.ReceiveMessage {
	b.log.Debugln("Creating receiveChannel for Plugin", plugin.GetName())
	b.receivers[plugin] = make(chan events.ReceiveMessage)
	return b.receivers[plugin]
}

// GetSendMessageChannel returns the channel which is used to
// send messages using the bot. For SlackBot these messages
// can be normal channel messages, whispers
func (b Bot) GetSendMessageChannel() chan events.SendMessage {
	return b.sendMessageChan
}

// GetCommandChannel gives a channel to control the bot from
// a plugin
func (b Bot) GetCommandChannel() chan events.Command {
	return b.commandChan
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
			case "message":
				b.handleEventMessage(message)
			default:
				b.log.Warnf("Received unhandled event %s: %s", event, message)
			}
		} else {
			b.log.Warnf("Received unhandled message: %s", message)
		}
	}
}

// CreateSlackBot creates a new instance of a SlackBot
func CreateSlackBot(cfg config.SlackConfig) (*Bot, error) {
	log := logging.Get("SlackBot")
	log.Printf("SlackBot is CREATING itself")

	b := Bot{
		config:          cfg,
		log:             log,
		receivers:       make(map[plugins.Plugin]chan events.ReceiveMessage),
		sendMessageChan: make(chan events.SendMessage),
		commandChan:     make(chan events.Command),

		channels: newChannelManager(),
	}

	if len(b.config.Token) == 0 {
		return nil, fmt.Errorf("No Slack token defined in config file")
	}

	rtmConnectResponse, err := b.RtmConnect()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Slack servers: %s", err)
	}
	b.rtmConnectResponse = rtmConnectResponse

	ws, err := b.dialGateway(rtmConnectResponse.URL)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	b.ws = ws

	err = b.populateChannelList()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &b, nil
}

func (b *Bot) populateChannelList() error {
	conversations, err := b.getConversationsList()
	if err != nil {
		return err
	}
	if !conversations.Ok {
		return fmt.Errorf("We received a NOT OK when we tried to get the Conversations List")
	}

	for _, channel := range conversations.Channels {
		b.channels.addKnownChannel(channel)
	}
	return nil
}

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.sendMessageChan {
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
		}
	}
}

func (b *Bot) startCommandChannelReceiver() {
	for cmd := range b.commandChan {
		switch cmd.Command {
		case string("DemoCommand"):
			b.log.Infoln("Received DemoCommand with server name" + cmd.Payload)
		default:
			b.log.Errorln("Received unhandled command" + cmd.Command)
		}
	}
}

// Start the Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	b.log.Infoln("SlackBot is STARTING")
	go b.startSlackBot(doneChannel)
	go b.startSendChannelReceiver()
	go b.startCommandChannelReceiver()
	b.log.Infoln("SlackBot is RUNNING")
}

// Stop the Bot
func (b *Bot) Stop() {
	b.log.Infoln("SlackBot is SHUTING DOWN")
	b.log.Infof("SlackBot Stats:\n%s", b.stats.toString())
	err := b.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		b.log.Errorln("write close:", err)
	}

	b.disconnectReceivers()

	b.log.Infoln("SlackBot is SHUT DOWN")
}

// Status returns the current status of the SlackBot
func (b *Bot) Status() botinterface.BotStatus {
	status := botinterface.BotStatus{
		Running: true,
		Fail:    false,
		Fatal:   false}
	return status
}

// AddPlugin takes as argument a plugin interface and
// adds it to the SlackBot by connecting all the required
// channels and starting it
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	plugin.ConnectChannels(b.GetReceiveMessageChannel(plugin), b.GetSendMessageChannel(), b.GetCommandChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
}

func (b *Bot) disconnectReceivers() {
	for plugin, pluginChannel := range b.receivers {
		b.log.Debugln("Disconnecting Plugin", plugin.GetName())
		defer close(pluginChannel)
	}
}
