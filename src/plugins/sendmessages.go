package plugins

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"api"
	"events"
	"logging"
)

// SendMessagesPlugin struct holds the private variables for a SendMessagesPlugin
type SendMessagesPlugin struct {
	log *logrus.Entry

	botSendChannel    chan events.SendMessage
	botCommandChannel chan events.Command

	isStarted bool
}

// CreateSendMessagesPlugin returns the struct for a new SendMessagesPlugin
func CreateSendMessagesPlugin() (SendMessagesPlugin, error) {
	ep := SendMessagesPlugin{}
	ep.log = logging.Get("SendMessagesPlugin")
	ep.log.Printf("SendMessagesPlugin is CREATING itself")
	return ep, nil
}

func (p *SendMessagesPlugin) handleSendMessage(ident string, content string) {
	if p.IsStarted() {
		select {
		case p.botSendChannel <- events.SendMessage{Type: events.MESSAGE, Ident: ident, Content: content}:
		default:
		}
	}
}

// Message is a simple struct holding Ident and Content for a message to send
type Message struct {
	Ident   string `json:"ident,omitempty"`
	Content string `json:"content,omitempty"`
}

func (p *SendMessagesPlugin) sendMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err == nil && len(message.Ident) > 0 {
		p.handleSendMessage(message.Ident, message.Content)
		io.WriteString(w, `{"sent":true}`)
	} else {
		io.WriteString(w, `{"sent":false, "error":"Invalid Request"}`)
	}
}

// RegisterToRestAPI registers all endpoints of the plugin to the AbyleBotter REST API
func (p *SendMessagesPlugin) RegisterToRestAPI() {
	api.AttachModulePost("/plugins/sendmessages", p.sendMessage)
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
