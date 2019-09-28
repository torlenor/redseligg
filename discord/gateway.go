package discord

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type gatewayResponse struct {
	URL string `json:"url"`
}

func (b *Bot) getGateway() string {
	log.Debugln("DiscordBot: Requesting the Discord gateway address")
	response, err := b.apiCall("/gateway", "GET", "")
	if err != nil {
		log.Fatalln("FATAL: Could not get the Discord gateway:", err)
	}

	var dat map[string]interface{}

	if err := json.Unmarshal(response, &dat); err != nil {
		log.Fatalln("FATAL: Could not parse the response to our Discord gateway request:", err)
	}

	url := dat["url"].(string)
	log.Debugf("Received Discord gateway address: %s", url)
	return url
}

func dialGateway(gatewayURL string) *websocket.Conn {
	log.Debugln("Dialing the Discord gateway")
	c, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		log.Fatalln("FATAL: Could not dial the Discord gateway:", err)
	}

	return c
}
