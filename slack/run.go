package slack

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

func (b *Bot) run() {
	for {
		_, message, err := b.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				b.log.Debugln("Connection closed normally: ", err)
			} else {
				b.log.Errorln("UNHANDLED ERROR: ", err)
			}
			break
		}

		var data map[string]interface{}

		if err := json.Unmarshal(message, &data); err != nil {
			b.log.Errorln("UNHANDLED ERROR: ", err)
			continue
		}

		if event, ok := data["type"]; ok { // Dispatch to event handlers
			b.eventDispatcher(event, message)
		} else if _, ok := data["ok"]; ok {
			ackMessage := EventAck{}
			if err := json.Unmarshal(message, &ackMessage); err != nil {
				b.log.Errorln("UNHANDLED ERROR: ", err)
			} else {
				b.log.Debugf("Received an ACK to a message we sent, not used yet: %s", message)
			}
		} else {
			b.log.Warnf("Received unhandled message: %s", message)
		}
	}
}

func pingSender(interval time.Duration, f func() error, stop chan bool) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-stop:
			ticker.Stop()
			return
		case <-ticker.C:
			f()
		}
	}
}

func (b *Bot) onFail() {
	b.log.Debugf("TODO: Implement recover onFail")
}
