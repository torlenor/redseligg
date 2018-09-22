package discord

import (
	"time"

	"events"

	"github.com/mitchellh/mapstructure"
)

func (b Bot) getMessageType(mc messageCreate) events.MessageType {
	if val, ok := b.knownChannels[mc.ChannelID]; ok {
		if len(val.Recipients) == 1 {
			return events.WHISPER
		}
	}
	return events.MESSAGE
}

func (b Bot) dispatchMessage(newMessageCreate messageCreate) {
	receiveMessage := events.ReceiveMessage{Type: b.getMessageType(newMessageCreate), Ident: newMessageCreate.Author.ID, Content: newMessageCreate.Content}
	select {
	case b.receiveMessageChan <- receiveMessage:
	default:
	}
}

func (b Bot) handleMessageCreate(data map[string]interface{}) {
	// log.Printf("MESSAGE_CREATE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newMessageCreate messageCreate
	err := mapstructure.Decode(data["d"], &newMessageCreate)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: MESSAGE_CREATE", err)
	}
	// FIXME: Workaround for ChannelID not decoded correctly
	if str, ok := data["d"].(map[string]interface{})["channel_id"].(string); ok {
		newMessageCreate.ChannelID = str
	}
	if str, ok := data["d"].(map[string]interface{})["timestamp"].(string); ok {
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			log.Println("DiscordBot: UNHANDELED ERROR: TYPING_START", err)
		}
		newMessageCreate.Timestamp = t
	}
	log.Printf("DiscordBot: Received: MESSAGE_CREATE from User = %s, Content = %s, Timestamp = %s, ChannelID = %s", newMessageCreate.Author.Username, newMessageCreate.Content, newMessageCreate.Timestamp, newMessageCreate.ChannelID)

	snowflakeID := newMessageCreate.Author.ID

	if snowflakeID != b.ownSnowflakeID {
		b.dispatchMessage(newMessageCreate)
	}
}

func (b *Bot) handleReady(data map[string]interface{}) {
	// log.Printf("READY: data['t']: %s, data['d'] %s", data["t"], data["d"])
	b.ownSnowflakeID = extractOwnSnowflakeID(data)

	var newReady ready
	err := mapstructure.Decode(data["d"], &newReady)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: READY", err)
	}
	// log.Println("READY: ", newReady.toString())
	log.Printf("DiscordBot: Received: READY for Bot User = %s, UserID = %s, SnowflakeID = %s", newReady.User.Username, newReady.User.ID, b.ownSnowflakeID)
}

func (b *Bot) handleGuildCreate(data map[string]interface{}) {
	// log.Printf("GUILD_CREATE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newGuildCreate guildCreate
	err := mapstructure.Decode(data["d"], &newGuildCreate)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: GUILD_CREATE", err)
		return
	}

	newGuild := guild{}
	newGuild.channel = newGuildCreate.Channels
	newGuild.memberCount = newGuildCreate.MemberCount
	newGuild.name = newGuildCreate.Name
	newGuild.snowflakeID = newGuildCreate.ID

	b.guilds[newGuild.name] = newGuild
	b.guildNameToID[newGuild.name] = newGuild.snowflakeID

	log.Println("GUILD_CREATE: Added new Guild:", newGuild.name)
}

func handlePresenceUpdate(data map[string]interface{}) {
	// log.Printf("PRESENCE_UPDATE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newPresenceUpdate presenceUpdate
	err := mapstructure.Decode(data["d"], &newPresenceUpdate)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: PRESENCE_UPDATE", err)
	}
	// FIXME: Workaround for GuildID not decoded correctly
	newPresenceUpdate.GuildID = data["d"].(map[string]interface{})["guild_id"].(string)
	// log.Println("PRESENCE_UPDATE: ", newPresenceUpdate.toString())
	log.Printf("DiscordBot: Received: PRESENCE_UPDATE for UserID = %s", newPresenceUpdate.User.ID)
}

func handlePresenceReplace(data map[string]interface{}) {
	log.Printf("DiscordBot: PRESENCE_REPLACE: data['t']: %s, data['d']: %s", data["t"], data["d"])
}

