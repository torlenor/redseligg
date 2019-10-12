package slack

import "github.com/torlenor/abylebotter/events"

func (b *Bot) pluginMessageReceiver() {
	for sendMsg := range b.plugins.SendChannel() {
		switch sendMsg.Type {
		case events.MESSAGE:
			err := b.sendMessage(sendMsg.ChannelID, sendMsg.Content)
			if err != nil {
				b.log.Errorln("Error sending message:", err)
			}
		case events.WHISPER:
			var userID string
			if len(sendMsg.UserID) != 0 {
				userID = sendMsg.UserID
			} else if len(sendMsg.User) != 0 {
				var err error
				userID, err = b.users.getUserNameByID(sendMsg.User)
				if err != nil {
					b.log.Errorf("User not found, not sending Whisper")
				}
			} else {
				b.log.Errorf("Plugin did not provide User or UserID, not sending Whisper")
			}
			err := b.sendWhisper(userID, sendMsg.Content)
			if err != nil {
				b.log.Errorln("Error sending whisper:", err)
			}
		default:
			b.log.Warnf("Bot does not support Send Event %s", sendMsg.Type)
		}
	}
}
