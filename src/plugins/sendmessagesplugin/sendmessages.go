package sendmessagesplugin

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"api"
	"events"
	"logging"
)

var registerToAPI = func(path string, f func(http.ResponseWriter, *http.Request)) {
	api.AttachModulePost(path, f)
}

// SendMessagesPlugin struct holds the private variables for a SendMessagesPlugin
type SendMessagesPlugin struct {
	log *logrus.Entry

	botSendChannel    chan events.SendMessage
	botCommandChannel chan events.Command

	isStarted bool
}

// GetName returns the name of the plugin
func (p *SendMessagesPlugin) GetName() string {
	return "SendMessagesPlugin"
}

// CreateSendMessagesPlugin returns the struct for a new SendMessagesPlugin
func CreateSendMessagesPlugin() (SendMessagesPlugin, error) {
	ep := SendMessagesPlugin{}
	ep.log = logging.Get("SendMessagesPlugin")
	ep.log.Printf("SendMessagesPlugin is CREATING itself")
	return ep, nil
}

func (p *SendMessagesPlugin) handleSendMessage(ident string, content string) {
	select {
	case p.botSendChannel <- events.SendMessage{Type: events.MESSAGE, Ident: ident, Content: content}:
	default:
	}
}

// Message is a simple struct holding Ident and Content for a message to send
type Message struct {
	Ident   string `json:"ident,omitempty"`
	Content string `json:"content,omitempty"`
}

func (p *SendMessagesPlugin) sendMessage(message Message) error {
	if p.IsStarted() == false {
		return errors.New("Plugin not started")
	} else if len(message.Ident) == 0 {
		return errors.New("Invalid Request, Ident is empty")
	}

	p.handleSendMessage(message.Ident, message.Content)
	return nil
}

func (p *SendMessagesPlugin) sendMessageRequest(w http.ResponseWriter, r *http.Request) {
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		io.WriteString(w, `{"sent":false, "error":"Invalid Request"}`)
	}
	if err := p.sendMessage(message); err != nil {
		io.WriteString(w, `{"sent":false, "error":"`+err.Error()+`"}`)
	}
	io.WriteString(w, `{"sent":true}`)
}

// RegisterToRestAPI registers all endpoints of the plugin to the AbyleBotter REST API
func (p *SendMessagesPlugin) RegisterToRestAPI() {
	registerToAPI("/plugins/sendmessages", p.sendMessageRequest)
}

// Start the SendMessagesPlugin
func (p *SendMessagesPlugin) Start() {
	p.isStarted = true
}

// Stop the SendMessagesPlugin
func (p *SendMessagesPlugin) Stop() {
	p.isStarted = false
}

// IsStarted reports if the SendMessagesPlugin is running or not
func (p *SendMessagesPlugin) IsStarted() bool {
	return p.isStarted
}

// ConnectChannels connects the given receive, send and command channels to the plugin
func (p *SendMessagesPlugin) ConnectChannels(receiveChannel chan events.ReceiveMessage,
	sendChannel chan events.SendMessage,
	commandCHannel chan events.Command) error {
	p.botSendChannel = sendChannel
	p.botCommandChannel = commandCHannel

	return nil
}
