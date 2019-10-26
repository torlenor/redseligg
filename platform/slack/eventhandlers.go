package slack

import (
	"encoding/json"
	"regexp"

	"github.com/torlenor/abylebotter/model"
)

func cleanupMessage(msg string) string {
	var re = regexp.MustCompile(`(?m)<(https?://(\w|\.)*)(>|.*>)`)
	return re.ReplaceAllString(msg, "$1")
}

func (b *Bot) eventDispatcher(event interface{}, message []byte) {
	switch event {
	case "hello":
	case "message":
		b.handleEventMessage(message)
	case "desktop_notification":
		b.handleEventDesktopNotification(message)
	case "user_typing":
		b.handleEventUserTyping(message)
	case "channel_created":
		b.handleEventChannelCreated(message)
	case "channel_deleted":
		b.handleEventChannelDeleted(message)
	case "channel_joined":
		b.handleEventChannelJoined(message)
	case "channel_left":
		b.handleEventChannelLeft(message)
	case "member_joined_channel":
		b.handleEventMemberJoinedChannel(message)
	case "group_joined":
		b.handleEventGroupJoined(message)
	case "user_changed":
		b.handleEventUserChanged(message)
	case "team_join":
		b.handleEventTeamJoin(message)
	case "dnd_updated_user":
		b.handleEventDnDUpdatedUser(message)
	case "im_created":
		b.handleEventIMCreated(message)
	case "pong":
		b.handleEventPong(message)
	default:
		b.log.Warnf("Received unhandled event %s: %s", event, message)
	}
}

func (b *Bot) handleEventMessage(data []byte) {
	var message eventMessage

	if err := json.Unmarshal(data, &message); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	if message.Subtype != "message_deleted" {
		receiveMessage := model.Post{UserID: message.User, ChannelID: message.Channel, Content: cleanupMessage(message.Text)}
		for _, plugin := range b.plugins {
			plugin.OnPost(receiveMessage)
		}
	} else {
		b.log.Debugf("Received message::message_deleted event on Channel ID %s", message.Channel)
	}
}

func (b *Bot) handleEventUserTyping(data []byte) {
	var userTyping eventUserTyping

	if err := json.Unmarshal(data, &userTyping); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received UserTyping event from User ID %s on Channel ID %s", userTyping.User, userTyping.Channel)
}

func (b *Bot) handleEventDesktopNotification(data []byte) {
	var desktopNotification eventDesktopNotification

	if err := json.Unmarshal(data, &desktopNotification); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received DesktopNotification event from Channel ID %s", desktopNotification.Channel)
}

func (b *Bot) handleEventChannelCreated(data []byte) {
	var channelCreated eventChannelCreated

	if err := json.Unmarshal(data, &channelCreated); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelCreated event from Channel Name %s", channelCreated.Channel.Name)
}

func (b *Bot) handleEventChannelJoined(data []byte) {
	var channelJoined eventChannelJoined

	if err := json.Unmarshal(data, &channelJoined); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelJoined event from Channel Name %s", channelJoined.Channel.Name)
}

func (b *Bot) handleEventChannelLeft(data []byte) {
	var channelLeft eventChannelLeft

	if err := json.Unmarshal(data, &channelLeft); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelLeft event from Channel ID %s", channelLeft.Channel)
}

func (b *Bot) handleEventChannelDeleted(data []byte) {
	var channelDeleted eventChannelDeleted

	if err := json.Unmarshal(data, &channelDeleted); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received ChannelDeleted event for Channel ID %s", channelDeleted.Channel)
}

func (b *Bot) handleEventMemberJoinedChannel(data []byte) {
	var memberJoinedChannel eventMemberJoinedChannel

	if err := json.Unmarshal(data, &memberJoinedChannel); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received MemberJoinedChannel event, user ID %s -> Channel ID %s", memberJoinedChannel.User, memberJoinedChannel.Channel)
}

func (b *Bot) handleEventGroupJoined(data []byte) {
	var groupJoined eventGroupJoined

	if err := json.Unmarshal(data, &groupJoined); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received GroupJoined event from Channel Name %s", groupJoined.Channel.Name)
}

func (b *Bot) handleEventPong(data []byte) {
	b.receivePong(data)
}

func (b *Bot) handleEventUserChanged(data []byte) {
	var userChanged eventUser

	if err := json.Unmarshal(data, &userChanged); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received UserChanged event for User %s", userChanged.User.Name)
}

func (b *Bot) handleEventTeamJoin(data []byte) {
	var teamJoin eventUser

	if err := json.Unmarshal(data, &teamJoin); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received TeamJoin event for User %s", teamJoin.User.Name)
}

func (b *Bot) handleEventDnDUpdatedUser(data []byte) {
	var dndUpdatedUser eventDnDUpdatedUser

	if err := json.Unmarshal(data, &dndUpdatedUser); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received DnDUpdatedUser event for User ID %s", dndUpdatedUser.User)
}

func (b *Bot) handleEventIMCreated(data []byte) {
	var imCreated eventIMCreated

	if err := json.Unmarshal(data, &imCreated); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Debugf("Received IMCreated event for User ID %s", imCreated.User)
}