func handleTypingStart(data map[string]interface{}) {
	// log.Printf("TYPING_START: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newTypingStart typingStart
	err := mapstructure.Decode(data["d"], &newTypingStart)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: TYPING_START", err)
	}
	// FIXME: Workaround for GuildID not decoded correctly
	if str, ok := data["d"].(map[string]interface{})["user_id"].(string); ok {
		newTypingStart.UserID = str
	}
	// FIXME: Workaround for GuildID not decoded correctly
	if str, ok := data["d"].(map[string]interface{})["guild_id"].(string); ok {
		newTypingStart.GuildID = str
	}
	// FIXME: Workaround for ChannelID not decoded correctly
	if str, ok := data["d"].(map[string]interface{})["channel_id"].(string); ok {
		newTypingStart.ChannelID = str
	}
	// FIXME: Workaround for JoinedAt not decoded correctly
	if val, ok := data["d"].(map[string]interface{})["member"]; ok {
		if str, ok := val.(map[string]interface{})["join_at"].(string); ok {
			t, err := time.Parse(time.RFC3339, str)
			if err != nil {
				log.Println("DiscordBot: UNHANDELED ERROR: TYPING_START", err)
			}
			newTypingStart.Member.JoinedAt = t
		}
	}
	// log.Println("TYPING_START: ", newTypingStart.toString())
	log.Println("DiscordBot: Received: TYPING_START User = " + newTypingStart.Member.User.Username)
}

func (b Bot) handleChannelCreate(data map[string]interface{}) {
	// log.Printf("CHANNEL_CREATE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newChannelCreate channelCreate
	err := mapstructure.Decode(data["d"], &newChannelCreate)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: CHANNEL_CREATE", err)
	}
	// log.Println("CHANNEL_CREATE: ", newChannelCreate.toString())
	log.Printf("DiscordBot: Received: CHANNEL_CREATE with ID = %s", newChannelCreate.ID)
	b.knownChannels[newChannelCreate.ID] = newChannelCreate
}

func handleMessageReactionAdd(data map[string]interface{}) {
	// log.Printf("MESSAGE_REACTION_ADD: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newMessageReactionAdd messageReactionAdd
	err := mapstructure.Decode(data["d"], &newMessageReactionAdd)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: MESSAGE_REACTION_ADD", err)
	}
	log.Print("DiscordBot: Received: MESSAGE_REACTION_ADD")
}

func handleMessageReactionRemove(data map[string]interface{}) {
	// log.Printf("MESSAGE_REACTION_REMOVE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newMessageReactionRemove messageReactionRemove
	err := mapstructure.Decode(data["d"], &newMessageReactionRemove)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: MESSAGE_REACTION_REMOVE", err)
	}
	log.Print("DiscordBot: Received: MESSAGE_REACTION_REMOVE")
}

func handleMessageDelete(data map[string]interface{}) {
	// log.Printf("MESSAGE_DELETE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newMessageDelete messageDelete
	err := mapstructure.Decode(data["d"], &newMessageDelete)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: MESSAGE_DELETE", err)
	}
	log.Print("DiscordBot: Received: MESSAGE_DELETE")
}

func handleMessageUpdate(data map[string]interface{}) {
	// log.Printf("MESSAGE_UPDATE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newMessageUpdate messageUpdate
	err := mapstructure.Decode(data["d"], &newMessageUpdate)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: MESSAGE_UPDATE", err)
	}
	log.Print("DiscordBot: Received: MESSAGE_UPDATE")
}

func handleCHannelPinsUpdate(data map[string]interface{}) {
	// log.Printf("CHANNEL_PINS_UPDATE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newChannelPinsUpdate channelPinsUpdate
	err := mapstructure.Decode(data["d"], &newChannelPinsUpdate)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: CHANNEL_PINS_UPDATE", err)
	}
	log.Print("DiscordBot: Received: CHANNEL_PINS_UPDATE")
}

func handleGuildMemberUpdate(data map[string]interface{}) {
	// log.Printf("GUILD_MEMBER_UPDATE: data['t']: %s, data['d']: %s", data["t"], data["d"])
	var newGuildMemberUpdate guildMemberUpdate
	err := mapstructure.Decode(data["d"], &newGuildMemberUpdate)
	if err != nil {
		log.Println("DiscordBot: UNHANDELED ERROR: GUILD_MEMBER_UPDATE", err)
	}
	log.Print("DiscordBot: Received: GUILD_MEMBER_UPDATE")
}

func handleUnknown(data map[string]interface{}) {
	log.Printf("DiscordBot: TODO: %s", data["t"])
}
