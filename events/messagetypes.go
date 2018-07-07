package events

type MessageType int

const ( // iota is reset to 0
	UNKNOWN MessageType = iota
	WHISPER
	MESSAGE
)

var MessageTypes = [...]string{
	"UNKNOWN",
	"WHISPER",
	"MESSAGE",
}

func (messageType MessageType) String() string {
	return MessageTypes[messageType]
}
