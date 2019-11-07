package mattermost

import (
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (b *Bot) dialGateway(gatewayURL string) (*websocket.Conn, error) {
	b.log.Debugf("Dialing the Mattermost gateway: %s", gatewayURL)
	c, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not dial the Mattermost gateway")
	}

	return c, nil
}

func (b *Bot) authWs() {
	b.lastWsSeqNumber++
	ident := []byte(`{
		"seq": ` + strconv.Itoa(int(b.lastWsSeqNumber)) + `,
		"action": "authentication_challenge",
		"data": {
		  "token": "` + b.token + `"
		}
	  }`)

	b.log.Println("Sending AUTH to gateway")

	err := b.ws.WriteMessage(websocket.TextMessage, ident)
	if err != nil {
		b.log.Fatal("Error sending AUTH to gateway:", err)
	}
}
