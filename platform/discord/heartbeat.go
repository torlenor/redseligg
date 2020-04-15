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
	ws webSocketClient
}

func heartBeat(interval time.Duration, hbSender heartBeatSender, stop chan bool, seqNumber chan int, onFail func()) {
	log.Debugf("Starting heartbeat with interval: %d ms", interval.Milliseconds())
	ticker := time.NewTicker(interval)

	var currentSeqNumber int

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Debugf("Sending heartbeat (seq number = %d)", currentSeqNumber)
			err := hbSender.sendHeartBeat(currentSeqNumber)
			if err != nil {
				log.Errorln("UNHANDLED ERROR in heartbeat:", err)
				go onFail()
			}
		case currentSeqNumber = <-seqNumber:
		case <-stop:
			return
		}
	}
}

func (hbs *discordHeartBeatSender) sendHeartBeat(seqNumber int) error {
	hb := []byte(`{"op":1,"d":` + strconv.Itoa(seqNumber) + `}`)
	return hbs.ws.SendMessage(websocket.TextMessage, hb)
}
