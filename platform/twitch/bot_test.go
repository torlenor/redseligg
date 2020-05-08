package twitch

import (
	"context"
	"testing"
	"time"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/ws"
)

func Test_CreateDiscordBot(t *testing.T) {
	ws := &ws.MockClient{}
	storage := storage.MockStorage{}
	cfg := botconfig.TwitchConfig{}

	_, err := CreateTwitchBot(cfg, &storage, ws)
	if err != nil {
		t.Fatalf("Creating the bot should not have failed")
	}
}

func Test_DiscordBot_Run(t *testing.T) {
	// assert := assert.New(t)

	ws := &ws.MockClient{}
	storage := storage.MockStorage{}
	cfg := botconfig.TwitchConfig{}

	bot, err := CreateTwitchBot(cfg, &storage, ws)
	if err != nil {
		t.Fatalf("Creating the bot should not have failed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go bot.Run(ctx)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}
