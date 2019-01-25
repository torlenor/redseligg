package mattermost

// func (b *Bot) dispatchMessage() {
// 	var receiveMessage events.ReceiveMessage
// 	if events.MESSAGE == events.MESSAGE {
// 		b.stats.messagesReceived++
// 		receiveMessage = events.ReceiveMessage{Type: events.WHISPER, Ident: "", Content: ""}
// 	} else {
// 		b.stats.whispersReceived++
// 		receiveMessage = events.ReceiveMessage{Type: events.MESSAGE, Ident: "", Content: ""}
// 	}

// 	for plugin, pluginChannel := range b.receivers {
// 		b.log.Debugln("Notifying plugin", plugin.GetName(), "about new message/whisper")
// 		select {
// 		case pluginChannel <- receiveMessage:
// 		default:
// 		}
// 	}
// }

func (b *Bot) handleUnknown(data map[string]interface{}) {
	b.log.Debugf("TODO HANDLE UNKNOWN EVENT: %s", data["t"])
}
