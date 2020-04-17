package discord

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

func sendIdent(token string, ws webSocketClient) error {
	ident := []byte(`{"op": 2,
			"d": {
				"token": "` + token + `",
				"properties": {
					"$os": "linux",
					"$browser": "abylebotter",
					"$device": "abylebotter"
				  },
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

func sendResume(token string, sessionID string, seq int, ws webSocketClient) error {
	ident := []byte(`{"op": 6,
		"d": {
			"token": "` + token + `",
		  "session_id": "` + sessionID + `",
		  "seq": ` + strconv.Itoa(seq) + `
		}
	  }`)

	log.Trace("Sending RESUME to gateway")

	err := ws.SendMessage(websocket.TextMessage, ident)
	if err != nil {
		return fmt.Errorf("Error sending RESUME to gateway: %s", err)
	}

	return nil
}
