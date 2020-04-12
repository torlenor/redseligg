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
	case "reaction_added":
		fallthrough
	case "reaction_removed":
		b.handleEventReactionAddedOrRemoved(message)
	default:
		b.log.Warnf("Received unhandled event %s: %s", event, message)
	}
}

func (b *Bot) handleEventMessage(data []byte) {
	var message eventMessage

	if err := json.Unmarshal(data, &message); err != nil {
		b.log.Errorln("Unable to handle Event Message, error unmarshalling JSON: ", err)
		return
	}

	if message.Subtype != "message_deleted" {
		user, err := b.users.getUserByID(message.User)
		if err != nil {
			b.log.Warnf("Was not able to determine User from message. User ID %s, error: %s", message.User, err)
		}
		receiveMessage := model.Post{User: model.User{ID: message.User, Name: user.Name}, ChannelID: message.Channel, Content: cleanupMessage(message.Text)}
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
		b.log.Errorln("Unable to handle Event UserTyping, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received UserTyping event from User ID %s on Channel ID %s", userTyping.User, userTyping.Channel)
}

func (b *Bot) handleEventDesktopNotification(data []byte) {
	var desktopNotification eventDesktopNotification

	if err := json.Unmarshal(data, &desktopNotification); err != nil {
		b.log.Errorln("Unable to handle Event DesktopNotification, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received DesktopNotification event from Channel ID %s", desktopNotification.Channel)
}

func (b *Bot) handleEventChannelCreated(data []byte) {
	var channelCreated eventChannelCreated

	if err := json.Unmarshal(data, &channelCreated); err != nil {
		b.log.Errorln("Unable to handle Event ChannelCreated, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received ChannelCreated event from Channel Name %s", channelCreated.Channel.Name)
}

func (b *Bot) handleEventChannelJoined(data []byte) {
	var channelJoined eventChannelJoined

	if err := json.Unmarshal(data, &channelJoined); err != nil {
		b.log.Errorln("Unable to handle Event ChannelJoined, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received ChannelJoined event from Channel Name %s", channelJoined.Channel.Name)
}

func (b *Bot) handleEventChannelLeft(data []byte) {
	var channelLeft eventChannelLeft

	if err := json.Unmarshal(data, &channelLeft); err != nil {
		b.log.Errorln("Unable to handle Event ChannelLeft, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received ChannelLeft event from Channel ID %s", channelLeft.Channel)
}

func (b *Bot) handleEventChannelDeleted(data []byte) {
	var channelDeleted eventChannelDeleted

	if err := json.Unmarshal(data, &channelDeleted); err != nil {
		b.log.Errorln("Unable to handle Event ChannelDeleted, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received ChannelDeleted event for Channel ID %s", channelDeleted.Channel)
}

func (b *Bot) handleEventMemberJoinedChannel(data []byte) {
	var memberJoinedChannel eventMemberJoinedChannel

	if err := json.Unmarshal(data, &memberJoinedChannel); err != nil {
		b.log.Errorln("Unable to handle Event ChannelJoined, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received MemberJoinedChannel event, user ID %s -> Channel ID %s", memberJoinedChannel.User, memberJoinedChannel.Channel)
}

func (b *Bot) handleEventGroupJoined(data []byte) {
	var groupJoined eventGroupJoined

	if err := json.Unmarshal(data, &groupJoined); err != nil {
		b.log.Errorln("Unable to handle Event GroupJoined, error unmarshalling JSON: ", err)
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
		b.log.Errorln("Unable to handle Event EventUserChanged, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received UserChanged event for User %s", userChanged.User.Name)
}

func (b *Bot) handleEventTeamJoin(data []byte) {
	var teamJoin eventUser

	if err := json.Unmarshal(data, &teamJoin); err != nil {
		b.log.Errorln("Unable to handle Event TeamJoin, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received TeamJoin event for User %s", teamJoin.User.Name)
}

func (b *Bot) handleEventDnDUpdatedUser(data []byte) {
	var dndUpdatedUser eventDnDUpdatedUser

	if err := json.Unmarshal(data, &dndUpdatedUser); err != nil {
		b.log.Errorln("Unable to handle Event DnDUpdatedUser, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received DnDUpdatedUser event for User ID %s", dndUpdatedUser.User)
}

func (b *Bot) handleEventIMCreated(data []byte) {
	var imCreated eventIMCreated

	if err := json.Unmarshal(data, &imCreated); err != nil {
		b.log.Errorln("Unable to handle Event IMCreated, error unmarshalling JSON: ", err)
		return
	}

	b.log.Debugf("Received IMCreated event for User ID %s", imCreated.User)
}

func (b *Bot) handleEventReactionAddedOrRemoved(data []byte) {
	var reaction eventReactionAddedOrRemoved

	if err := json.Unmarshal(data, &reaction); err != nil {
		b.log.Errorln("Unable to handle Event Reaction Added/Removed, error unmarshalling JSON: ", err)
		return
	}

	var reactionType string
	switch reaction.Type {
	case "reaction_added":
		// example: {"type":"reaction_added","user":"UNL92ERS4","item":{"type":"message","channel":"G011C8YPGET","ts":"1586690851.001000"},"reaction":"wink","item_user":"UNL92ERS4","event_ts":"1586690859.001100","ts":"1586690859.001100"}
		reactionType = "added"
	case "reaction_removed":
		// example: {"type":"reaction_removed","user":"UNL92ERS4","item":{"type":"message","channel":"G011C8YPGET","ts":"1586690851.001000"},"reaction":"wink","item_user":"UNL92ERS4","event_ts":"1586691109.001200","ts":"1586691109.001200"}
		reactionType = "removed"
	default:
		b.log.Warnf("Received unknown Event Reaction of type %s on Channel ID %s", reaction.Type, reaction.Item.Channel)
		return
	}

	forPlugin := model.Reaction{
		Message: model.MessageIdentifier{
			ID: reaction.Item.Ts, Channel: reaction.Item.Channel,
		},

		Type:     reactionType,
		Reaction: reaction.Reaction,
		User:     model.User{ID: reaction.User},
	}

	for _, plugin := range b.plugins {
		plugin.OnReactionAdded(forPlugin)
	}
}
