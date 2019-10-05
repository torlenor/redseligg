package echoplugin

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
)

// EchoPlugin struct holds the private variables for a EchoPlugin
type EchoPlugin struct {
	log *logrus.Entry

	botReceiveChannel <-chan events.ReceiveMessage
	botSendChannel    chan events.SendMessage

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
	ep := EchoPlugin{
		log: log,
	}

	return ep, nil
}

// SetOnlyOnWhisper tells the EchoPlugin that it should only
// echo on WHISPER type messages
func (p *EchoPlugin) SetOnlyOnWhisper(onlyOnWhisper bool) {
	p.onlyOnWhisper = onlyOnWhisper
}

func (p *EchoPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	p.log.Tracef("Received Message with Type = %s, UserID = %s, content = %s", receivedMessage.Type.String(), receivedMessage.UserID, receivedMessage.Content)
	msg := strings.Trim(receivedMessage.Content, " ")
	if p.isStarted && (!p.onlyOnWhisper || receivedMessage.Type == events.WHISPER) && strings.HasPrefix(msg, "!echo") {
		p.log.Tracef("Echoing message back to user = %s, content = %s", receivedMessage.User, stripCmd(msg, "echo"))
		p.botSendChannel <- events.SendMessage{Type: receivedMessage.Type, ChannelID: receivedMessage.ChannelID, Content: stripCmd(msg, "echo")}
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

// ConnectChannels connects the given receive and send channels to the plugin
func (p *EchoPlugin) ConnectChannels(receiveChannel <-chan events.ReceiveMessage,
	sendChannel chan events.SendMessage) error {
	p.botReceiveChannel = receiveChannel
	p.botSendChannel = sendChannel

	return nil
}
