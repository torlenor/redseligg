package events

// Command describes the struct to send commands over the command channel between plugin and bot
type Command struct {
	Command string
	Payload string
}
