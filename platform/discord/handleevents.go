package discord

import (
	"encoding/json"

	"github.com/torlenor/redseligg/model"
)

type messageType int

// MessageTypes
const ( // iota is reset to 0
	UNKNOWN messageType = iota
	WHISPER
	MESSAGE
)

var messageTypes = [...]string{
	"UNKNOWN",
	"WHISPER",
	"MESSAGE",
}

func (messageType messageType) String() string {
	return messageTypes[messageType]
}

func (b *Bot) getMessageType(mc messageCreate) messageType {
	if val, ok := b.knownChannels[mc.ChannelID]; ok {
		if len(val.Recipients) == 1 {
			return WHISPER
		}
	}
	return MESSAGE
}

func (b *Bot) dispatchMessage(msg messageCreate) {
	var receiveMessage model.Post
	receiveMessage = model.Post{ServerID: msg.GuildID, User: model.User{ID: msg.Author.ID, Name: combineUsernameAndDiscriminator(msg.Author.Username, msg.Author.Discriminator)}, ChannelID: msg.ChannelID, Content: msg.Content}
	if b.getMessageType(msg) == WHISPER {
		receiveMessage.IsPrivate = true
	}

	for _, plugin := range b.plugins {
		plugin.OnPost(receiveMessage)
	}

	b.Dispatcher.OnPost(receiveMessage)
}

func (b *Bot) handleMessageCreate(data json.RawMessage) {
	var newMessageCreate messageCreate
	err := json.Unmarshal(data, &newMessageCreate)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_CREATE", err)
		return
	}

	log.Tracef("Received: MESSAGE_CREATE from User = %s, Content = %s, Timestamp = %s, ChannelID = %s", newMessageCreate.Author.Username, newMessageCreate.Content, newMessageCreate.Timestamp, newMessageCreate.ChannelID)
	b.dispatchMessage(newMessageCreate)
}

func (b *Bot) handleReady(data json.RawMessage) {
	var newReady ready
	err := json.Unmarshal(data, &newReady)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: READY", err)
		return
	}
	b.ownSnowflakeID = newReady.User.ID
	b.sessionID = newReady.SessionID

	log.Tracef("Received: READY for Bot User = %s, UserID = %s, SnowflakeID = %s", newReady.User.Username, newReady.User.ID, b.ownSnowflakeID)
}

func (b *Bot) handleGuildCreate(data json.RawMessage) {
	var newGuildCreate guildCreate
	err := json.Unmarshal(data, &newGuildCreate)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: GUILD_CREATE", err)
		return
	}

	b.guilds[newGuildCreate.ID] = newGuildCreate
	b.guildNameToID[newGuildCreate.Name] = newGuildCreate.ID

	log.Traceln("GUILD_CREATE: Added new Guild:", newGuildCreate.Name)
}

func (b *Bot) handlePresenceUpdate(data json.RawMessage) {
	var newPresenceUpdate presenceUpdate
	err := json.Unmarshal(data, &newPresenceUpdate)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: PRESENCE_UPDATE", err)
		return
	}

	log.Tracef("Received: PRESENCE_UPDATE for UserID = %s", newPresenceUpdate.User.ID)
}

func (b *Bot) handlePresenceReplace(data json.RawMessage) {
	log.Warnf("NOT_IMPLEMENTED: PRESENCE_REPLACE")
}

func (b *Bot) handleTypingStart(data json.RawMessage) {
	var newTypingStart typingStart
	err := json.Unmarshal(data, &newTypingStart)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: TYPING_START", err)
		return
	}

	log.Tracef("Received: TYPING_START User = %s", newTypingStart.Member.User.Username)
}

func (b *Bot) addKnownChannel(channel channelCreate) {
	b.knownChannels[channel.ID] = channel
}

func (b *Bot) handleChannelCreate(data json.RawMessage) {
	var newChannelCreate channelCreate
	err := json.Unmarshal(data, &newChannelCreate)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: CHANNEL_CREATE", err)
		return
	}

	log.Tracef("Received: CHANNEL_CREATE with ID = %s", newChannelCreate.ID)

	b.addKnownChannel(newChannelCreate)
}

