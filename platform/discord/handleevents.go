package discord

import (
	"github.com/torlenor/abylebotter/model"
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
	receiveMessage = model.Post{User: model.User{ID: msg.Author.ID, Name: combineUsernameAndDiscriminator(msg.Author.Username, msg.Author.Discriminator)}, ChannelID: msg.ChannelID, Content: msg.Content}
	if b.getMessageType(msg) == WHISPER {
		receiveMessage.IsPrivate = true
	}

	for _, plugin := range b.plugins {
		plugin.OnPost(receiveMessage)
	}
}

func (b *Bot) handleMessageCreate(data map[string]interface{}) {
	newMessageCreate, err := decodeMessageCreate(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_CREATE", err)
		return
	}

	log.Debugf("Received: MESSAGE_CREATE from User = %s, Content = %s, Timestamp = %s, ChannelID = %s", newMessageCreate.Author.Username, newMessageCreate.Content, newMessageCreate.Timestamp, newMessageCreate.ChannelID)

	snowflakeID := newMessageCreate.Author.ID

	if snowflakeID != b.ownSnowflakeID {
		b.dispatchMessage(newMessageCreate)
	}
}

func (b *Bot) handleReady(data map[string]interface{}) {
	newReady, err := decodeReady(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: READY", err)
		return
	}
	b.ownSnowflakeID = newReady.User.ID

	log.Debugf("Received: READY for Bot User = %s, UserID = %s, SnowflakeID = %s", newReady.User.Username, newReady.User.ID, b.ownSnowflakeID)
}

func (b *Bot) handleGuildCreate(data map[string]interface{}) {
	newGuildCreate, err := decodeGuildCreate(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: GUILD_CREATE", err)
		return
	}

	b.guilds[newGuildCreate.ID] = newGuildCreate
	b.guildNameToID[newGuildCreate.Name] = newGuildCreate.ID

	log.Debugln("GUILD_CREATE: Added new Guild:", newGuildCreate.Name)
}

func (b *Bot) handlePresenceUpdate(data map[string]interface{}) {
	newPresenceUpdate, err := decodePresenceUpdate(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: PRESENCE_UPDATE", err)
		return
	}

	log.Debugf("Received: PRESENCE_UPDATE for UserID = %s", newPresenceUpdate.User.ID)
}

func (b *Bot) handlePresenceReplace(data map[string]interface{}) {
	log.Warnf("NOT_IMPLEMENTED: PRESENCE_REPLACE: data['t']: %s, data['d']: %s", data["t"], data["d"])
}

func (b *Bot) handleTypingStart(data map[string]interface{}) {
	newTypingStart, err := decodeTypingStart(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: TYPING_START", err)
		return
	}

	log.Debugf("Received: TYPING_START User = %s", newTypingStart.Member.User.Username)
}

func (b *Bot) addKnownChannel(channel channelCreate) {
	b.knownChannels[channel.ID] = channel
}

func (b *Bot) handleChannelCreate(data map[string]interface{}) {
	newChannelCreate, err := decodeChannelCreate(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: CHANNEL_CREATE", err)
		return
	}

	log.Debugf("Received: CHANNEL_CREATE with ID = %s", newChannelCreate.ID)
	b.addKnownChannel(newChannelCreate)
}

func (b *Bot) handleMessageReactionAdd(data map[string]interface{}) {
	newMessageReactionAdd, err := decodeMessageReactionAdd(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_REACTION_ADD", err)
		return
	}

	emoji, err := getAbyleBotterEmojiFromDiscordEmoji(newMessageReactionAdd.Emoji.Name)

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
}

func (b *Bot) handleMessageReactionRemove(data map[string]interface{}) {
	newMessageReactionRemove, err := decodeMessageReactionRemove(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_REACTION_REMOVE", err)
		return
	}

	emoji, err := getAbyleBotterEmojiFromDiscordEmoji(newMessageReactionRemove.Emoji.Name)

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
}

func (b *Bot) handleMessageDelete(data map[string]interface{}) {
	newMessageReactionDelete, err := decodeMessageDelete(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_DELETE", err)
		return
	}

	log.Debugln("Received: MESSAGE_DELETE", newMessageReactionDelete)
}

func (b *Bot) handleMessageUpdate(data map[string]interface{}) {
	newMessageUpdate, err := decodeMessageUpdate(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: MESSAGE_UPDATE", err)
		return
	}

	log.Debugln("Received: MESSAGE_UPDATE", newMessageUpdate)
}

func (b *Bot) handleChannelPinsUpdate(data map[string]interface{}) {
	newChannelPinsUpdate, err := decodeChannelPinsUpdate(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: CHANNEL_PINS_UPDATE", err)
		return
	}

	log.Debugln("Received: CHANNEL_PINS_UPDATE", newChannelPinsUpdate)
}

func (b *Bot) handleGuildMemberUpdate(data map[string]interface{}) {
	newGuildMemberUpdate, err := decodeGuildMemberUpdate(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: GUILD_MEMBER_UPDATE", err)
		return
	}

	log.Debugln("Received: GUILD_MEMBER_UPDATE", newGuildMemberUpdate)
}

func (b *Bot) handlePresencesReplace(data map[string]interface{}) {
	newPresencesReplace, err := decodePresencesReplace(data)
	if err != nil {
		log.Errorln("UNHANDLED ERROR: PRESENCES_REPLACE", err)
		return
	}

	log.Debugln("Received: PRESENCES_REPLACE", newPresencesReplace)
}

func (b *Bot) handleUnknown(data map[string]interface{}) {
	log.Debugf("TODO HANDLE UNKNOWN EVENT: %s", data["t"])
}
