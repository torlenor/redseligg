package twitch

import (
	"context"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin"
	"github.com/torlenor/abylebotter/storage"
)

var (
	log = logging.Get("TwitchBot")
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	storage storage.Storage
	plugins []plugin.Hooks

	cfg botconfig.TwitchConfig
}

// CreateTwitchBot creates a new instance of a TwitchBot
func CreateTwitchBot(cfg botconfig.TwitchConfig, storage storage.Storage) (*Bot, error) {
	log.Info("TwitchBot is CREATING itself")

	b := Bot{
		storage: storage,
		cfg:     cfg,
	}

	return &b, nil
}

// Run the Bot (blocking)
func (b *Bot) Run(ctx context.Context) error {
	// RUN SOMETHING

	for _, plugin := range b.plugins {
		plugin.OnRun()
	}

	<-ctx.Done()

	for _, plugin := range b.plugins {
		plugin.OnStop()
	}

	// STOP SOMETHING

	return nil
}

// AddPlugin takes as argument a plugin and
// adds it to the bot providing it with the API
func (b *Bot) AddPlugin(plugin platform.BotPlugin) {
	plugin.SetAPI(b)
	b.plugins = append(b.plugins, plugin)
}

// GetInfo returns information about the Bot
func (b *Bot) GetInfo() platform.BotInfo {
	return platform.BotInfo{
		BotID:    "",
		Platform: "Twitch",
		Healthy:  true,
		Plugins:  []platform.PluginInfo{},
	}
}
