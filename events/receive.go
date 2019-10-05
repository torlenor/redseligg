package events

// ReceiveMessage is used to notify about received messages
type ReceiveMessage struct {
	Type MessageType

	ChannelID string // ChannelID is a unique ID on which the Bot identifies the channel on which this message was seen
	Channel   string // Channel is a clear text name of a Channel in which the message was seen
	UserID    string // UserID is a unique ID on which the Bot identifies the User which sent the message
	User      string // User is a clear text user name of the sender of the message

	Content string
}
