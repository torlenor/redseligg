package echoplugin

import (
	"strings"

	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
)

// EchoPlugin struct holds the private variables for a EchoPlugin
type EchoPlugin struct {
	botReceiveChannel chan events.ReceiveMessage
	botSendChannel    chan events.SendMessage
	botCommandChannel chan events.Command

	onlyOnWhisper bool

	isStarted bool
}

// GetName returns the name of the plugin
func (p *EchoPlugin) GetName() string {
	return "EchoPlugin"
}

// CreateEchoPlugin returns the struct for a new EchoPlugin
func CreateEchoPlugin() (EchoPlugin, error) {
	log := logging.Get("EchoPlugin")

	log.Printf("EchoPlugin is CREATING itself")
	ep := EchoPlugin{}
	return ep, nil
}

// SetOnlyOnWhisper tells the EchoPlugin that it should only
// echo on WHISPER type messages
func (p *EchoPlugin) SetOnlyOnWhisper(onlyOnWhisper bool) {
	p.onlyOnWhisper = onlyOnWhisper
}

func (p *EchoPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	log := logging.Get("EchoPlugin")

	log.Printf("Received Message with Type = %s, Ident = %s, content = %s", receivedMessage.Type.String(), receivedMessage.Ident, receivedMessage.Content)
	msg := strings.Trim(receivedMessage.Content, " ")
	if p.isStarted && (!p.onlyOnWhisper || receivedMessage.Type == events.WHISPER) && strings.HasPrefix(msg, "!echo") {
		log.Printf("Echoing message back to user = %s, content = %s", receivedMessage.Ident, stripCmd(msg, "echo"))
		p.botSendChannel <- events.SendMessage{Type: receivedMessage.Type, Ident: receivedMessage.Ident, Content: stripCmd(msg, "echo")}
	}
}

func (p *EchoPlugin) receiveMessageRunner() {
	log := logging.Get("EchoPlugin")

	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	log.Printf("Automatically SHUTTING DOWN because bot closed the receive channel")
	p.isStarted = false
}

// Start the EchoPlugin
func (p *EchoPlugin) Start() {
	p.isStarted = true
	go p.receiveMessageRunner()
}

// Stop the EchoPlugin
func (p *EchoPlugin) Stop() {
	p.isStarted = false
}

// IsStarted reports if the EchoPlugin is running or not
func (p *EchoPlugin) IsStarted() bool {
	return p.isStarted
}

// ConnectChannels connects the given receive, send and command channels to the plugin
func (p *EchoPlugin) ConnectChannels(receiveChannel chan events.ReceiveMessage,
	sendChannel chan events.SendMessage,
	commandCHannel chan events.Command) error {
	p.botReceiveChannel = receiveChannel
	p.botSendChannel = sendChannel
	p.botCommandChannel = commandCHannel

	return nil
}
