package matrix

import (
	"context"
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
	log.Println("MatrixBot is RUNNING")
}

// Run the Matrix Bot (blocking)
func (b *Bot) Run(ctx context.Context) error {
	b.Start()

	<-ctx.Done()

	b.Stop()

	return nil
}

// Stop the Matrix Bot
func (b *Bot) Stop() {
	log.Println("MatrixBot is SHUTING DOWN")

	b.pollingDone <- true

	log.Println("MatrixBot is SHUT DOWN")
}
