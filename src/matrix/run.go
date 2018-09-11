package matrix

import (
	"time"
)

func (b *Bot) handlePolling() error {
	return b.callSync()
}

func (b *Bot) startBot(doneChannel chan struct{}) {
	defer close(doneChannel)
	// do some message polling or whatever until stopped
	tickChan := time.Tick(b.pollingInterval)

	for {
		select {
		case <-tickChan:
			b.handlePolling()
		case <-b.pollingDone:
			return
		}
	}
}

// Start the Matrix Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	log.Println("MatrixBot is STARTING")
	go b.startBot(doneChannel)
	go b.startSendChannelReceiver()
	go b.startCommandChannelReceiver()
	log.Println("MatrixBot is RUNNING")
}

// Stop the Matrix Bot
func (b Bot) Stop() {
	log.Println("MatrixBot is SHUTING DOWN")

	b.pollingDone <- true

	defer close(b.receiveMessageChan)

	log.Println("MatrixBot is SHUT DOWN")
}
