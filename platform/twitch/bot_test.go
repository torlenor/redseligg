package twitch

import (
	"context"
	"testing"
	"time"

	"github.com/torlenor/redseligg/botconfig"
	"github.com/torlenor/redseligg/commanddispatcher"

	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/ws"
)

func Test_CreateTwitchBot(t *testing.T) {
	ws := &ws.MockClient{}
	dispatcher := commanddispatcher.CommandDispatcher{}
	storage := storage.MockStorage{}
	cfg := botconfig.TwitchConfig{}

	_, err := CreateTwitchBot(cfg, &storage, &dispatcher, ws)
	if err != nil {
		t.Fatalf("Creating the bot should not have failed")
	}
}

func Test_TwitchBot_Run(t *testing.T) {
	ws := &ws.MockClient{}
	dispatcher := commanddispatcher.CommandDispatcher{}
	storage := storage.MockStorage{}
	cfg := botconfig.TwitchConfig{}

	bot, err := CreateTwitchBot(cfg, &storage, &dispatcher, ws)
	if err != nil {
		t.Fatalf("Creating the bot should not have failed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go bot.Run(ctx)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}
