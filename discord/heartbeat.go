package discord

import (
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func (b *Bot) heartBeat(interval int, ws *websocket.Conn) {
	log.Printf("DiscordBot: Starting heartbeat with interval: %d ms", interval)
	ticker := time.NewTicker(time.Duration(interval/1000) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			hb := []byte(`{"op":1,"d":` + strconv.Itoa(b.currentSeqNumber) + `}`)

			log.Printf("DiscordBot: Sending heartbeat (seq number = %d)", b.currentSeqNumber)
			err := ws.WriteMessage(websocket.TextMessage, hb)
			if err != nil {
				log.Println("DiscordBot: UNHANDELED ERROR in heartbeat:", err)
				return
			}
		}
	}
}
