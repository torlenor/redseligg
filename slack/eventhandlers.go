package slack

import (
	"encoding/json"
	"regexp"

	"github.com/torlenor/abylebotter/events"
)

func cleanupMessage(msg string) string {
	var re = regexp.MustCompile(`(?m)<(https?://(\w|\.)*)(>|.*>)`)
	return re.ReplaceAllString(msg, "$1")
}

func (b *Bot) handleEventMessage(data []byte) {
	var message EventMessage

	if err := json.Unmarshal(data, &message); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	if message.Subtype != "message_deleted" {
		receiveMessage := events.ReceiveMessage{Type: events.MESSAGE, UserID: message.User, ChannelID: message.Channel, Content: cleanupMessage(message.Text)}
		b.plugins.Send(receiveMessage)
	} else {
		b.log.Debugf("Received message::message_deleted event on Channel ID %s", message.Channel)
	}
}

func (b *Bot) handleEventUserTyping(data []byte) {
	var userTyping EventUserTyping

	if err := json.Unmarshal(data, &userTyping); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received UserTyping event from User ID %s on Channel ID %s", userTyping.User, userTyping.Channel)
}

func (b *Bot) handleEventDesktopNotification(data []byte) {
	var desktopNotification EventDesktopNotification

	if err := json.Unmarshal(data, &desktopNotification); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received DesktopNotification event from Channel ID %s", desktopNotification.Channel)
}

func (b *Bot) handleEventChannelCreated(data []byte) {
	var channelCreated EventChannelCreated

	if err := json.Unmarshal(data, &channelCreated); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelCreated event from Channel Name %s", channelCreated.Channel.Name)
}

func (b *Bot) handleEventChannelJoined(data []byte) {
	var channelJoined EventChannelJoined

	if err := json.Unmarshal(data, &channelJoined); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelJoined event from Channel Name %s", channelJoined.Channel.Name)
}

func (b *Bot) handleEventChannelLeft(data []byte) {
	var channelLeft EventChannelLeft

	if err := json.Unmarshal(data, &channelLeft); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelLeft event from Channel ID %s", channelLeft.Channel)
}

func (b *Bot) handleEventChannelDeleted(data []byte) {
	var channelDeleted EventChannelDeleted

	if err := json.Unmarshal(data, &channelDeleted); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelDeleted event for Channel ID %s", channelDeleted.Channel)
}

func (b *Bot) handleEventMemberJoinedChannel(data []byte) {
	var memberJoinedChannel EventMemberJoinedChannel

	if err := json.Unmarshal(data, &memberJoinedChannel); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received MemberJoinedChannel event, user ID %s -> Channel ID %s", memberJoinedChannel.User, memberJoinedChannel.Channel)
}

func (b *Bot) handleEventGroupJoined(data []byte) {
	var groupJoined EventGroupJoined

	if err := json.Unmarshal(data, &groupJoined); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received GroupJoined event from Channel Name %s", groupJoined.Channel.Name)
}
