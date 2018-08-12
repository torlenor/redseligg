package discord

import (
	"log"

	"github.com/gorilla/websocket"
)

func sendIdent(token string, ws *websocket.Conn) {
	ident := []byte(`{"op": 2,
			"d": {
				"token": "` + token + `",
				"properties": {},
				"compress": false,
				"large_threshold": 250
			}
}`)

	log.Println("DiscordBot: Sending IDENT to gateway")

	err := ws.WriteMessage(websocket.TextMessage, ident)
	if err != nil {
		log.Fatal("DiscordBot: FATAL: Error sending IDENT to gateway:", err)
	}
}
