package plugins

import (
	"logging"
	"strings"

	"events"
)

// EchoPlugin struct holds the private variables for a EchoPlugin
type EchoPlugin struct {
	botReceiveChannel chan events.ReceiveMessage
	botSendChannel    chan events.SendMessage
	botCommandChannel chan events.Command

	onlyOnWhisper bool
}

// CreateEchoPlugin returns the struct for a new EchoPlugin
func CreateEchoPlugin() EchoPlugin {
	log := logging.Get("EchoPlugin")

	log.Printf("EchoPlugin is CREATING itself")
	ep := EchoPlugin{}
	return ep
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
	if (!p.onlyOnWhisper || receivedMessage.Type == events.WHISPER) && strings.HasPrefix(msg, "!echo") {
		log.Printf("Echoing message back to user = %s, content = %s", receivedMessage.Ident, stripCmd(msg, "echo"))
		select {
		case p.botSendChannel <- events.SendMessage{Type: events.WHISPER, Ident: receivedMessage.Ident, Content: stripCmd(msg, "echo")}:
		default:
		}
	}
}

func (p *EchoPlugin) receiveMessageRunner() {
	log := logging.Get("EchoPlugin")

	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	log.Printf("Automatically SHUTTING DOWN because bot closed the receive channel")
}

// Start the EchoPlugin
func (p *EchoPlugin) Start() {
	go p.receiveMessageRunner()
}

// Stop the EchoPlugin
func (p *EchoPlugin) Stop() {

}

// ConnectChannels connects the given receive, send and command channels to the plugin
func (p *EchoPlugin) ConnectChannels(receiveChannel chan events.ReceiveMessage, sendChannel chan events.SendMessage, commandCHannel chan events.Command) {
	p.botReceiveChannel = receiveChannel
	p.botSendChannel = sendChannel
	p.botCommandChannel = commandCHannel
}
