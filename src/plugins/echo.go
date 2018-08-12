package plugins

import (
	"log"
	"strings"

	"botinterface"
	"events"
)

// Plugin type interface
type Plugin interface {
	RegisterBot(bot botinterface.Bot)
	ReceiveMessage()
}

// EchoPlugin struct holds the private variables for a EchoPlugin
type EchoPlugin struct {
	botReceiveChannel chan events.ReceiveMessage
	botSendChannel    chan events.SendMessage
	onlyOnWhisper     bool
}

// CreateEchoPlugin returns the struct for a new EchoPlugin
func CreateEchoPlugin(receiveChannel chan events.ReceiveMessage, sendChannel chan events.SendMessage) EchoPlugin {
	log.Printf("EchoPlugin: EchoPlugin is CREATING itself")
	ep := EchoPlugin{botReceiveChannel: receiveChannel,
		botSendChannel: sendChannel}
	return ep
}

// SetOnlyOnWhisper tells the EchoPlugin that it should only
// echo on WHISPER type messages
func (p *EchoPlugin) SetOnlyOnWhisper(onlyOnWhisper bool) {
	p.onlyOnWhisper = onlyOnWhisper
}

func (p EchoPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	log.Printf("EchoPlugin: Received Message with Type = %s, Ident = %s, content = %s", receivedMessage.Type.String(), receivedMessage.Ident, receivedMessage.Content)
	msg := strings.Trim(receivedMessage.Content, " ")
	if (!p.onlyOnWhisper || receivedMessage.Type == events.WHISPER) && strings.HasPrefix(msg, "!echo") {
		log.Printf("EchoPlugin: Echoing message back to user = %s, content = %s", receivedMessage.Ident, stripCmd(msg, "echo"))
		select {
		case p.botSendChannel <- events.SendMessage{Type: events.WHISPER, Ident: receivedMessage.Ident, Content: stripCmd(msg, "echo")}:
		default:
		}
	}
}

func (p EchoPlugin) receiveMessageRunner() {
	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	log.Printf("EchoPlugin: Automatically SHUTTING DOWN because bot closed the receive channel")
}

// Start the EchoPlugin
func (p EchoPlugin) Start() {
	go p.receiveMessageRunner()
}
