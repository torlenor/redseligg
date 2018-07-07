package events

type ReceiveMessage struct {
	Type    MessageType
	Ident   string
	Content string
}
