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
				break
			} else {
				b.log.Errorln("UNHANDLED ERROR: ", err)
				continue
			}
		}

		var data map[string]interface{}

		if err := json.Unmarshal(message, &data); err != nil {
			b.log.Errorf("Error unmarshalling received message from WebSocket: %s, message was: %s", err, data)
			continue
		}

		if event, ok := data["type"]; ok { // Dispatch to event handlers
			b.eventDispatcher(event, message)
		} else if _, ok := data["ok"]; ok {
			ackMessage := eventAck{}
			if err := json.Unmarshal(message, &ackMessage); err != nil {
				b.log.Errorln("Unable to handle ACK, error unmarshalling JSON:", err)
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
	b.log.Warnf("Encountered an error, trying to restart the bot...")
	b.stopPingWatchdog()
	err := b.ws.SendMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		b.log.Warnln("Error when writing close message to ws:", err)
	}
	b.wg.Wait()
	b.ws.Stop()

	rtmConnectResponse, err := b.RtmConnect()
	if err != nil {
		b.log.Errorf("Error connecting to Slack servers: %s", err)
		return
	}
	b.rtmURL = rtmConnectResponse.URL
	err = b.ws.Dial(b.rtmURL)
	if err != nil {
		b.log.Errorln("Could not dial Slack RTM WebSocket, Slack Bot not operational:", err)
		return
	}

	b.startPingWatchdog()
	go func() {
		b.wg.Add(1)
		b.run()
		defer b.wg.Done()
	}()

	b.log.Info("Recovery attempt finished")
}
