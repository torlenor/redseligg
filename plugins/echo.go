package plugins

import (
	"log"

	"github.com/torlenor/AbyleBotter/botinterface"
	"github.com/torlenor/AbyleBotter/events"
)

// Plugin type interface
type Plugin interface {
	RegisterBot(bot botinterface.Bot)
	ReceiveMessage()
}

// EchoPlugin struct holds the private variables for a EchoPlugin
type EchoPlugin struct {
	botReceiveChannel chan events.ReceiveMessage
	onlyOnWhisper     bool
}

// CreateEchoPlugin returns the struct for a new EchoPlugin
func CreateEchoPlugin(receiveChannel chan events.ReceiveMessage) EchoPlugin {
	log.Printf("EchoBot: EchoBot is CREATING itself")
	ep := EchoPlugin{botReceiveChannel: receiveChannel}
	return ep
}

// SetOnlyOnWhisper tells the EchoPlugin that it should only
// echo on WHISPER type messages
func (p *EchoPlugin) SetOnlyOnWhisper(onlyOnWhisper bool) {
	p.onlyOnWhisper = onlyOnWhisper
}

func (p EchoPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	log.Printf("EchoBot: Received Message with Type = %s, Ident = %s, content = %s", receivedMessage.Type.String(), receivedMessage.Ident, receivedMessage.Content)
	if !p.onlyOnWhisper || receivedMessage.Type == events.WHISPER {
		log.Printf("EchoBot: Echoing message back to user = %s, content = %s", receivedMessage.Ident, receivedMessage.Content)
	}
}

func (p EchoPlugin) receiveMessageRunner() {
	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	log.Printf("EchoBot: Automatically SHUTTING DOWN because bot closed the receive channel")
}

// Start the EchoPlugin
func (p EchoPlugin) Start() {
	go p.receiveMessageRunner()
}
