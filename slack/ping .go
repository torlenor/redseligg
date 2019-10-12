package slack

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Ping is the message we are sending periodically to the server
type Ping struct {
	ID   int       `json:"id"`
	Type string    `json:"type"`
	Time time.Time `json:"time"` // custom field
}

// Pong is an answer to a Ping
type Pong struct {
	ReplyTo int       `json:"reply_to"`
	Type    string    `json:"type"`
	Time    time.Time `json:"time"` // custom field
}

func (b *Bot) sendPing() error {

	ping := Ping{
		ID:   b.idProvider.Get(),
		Type: "ping",
		Time: time.Now(),
	}

	err := b.ws.SendJSONMessage(ping)
	if err != nil {
		return errors.New("Sending Ping failed, probably we are already dead: " + err.Error())
	}

	return nil
}

func (b *Bot) receivePong(data []byte) {
	var pong Pong

	if err := json.Unmarshal(data, &pong); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.watchdog.Feed()
}
