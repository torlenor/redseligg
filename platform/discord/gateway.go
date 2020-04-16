package discord

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type gatewayResponse struct {
	URL string `json:"url"`
}

func (b *Bot) getGateway() (string, error) {
	log.Traceln("DiscordBot: Requesting the Discord gateway address")
	response, err := b.api.Call("/gateway", "GET", "")
	if err != nil {
		return "", fmt.Errorf("Could not get the Discord gateway: %s", err.Error())
	}

	var dat map[string]interface{}

	if err := json.Unmarshal(response.Body, &dat); err != nil {
		return "", fmt.Errorf("Could not parse the response to our Discord gateway request: %s", err.Error())
	}

	url := dat["url"].(string)
	log.Tracef("Received Discord gateway address: %s", url)
	return url, nil
}

func dialGateway(gatewayURL string) *websocket.Conn {
	log.Debugln("Dialing the Discord gateway")
	c, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		log.Fatalln("FATAL: Could not dial the Discord gateway:", err)
	}

	return c
}
