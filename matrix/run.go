package matrix

import (
	"time"
)

func (b *Bot) handlePolling() error {
	return b.callSync()
}

func (b *Bot) startBot() {
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
func (b *Bot) Start() {
	log.Println("MatrixBot is STARTING")
	go b.startBot()
	go b.startSendChannelReceiver()
	log.Println("MatrixBot is RUNNING")
}

// Stop the Matrix Bot
func (b *Bot) Stop() {
	log.Println("MatrixBot is SHUTING DOWN")

	b.pollingDone <- true

	b.disconnectReceivers()

	log.Println("MatrixBot is SHUT DOWN")
}

func (b *Bot) disconnectReceivers() {
	for plugin, pluginChannel := range b.receivers {
		log.Debugln("Disconnecting Plugin", plugin.GetName())
		defer close(pluginChannel)
	}
}
