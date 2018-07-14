package events

// MessageType enum
type MessageType int

// Known MessageTypes
const ( // iota is reset to 0
	UNKNOWN MessageType = iota
	WHISPER
	MESSAGE
)

// MessageTypes as strings
var MessageTypes = [...]string{
	"UNKNOWN",
	"WHISPER",
	"MESSAGE",
}

func (messageType MessageType) String() string {
	return MessageTypes[messageType]
}
