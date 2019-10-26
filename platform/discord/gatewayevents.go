package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type typingStart struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	UserID    string `json:"user_id"`
	Timestamp int    `json:"timestamp"`
	Member    struct {
		User struct {
			Username      string `json:"username"`
			ID            string `json:"id"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		} `json:"user"`
		Roles    []string  `json:"roles"`
		Mute     bool      `json:"mute"`
		JoinedAt time.Time `json:"joined_at"`
		Deaf     bool      `json:"deaf"`
	} `json:"member"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	// } `json:"d"`
}

func decodeTypingStart(data map[string]interface{}) (typingStart, error) {
	var newTypingStart typingStart
	err := mapstructure.Decode(data["d"], &newTypingStart)
	if err != nil {
		log.Debugln("Error automatically decoding typingStart, ignoring it")
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
				log.Errorln("UNHANDELED ERROR: decodeTypingStart:", err)
			}
			newTypingStart.Member.JoinedAt = t
		}
	}

	return newTypingStart, nil
}

func (t typingStart) toString() string {
	data, err := json.Marshal(t)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type presenceUpdate struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	Status  string      `json:"status"`
	Roles   []string    `json:"roles"`
	Nick    interface{} `json:"nick"`
	GuildID string      `json:"guild_id"`
	Game    struct {
		Type int    `json:"type"`
		Name string `json:"name"`
	} `json:"game"`
	// } `json:"d"`
}

func decodePresenceUpdate(data map[string]interface{}) (presenceUpdate, error) {
	var newPresenceUpdate presenceUpdate
	err := mapstructure.Decode(data["d"], &newPresenceUpdate)
	if err != nil {
		log.Debugln("Error automatically decoding presenceUpdate, ignoring it")
	}

	// FIXME: Workaround for GuildID not decoded correctly
	newPresenceUpdate.GuildID = data["d"].(map[string]interface{})["guild_id"].(string)

	return newPresenceUpdate, nil
}

func (p presenceUpdate) toString() string {
	data, err := json.Marshal(p)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageCreate struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	Type            int           `json:"type"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	Pinned          bool          `json:"pinned"`
	Nonce           interface{}   `json:"nonce"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	MentionEveryone bool          `json:"mention_everyone"`
	Member          struct {
		Roles    []string  `json:"roles"`
		Mute     bool      `json:"mute"`
		JoinedAt time.Time `json:"joined_at"`
		Deaf     bool      `json:"deaf"`
	} `json:"member"`
	ID              string        `json:"id"`
	Embeds          []interface{} `json:"embeds"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Content         string        `json:"content"`
	ChannelID       string        `json:"channel_id"`
	Author          struct {
		Username      string      `json:"username"`
		ID            string      `json:"id"`
		Discriminator string      `json:"discriminator"`
		Bot           bool        `json:"bot"`
		Avatar        interface{} `json:"avatar"`
	} `json:"author"`
	Attachments []interface{} `json:"attachments"`
	GuildID     string        `json:"guild_id"`
	// } `json:"d"`
}

func decodeMessageCreate(data map[string]interface{}) (messageCreate, error) {
	var newMessageCreate messageCreate
	err := mapstructure.Decode(data["d"], &newMessageCreate)
	if err != nil {
		log.Debugln("Error automatically decoding messageCreate, ignoring it")
	}

	// FIXME: Workaround for ChannelID not decoded correctly
	if str, ok := data["d"].(map[string]interface{})["channel_id"].(string); ok {
		newMessageCreate.ChannelID = str
	}
	if str, ok := data["d"].(map[string]interface{})["timestamp"].(string); ok {
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return messageCreate{}, errors.New("decodeMessageCreate:" + err.Error())
		}
		newMessageCreate.Timestamp = t
	}

	return newMessageCreate, nil
}

func (mc messageCreate) toString() string {
	data, err := json.Marshal(mc)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type channel struct {
	Type                 int           `json:"type"`
	Topic                string        `json:"topic,omitempty"`
	Position             int           `json:"position"`
	PermissionOverwrites []interface{} `json:"permission_overwrites"`
	Name                 string        `json:"name"`
	LastPinTimestamp     time.Time     `json:"last_pin_timestamp,omitempty"`
	LastMessageID        string        `json:"last_message_id,omitempty"`
	ID                   string        `json:"id"`
	UserLimit            int           `json:"user_limit,omitempty"`
	ParentID             string        `json:"parent_id,omitempty"`
	Bitrate              int           `json:"bitrate,omitempty"`
}

type guildCreate struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	VoiceStates       []interface{} `json:"voice_states"`
	VerificationLevel int           `json:"verification_level"`
	Unavailable       bool          `json:"unavailable"`
	SystemChannelID   interface{}   `json:"system_channel_id"`
	Splash            interface{}   `json:"splash"`
	Roles             []struct {
		Position    int    `json:"position"`
		Permissions int    `json:"permissions"`
		Name        string `json:"name"`
		Mentionable bool   `json:"mentionable"`
		Managed     bool   `json:"managed"`
		ID          string `json:"id"`
		Hoist       bool   `json:"hoist"`
		Color       int    `json:"color"`
	} `json:"roles"`
	Region    string `json:"region"`
	Presences []struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		Status string      `json:"status"`
		Game   interface{} `json:"game"`
	} `json:"presences"`
	OwnerID  string `json:"owner_id"`
	Name     string `json:"name"`
	MfaLevel int    `json:"mfa_level"`
	Members  []struct {
		User struct {
			Username      string `json:"username"`
			ID            string `json:"id"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		} `json:"user"`
		Roles    []string  `json:"roles"`
		Mute     bool      `json:"mute"`
		JoinedAt time.Time `json:"joined_at"`
		Deaf     bool      `json:"deaf"`
		Nick     string    `json:"nick,omitempty"`
	} `json:"members"`
	MemberCount                 int           `json:"member_count"`
	Lazy                        bool          `json:"lazy"`
	Large                       bool          `json:"large"`
	JoinedAt                    time.Time     `json:"joined_at"`
	ID                          string        `json:"id"`
	Icon                        string        `json:"icon"`
	Features                    []interface{} `json:"features"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	Emojis                      []interface{} `json:"emojis"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	Channels                    []channel     `json:"channels"`
	ApplicationID               interface{}   `json:"application_id"`
	AfkTimeout                  int           `json:"afk_timeout"`
	AfkChannelID                interface{}   `json:"afk_channel_id"`
	// } `json:"d"`
}

