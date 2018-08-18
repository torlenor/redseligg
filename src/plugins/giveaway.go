package plugins

import (
	"logging"
	"strings"
	"time"

	"events"
)

// GiveawayPlugin struct holds the private variables for a GiveawayPlugin
type GiveawayPlugin struct {
	botReceiveChannel chan events.ReceiveMessage
	botSendChannel    chan events.SendMessage
	botCommandChannel chan events.Command

	giveaways     []giveaway
	giveawayStart bool
}

type giveaway struct {
	userID    string
	channelID string
	users     []string
	started   time.Time
	runtime   time.Time
	winner    string
	isInPrep  bool
}

// CreateGiveawayPlugin returns the struct for a new GiveawayPlugin
func CreateGiveawayPlugin(receiveChannel chan events.ReceiveMessage, sendChannel chan events.SendMessage) GiveawayPlugin {
	log := logging.Get("GiveawayPlugin")

	log.Printf("GiveawayPlugin is CREATING itself")
	ep := GiveawayPlugin{botReceiveChannel: receiveChannel,
		botSendChannel: sendChannel}
	return ep
}

func (p *GiveawayPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	log := logging.Get("GiveawayPlugin")

	log.Printf("Received Message with Type = %s, Ident = %s, content = %s", receivedMessage.Type.String(), receivedMessage.Ident, receivedMessage.Content)
	if p.giveawayStart == false {
		msg := strings.Trim(receivedMessage.Content, " ")
		if strings.HasPrefix(msg, "!giveaway") {
			log.Printf("Echoing message back to user = %s, content = %s", receivedMessage.Ident, stripCmd(msg, "echo"))
			select {
			case p.botSendChannel <- events.SendMessage{Type: events.WHISPER, Ident: receivedMessage.Ident, Content: stripCmd(msg, "echo")}:
			default:
			}
		}
	} else {
		// already a giveaway start in progress
		for i := range p.giveaways {
			if p.giveaways[i].userID == receivedMessage.Ident {
				// Found!
			}
		}
	}
}

func (p GiveawayPlugin) receiveMessageRunner() {
	log := logging.Get("GiveawayPlugin")

	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	log.Printf("Automatically SHUTTING DOWN because bot closed the receive channel")
}

// Start the GiveawayPlugin
func (p *GiveawayPlugin) Start() {
	p.giveawayStart = false
	go p.receiveMessageRunner()
}

// Stop the GiveawayPlugin
func (p *GiveawayPlugin) Stop() {

}

// ConnectChannels connects the given receive, send and command channels to the plugin
func (p *GiveawayPlugin) ConnectChannels(receiveChannel chan events.ReceiveMessage, sendChannel chan events.SendMessage, commandCHannel chan events.Command) {
	p.botReceiveChannel = receiveChannel
	p.botSendChannel = sendChannel
	p.botCommandChannel = commandCHannel
}
