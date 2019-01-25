package mattermost

import (
	"botinterface"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logging"
	"net/http"
	"plugins"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"config"
	"events"

	"github.com/gorilla/websocket"
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
	commandChan     chan events.Command

	token string

	receivers map[plugins.Plugin]chan events.ReceiveMessage

	knownPlugins []plugins.Plugin

	stats stats

	log *logrus.Entry

	lastWsSeqNumber uint32

	User UserObject
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

// GetCommandChannel gives a channel to control the bot from
// a plugin
func (b Bot) GetCommandChannel() chan events.Command {
	return b.commandChan
}

type apiResponse struct {
	header http.Header
	body   []byte
}

func (b *Bot) apiCall(path string, method string, body string) (*apiResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, b.config.Server+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+b.token)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &apiResponse{
		body:   responseBody,
		header: response.Header,
	}, nil
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

type loginResponse struct {
	AccessToken string `json:"access_token"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
	DeviceID    string `json:"device_id"`
}

func (b *Bot) login() error {
	// get login server
	response, err := b.apiCall("/api/v4/users/login", "POST", `{"login_id":"`+b.config.Username+`","password":"`+b.config.Password+`"}`)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}

	if val, ok := response.header["Token"]; ok {
		if len(val) > 0 {
			b.token = val[0]
		}
	} else {
		return errors.New("could not login: Response: " + string(response.body))
	}

	user := UserObject{}
	err = json.Unmarshal(response.body, &user)
	return err
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
		commandChan:     make(chan events.Command),
		lastWsSeqNumber: 0,
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
			b.log.Errorln("Received unhandeled command" + cmd.Command)
		}
	}
}

// Start the Discord Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	b.log.Infoln("MattermostBot is STARTING")
	go b.startMattermostBot(doneChannel)
	go b.startSendChannelReceiver()
	go b.startCommandChannelReceiver()
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
	plugin.ConnectChannels(b.GetReceiveMessageChannel(plugin), b.GetSendMessageChannel(), b.GetCommandChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
}

func (b *Bot) disconnectReceivers() {
	for plugin, pluginChannel := range b.receivers {
		b.log.Debugln("Disconnecting Plugin", plugin.GetName())
		defer close(pluginChannel)
	}
}
