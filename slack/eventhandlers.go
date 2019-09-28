package slack

import (
	"encoding/json"

	"github.com/torlenor/abylebotter/events"
)

func (b *Bot) handleEventMessage(data []byte) {
	var message EventMessage

	if err := json.Unmarshal([]byte(data), &message); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Printf("%s", data)

	receiveMessage := events.ReceiveMessage{Type: events.MESSAGE, Ident: message.Channel, Content: message.Text}

	for plugin, pluginChannel := range b.receivers {
		b.log.Debugln("Notifying plugin", plugin.GetName(), "about new message/whisper")
		select {
		case pluginChannel <- receiveMessage:
		default:
		}
	}

}
