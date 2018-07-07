package botinterface

import "github.com/torlenor/AbyleBotter/events"

// Bot type interface
type Bot interface {
	Start(doneChannel chan struct{})
	Stop()

	GetReceiveMessageChannel() chan events.ReceiveMessage
}
