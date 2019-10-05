package httppingplugin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
)

// HTTPPingPlugin struct holds the private variables for a HTTPPingPlugin
type HTTPPingPlugin struct {
	botReceiveChannel <-chan events.ReceiveMessage
	botSendChannel    chan events.SendMessage

	isStarted bool
}

// GetName returns the name of the plugin
func (p *HTTPPingPlugin) GetName() string {
	return "HTTPPingPlugin"
}

// CreateHTTPPingPlugin returns the struct for a new HTTPPingPlugin
func CreateHTTPPingPlugin() (HTTPPingPlugin, error) {
	log := logging.Get("HTTPPingPlugin")

	log.Printf("HTTPPingPlugin is CREATING itself")
	ep := HTTPPingPlugin{}
	return ep, nil
}

func httpPing(u string) (int, error) {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return -1, fmt.Errorf("Not a valid url")
	}

	start := time.Now()
	resp, err := http.Get(u)
	elapsed := int(time.Since(start).Nanoseconds() / 1000 / 1000)
	if err != nil {
		return elapsed, fmt.Errorf("Error pinging the url: %s", err)
	}
	defer resp.Body.Close()
	return elapsed, nil
}

func (p *HTTPPingPlugin) handleReceivedMessage(receivedMessage events.ReceiveMessage) {
	log := logging.Get("HTTPPingPlugin")

	log.Printf("Received Message with Type = %s, Ident = %s, content = %s", receivedMessage.Type.String(), receivedMessage.Ident, receivedMessage.Content)
	msg := strings.Trim(receivedMessage.Content, " ")
	if p.isStarted && strings.HasPrefix(msg, "!httpping") {
		u := stripCmd(msg, "httpping")
		timeMs, err := httpPing(u)
		var response string
		if err != nil {
			response = fmt.Sprintf("FAIL (%s). Request took %d ms", err, timeMs)
		} else {
			response = fmt.Sprintf("SUCCESS. Request took %d ms", timeMs)
		}
		p.botSendChannel <- events.SendMessage{Type: receivedMessage.Type, Ident: receivedMessage.Ident, Content: response}
	}
}

func (p *HTTPPingPlugin) receiveMessageRunner() {
	log := logging.Get("HTTPPingPlugin")

	for receivedMessage := range p.botReceiveChannel {
		p.handleReceivedMessage(receivedMessage)
	}
	log.Printf("Automatically SHUTTING DOWN because bot closed the receive channel")
	p.isStarted = false
}

// Start the HTTPPingPlugin
func (p *HTTPPingPlugin) Start() {
	p.isStarted = true
	go p.receiveMessageRunner()
}

// Stop the HTTPPingPlugin
func (p *HTTPPingPlugin) Stop() {
	p.isStarted = false
}

// IsStarted reports if the HTTPPingPlugin is running or not
func (p *HTTPPingPlugin) IsStarted() bool {
	return p.isStarted
}

// ConnectChannels connects the given receive and send channels to the plugin
func (p *HTTPPingPlugin) ConnectChannels(receiveChannel <-chan events.ReceiveMessage,
	sendChannel chan events.SendMessage) error {
	p.botReceiveChannel = receiveChannel
	p.botSendChannel = sendChannel

	return nil
}
