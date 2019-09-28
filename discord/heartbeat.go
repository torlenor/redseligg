package discord

import (
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type heartBeatSender interface {
	sendHeartBeat(seqNumber int) error
}

type discordHeartBeatSender struct {
	ws *websocket.Conn
}

func heartBeat(interval int, hbSender heartBeatSender, stop chan struct{}, seqNumber chan int) {
	log.Debugf("Starting heartbeat with interval: %d ms", interval)
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	var currentSeqNumber int

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Debugf("Sending heartbeat (seq number = %d)", currentSeqNumber)
			err := hbSender.sendHeartBeat(currentSeqNumber)
			if err != nil {
				log.Errorln("UNHANDELED ERROR in heartbeat:", err)
			}
		case currentSeqNumber = <-seqNumber:
		case <-stop:
			return
		}
	}
}

func (hbs *discordHeartBeatSender) sendHeartBeat(seqNumber int) error {
	hb := []byte(`{"op":1,"d":` + strconv.Itoa(seqNumber) + `}`)
	return hbs.ws.WriteMessage(websocket.TextMessage, hb)
}