func decodeGuildCreate(data map[string]interface{}) (guildCreate, error) {
	var newGuildCreate guildCreate
	err := mapstructure.Decode(data["d"], &newGuildCreate)
	if err != nil {
		return guildCreate{}, errors.New("decodeGuildCreate" + err.Error())
	}

	return newGuildCreate, nil
}

func (gc guildCreate) toString() string {
	data, err := json.Marshal(gc)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type ready struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	V            int `json:"v"`
	UserSettings struct {
	} `json:"user_settings"`
	User struct {
		Verified      bool        `json:"verified"`
		Username      string      `json:"username"`
		MfaEnabled    bool        `json:"mfa_enabled"`
		ID            string      `json:"id"`
		Email         interface{} `json:"email"`
		Discriminator string      `json:"discriminator"`
		Bot           bool        `json:"bot"`
		Avatar        interface{} `json:"avatar"`
	} `json:"user"`
	SessionID       string        `json:"session_id"`
	Relationships   []interface{} `json:"relationships"`
	PrivateChannels []interface{} `json:"private_channels"`
	Presences       []interface{} `json:"presences"`
	Guilds          []struct {
		Unavailable bool   `json:"unavailable"`
		ID          string `json:"id"`
	} `json:"guilds"`
	Trace []string `json:"_trace"`
	// } `json:"d"`
}

func decodeReady(data map[string]interface{}) (ready, error) {
	var newReady ready
	err := mapstructure.Decode(data["d"], &newReady)
	if err != nil {
		return ready{}, errors.New("decodeReady" + err.Error())
	}

	return newReady, nil
}

func (r ready) toString() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type channelCreate struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	Type       int `json:"type"`
	Recipients []struct {
		Username      string `json:"username"`
		ID            string `json:"id"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
	} `json:"recipients"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
	// } `json:"d"`
}

func decodeChannelCreate(data map[string]interface{}) (channelCreate, error) {
	var newChannelCreate channelCreate
	err := mapstructure.Decode(data["d"], &newChannelCreate)
	if err != nil {
		return channelCreate{}, errors.New("decodeChannelCreate: " + err.Error())
	}

	return newChannelCreate, nil
}

func (cc channelCreate) toString() string {
	data, err := json.Marshal(cc)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageDelete struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	// } `json:"d"`
}

func decodeMessageDelete(data map[string]interface{}) (messageDelete, error) {
	var newMessageDelete messageDelete
	err := mapstructure.Decode(data["d"], &newMessageDelete)
	if err != nil {
		return messageDelete{}, errors.New("decodeMessageDelete: " + err.Error())
	}

	return newMessageDelete, nil
}

