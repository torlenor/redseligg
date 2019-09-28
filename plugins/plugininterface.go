package plugins

import "github.com/torlenor/abylebotter/events"

// Plugin type interface
type Plugin interface {
	ConnectChannels(receiveChannel chan events.ReceiveMessage,
		sendChannel chan events.SendMessage,
		commandCHannel chan events.Command) error

	GetName() string

	Start()
	Stop()
}