func (b *Bot) handleMessageReactionAdd(data json.RawMessage) {
	var newMessageReactionAdd messageReactionAdd
	err := json.Unmarshal(data, &newMessageReactionAdd)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_REACTION_ADD", err)
		return
	}

	emoji, err := getRedseliggEmojiFromDiscordEmoji(newMessageReactionAdd.Emoji.Name)
	if err != nil {
		log.Debugf("Could not map emoji %s, consider adding it to the mapping: %s", newMessageReactionAdd.Emoji.Name, err)
	}

	reaction := model.Reaction{
		Message: model.MessageIdentifier{
			ID:      newMessageReactionAdd.MessageID,
			Channel: newMessageReactionAdd.ChannelID,
		},
		Type:     "added",
		Reaction: emoji,
		User:     model.User{ID: newMessageReactionAdd.UserID},
	}

	for _, plugin := range b.plugins {
		plugin.OnReactionAdded(reaction)
	}

	log.Traceln("Received: MESSAGE_REACTION_ADD", newMessageReactionAdd)
}

func (b *Bot) handleMessageReactionRemove(data json.RawMessage) {
	var newMessageReactionRemove messageReactionRemove
	err := json.Unmarshal(data, &newMessageReactionRemove)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_REACTION_REMOVE", err)
		return
	}

	emoji, err := getRedseliggEmojiFromDiscordEmoji(newMessageReactionRemove.Emoji.Name)
	if err != nil {
		log.Debugf("Could not map emoji %s, consider adding it to the mapping: %s", newMessageReactionRemove.Emoji.Name, err)
	}

	reaction := model.Reaction{
		Message: model.MessageIdentifier{
			ID:      newMessageReactionRemove.MessageID,
			Channel: newMessageReactionRemove.ChannelID,
		},
		Type:     "removed",
		Reaction: emoji,
		User:     model.User{ID: newMessageReactionRemove.UserID},
	}

	for _, plugin := range b.plugins {
		plugin.OnReactionRemoved(reaction)
	}

	log.Traceln("Received: MESSAGE_REACTION_REMOVE", newMessageReactionRemove)
}

func (b *Bot) handleMessageDelete(data json.RawMessage) {
	var newMessageDelete messageDelete
	err := json.Unmarshal(data, &newMessageDelete)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_DELETE", err)
		return
	}

	log.Traceln("Received: MESSAGE_DELETE", newMessageDelete)
}

func (b *Bot) handleMessageUpdate(data json.RawMessage) {
	var newMessageUpdate messageUpdate
	err := json.Unmarshal(data, &newMessageUpdate)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_UPDATE", err)
		return
	}

	log.Traceln("Received: MESSAGE_UPDATE", newMessageUpdate)
}

func (b *Bot) handleChannelPinsUpdate(data json.RawMessage) {
	var newChannelPinsUpdate channelPinsUpdate
	err := json.Unmarshal(data, &newChannelPinsUpdate)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: CHANNEL_PINS_UPDATE", err)
		return
	}

	log.Traceln("Received: CHANNEL_PINS_UPDATE", newChannelPinsUpdate)
}

func (b *Bot) handleGuildMemberUpdate(data json.RawMessage) {
	var newGuildMemberUpdate guildMemberUpdate
	err := json.Unmarshal(data, &newGuildMemberUpdate)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: GUILD_MEMBER_UPDATE", err)
		return
	}

	log.Traceln("Received: GUILD_MEMBER_UPDATE", newGuildMemberUpdate)
}

func (b *Bot) handlePresencesReplace(data json.RawMessage) {
	var newPresencesReplace presenceUpdate
	err := json.Unmarshal(data, &newPresencesReplace)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: PRESENCES_REPLACE", err)
		return
	}

	log.Traceln("Received: PRESENCES_REPLACE", newPresencesReplace)
}
