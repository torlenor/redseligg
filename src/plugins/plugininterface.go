package plugins

import "events"

// Plugin type interface
type Plugin interface {
	ConnectChannels(receiveChannel chan events.ReceiveMessage,
		sendChannel chan events.SendMessage,
		commandCHannel chan events.Command) error

	Start()
	Stop()
}
