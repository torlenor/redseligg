package twitch

import (
	"context"
	"testing"
	"time"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/storage"
)

func Test_CreateDiscordBot(t *testing.T) {
	storage := storage.MockStorage{}
	cfg := botconfig.TwitchConfig{}

	_, err := CreateTwitchBot(cfg, &storage)
	if err != nil {
		t.Fatalf("Creating the bot should not have failed")
	}
}

func Test_DiscordBot_Run(t *testing.T) {
	// assert := assert.New(t)

	storage := storage.MockStorage{}
	cfg := botconfig.TwitchConfig{}

	bot, err := CreateTwitchBot(cfg, &storage)
	if err != nil {
		t.Fatalf("Creating the bot should not have failed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go bot.Run(ctx)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}
