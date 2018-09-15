package matrix

import "events"

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.sendMessageChan {
		switch sendMsg.Type {
		case events.MESSAGE:
			err := b.sendRoomMessage(sendMsg.Ident, sendMsg.Content)
			if err != nil {
				log.Println(err)
			}
		case events.WHISPER:
			err := b.sendWhisper(sendMsg.Ident, sendMsg.Content)
			if err != nil {
				log.Println(err)
			}
		default:
		}
	}
}

func (b *Bot) startCommandChannelReceiver() {
	for cmd := range b.commandChan {
		switch cmd.Command {
		case string("DemoCommand"):
			log.Println("Received DemoCommand with server name" + cmd.Payload)
		default:
			log.Println("Received unhandeled command" + cmd.Command)
		}
	}
}
