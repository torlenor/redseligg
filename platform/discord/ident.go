package discord

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func sendIdent(token string, ws webSocketClient) error {
	ident := []byte(`{"op": 2,
			"d": {
				"token": "` + token + `",
				"properties": {},
				"compress": false,
				"large_threshold": 250
			}
}`)

	log.Trace("Sending IDENT to gateway")

	err := ws.SendMessage(websocket.TextMessage, ident)
	if err != nil {
		return fmt.Errorf("Error sending IDENT to gateway: %s", err)
	}

	return nil
}
