package randomplugin

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
)

// RandomPlugin struct holds the private variables for a RandomPlugin
type RandomPlugin struct {
	botReceiveChannel <-chan events.ReceiveMessage
	botSendChannel    chan events.SendMessage

	isStarted bool
}

// GetName returns the name of the plugin
func (p *RandomPlugin) GetName() string {
	return "RandomPlugin"
}

// CreateRandomPlugin returns the struct for a new RandomPlugin
func CreateRandomPlugin() (RandomPlugin, error) {
	log := logging.Get("RandomPlugin")

	log.Printf("RandomPlugin is CREATING itself")
	ep := RandomPlugin{}

	rand.Seed(time.Now().UnixNano())

	return ep, nil
}

func random(max int) int {
	return rand.Intn(max + 1)
}

func (p *RandomPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	log := logging.Get("RandomPlugin")

	log.Printf("Received Message with Type = %s, Ident = %s, content = %s", receivedMessage.Type.String(), receivedMessage.ChannelID, receivedMessage.Content)
	msg := strings.Trim(receivedMessage.Content, " ")
	if p.isStarted && strings.HasPrefix(msg, "!roll") {
		u := stripCmd(msg, "roll")
		if len(msg) == len("!roll") && u == "!roll" {
			u = "100"
		}
		var response string
		num, err := strconv.Atoi(u)
		if err != nil {
			response = fmt.Sprintf("Not a number")
		} else if num <= 0 {
			response = fmt.Sprintf("Number must be > 0")
		} else {
			response = strconv.Itoa(random(num))
		}
		p.botSendChannel <- events.SendMessage{Type: receivedMessage.Type, ChannelID: receivedMessage.ChannelID, Content: response}
	}
}

func (p *RandomPlugin) receiveMessageRunner() {
	log := logging.Get("RandomPlugin")

	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	log.Printf("Automatically SHUTTING DOWN because bot closed the receive channel")
	p.isStarted = false
}

// Start the RandomPlugin
func (p *RandomPlugin) Start() {
	p.isStarted = true
	go p.receiveMessageRunner()
}

// Stop the RandomPlugin
func (p *RandomPlugin) Stop() {
	p.isStarted = false
}

// IsStarted reports if the RandomPlugin is running or not
func (p *RandomPlugin) IsStarted() bool {
	return p.isStarted
}

// ConnectChannels connects the given receive and send channels to the plugin
func (p *RandomPlugin) ConnectChannels(receiveChannel <-chan events.ReceiveMessage,
	sendChannel chan events.SendMessage) error {
	p.botReceiveChannel = receiveChannel
	p.botSendChannel = sendChannel

	return nil
}
