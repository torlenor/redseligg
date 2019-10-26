package model

// Post is a event of either incoming or outgoing messages
type Post struct {
	ChannelID string // ChannelID is a unique ID on which the Bot identifies the channel on which this message was seen
	Channel   string // Channel is a clear text name of a Channel in which the message was seen
	UserID    string // UserID is a unique ID on which the Bot identifies the User which sent the message
	User      string // User is a clear text user name of the sender of the message

	Content string

	IsPrivate bool // IsPrivate indicates it is a whisper or similar (depending on the Bot)
}
