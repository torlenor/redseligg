package matrix

import "github.com/torlenor/abylebotter/events"

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.sendMessageChan {
		switch sendMsg.Type {
		case events.MESSAGE:
			err := b.sendRoomMessage(sendMsg.ChannelID, sendMsg.Content)
			if err != nil {
				log.Println(err)
			}
		case events.WHISPER:
			err := b.sendWhisper(sendMsg.ChannelID, sendMsg.Content)
			if err != nil {
				log.Println(err)
			}
		default:
		}
	}
}
