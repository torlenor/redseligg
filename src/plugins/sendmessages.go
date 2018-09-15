package plugins

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"events"
	"logging"
)

// SendMessagesPlugin struct holds the private variables for a SendMessagesPlugin
type SendMessagesPlugin struct {
	log *logrus.Entry

	botReceiveChannel chan events.ReceiveMessage
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

func (p *SendMessagesPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	p.log.Printf("Received Message with Type = %s, Ident = %s, content = %s", receivedMessage.Type.String(), receivedMessage.Ident, receivedMessage.Content)
}

func (p *SendMessagesPlugin) handleSendMessage(ident string, content string) {
	if p.isStarted {
		select {
		case p.botSendChannel <- events.SendMessage{Type: events.MESSAGE, Ident: ident, Content: content}:
		default:
		}
	}
}

func (p *SendMessagesPlugin) receiveMessageRunner() {
	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	p.log.Printf("Automatically SHUTTING DOWN because bot closed the receive channel")
	p.isStarted = false
}

// Message is a simple struct holding Ident and Content for a message to send
type Message struct {
	Ident   string `json:"ident,omitempty"`
	Content string `json:"content,omitempty"`
}

func (p *SendMessagesPlugin) sendMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	_ = json.NewDecoder(r.Body).Decode(&message)
	p.handleSendMessage(message.Ident, message.Content)
	json.NewEncoder(w).Encode(message)
}

func (p *SendMessagesPlugin) startRestAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/plugins/sendmessages", p.sendMessage).Methods("POST")
	p.log.Fatal(http.ListenAndServe(":8000", router))
}

// Start the SendMessagesPlugin
func (p *SendMessagesPlugin) Start() {
	p.isStarted = true
	go p.startRestAPI()
	go p.receiveMessageRunner()
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
	p.botReceiveChannel = receiveChannel
	p.botSendChannel = sendChannel
	p.botCommandChannel = commandCHannel

	return nil
}
