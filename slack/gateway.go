package slack

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (b *Bot) dialGateway(gatewayURL string) (*websocket.Conn, error) {
	b.log.Debugf("Dialing the Slack gateway: %s", gatewayURL)
	c, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not dial the Slack gateway")
	}

	return c, nil
}
