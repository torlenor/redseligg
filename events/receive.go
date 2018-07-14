package events

// ReceiveMessage is used to notify about received messages
type ReceiveMessage struct {
	Type    MessageType
	Ident   string
	Content string
}