func (md messageDelete) toString() string {
	data, err := json.Marshal(md)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageReactionAdd struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     struct {
		Name     string      `json:"name"`
		ID       interface{} `json:"id"`
		Animated bool        `json:"animated"`
	} `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	// } `json:"d"`
}

func decodeMessageReactionAdd(data map[string]interface{}) (messageReactionAdd, error) {
	var newMessageReactionAdd messageReactionAdd
	err := mapstructure.Decode(data["d"], &newMessageReactionAdd)
	if err != nil {
		return messageReactionAdd{}, errors.New("decodeMessageReactionAdd: " + err.Error())
	}

	return newMessageReactionAdd, nil
}

func (mra messageReactionAdd) toString() string {
	data, err := json.Marshal(mra)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageReactionRemove struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     struct {
		Name     string      `json:"name"`
		ID       interface{} `json:"id"`
		Animated bool        `json:"animated"`
	} `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	// } `json:"d"`
}

func decodeMessageReactionRemove(data map[string]interface{}) (messageReactionRemove, error) {
	var newMessageReactionRemove messageReactionRemove
	err := mapstructure.Decode(data["d"], &newMessageReactionRemove)
	if err != nil {
		return messageReactionRemove{}, errors.New("decodeMessageReactionRemove: " + err.Error())
	}

	return newMessageReactionRemove, nil
}

func (mrr messageReactionRemove) toString() string {
	data, err := json.Marshal(mrr)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageUpdate struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	Type            int           `json:"type"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	Pinned          bool          `json:"pinned"`
	Nonce           interface{}   `json:"nonce"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	MentionEveryone bool          `json:"mention_everyone"`
	ID              string        `json:"id"`
	Embeds          []interface{} `json:"embeds"`
	EditedTimestamp time.Time     `json:"edited_timestamp"`
	Content         string        `json:"content"`
	ChannelID       string        `json:"channel_id"`
	Author          struct {
		Username      string `json:"username"`
		ID            string `json:"id"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
	} `json:"author"`
	Attachments []interface{} `json:"attachments"`
	GuildID     string        `json:"guild_id"`
	// } `json:"d"`
}

func decodeMessageUpdate(data map[string]interface{}) (messageUpdate, error) {
	var newMessageUpdate messageUpdate
	err := mapstructure.Decode(data["d"], &newMessageUpdate)
	if err != nil {
		return messageUpdate{}, errors.New("decodeMessageUpdate: " + err.Error())
	}

	return newMessageUpdate, nil
}

func (mu messageUpdate) toString() string {
	data, err := json.Marshal(mu)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type channelPinsUpdate struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	LastPinTimestamp time.Time `json:"last_pin_timestamp"`
	ChannelID        string    `json:"channel_id"`
	GuildID          string    `json:"guild_id"`
	// } `json:"d"`
}

func decodeChannelPinsUpdate(data map[string]interface{}) (channelPinsUpdate, error) {
	var newChannelPinsUpdate channelPinsUpdate
	err := mapstructure.Decode(data["d"], &newChannelPinsUpdate)
	if err != nil {
		return channelPinsUpdate{}, errors.New("decodeChannelPinsUpdate: " + err.Error())
	}

	return newChannelPinsUpdate, nil
}

func (cpu channelPinsUpdate) toString() string {
	data, err := json.Marshal(cpu)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type guildMemberUpdate struct {
	// T  string `json:"t"`
	// S  int    `json:"s"`
	// Op int    `json:"op"`
	// D  struct {
	User struct {
		Username      string      `json:"username"`
		ID            string      `json:"id"`
		Discriminator string      `json:"discriminator"`
		Bot           bool        `json:"bot"`
		Avatar        interface{} `json:"avatar"`
	} `json:"user"`
	Roles   []string    `json:"roles"`
	Nick    interface{} `json:"nick"`
	GuildID string      `json:"guild_id"`
	// } `json:"d"`
}

func decodeGuildMemberUpdate(data map[string]interface{}) (guildMemberUpdate, error) {
	var newGuildMemberUpdate guildMemberUpdate
	err := mapstructure.Decode(data["d"], &newGuildMemberUpdate)
	if err != nil {
		return guildMemberUpdate{}, errors.New("decodeGuildMemberUpdate: " + err.Error())
	}

	return newGuildMemberUpdate, nil
}

func (gmu guildMemberUpdate) toString() string {
	data, err := json.Marshal(gmu)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}

	return fmt.Sprintf("%s", data)
}

type presencesReplace struct {
	// T  string        `json:"t"`
	// S  int           `json:"s"`
	// Op int           `json:"op"`
	D []interface{} `json:"d"`
}

func decodePresencesReplace(data map[string]interface{}) (presencesReplace, error) {
	var newPresencesReplace presencesReplace
	err := mapstructure.Decode(data["d"], &newPresencesReplace)
	if err != nil {
		return presencesReplace{}, errors.New("decodePresencesReplace: " + err.Error())
	}

	return newPresencesReplace, nil
}

func (pr presencesReplace) toString() string {
	data, err := json.Marshal(pr)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}

	return fmt.Sprintf("%s", data)
}
