package platform

import (
	"context"

	"github.com/torlenor/redseligg/commanddispatcher"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storage"
)

// All currently supported features a platform can support
const (
	FeatureMessagePost    string = "FEATURE_MESSAGE_POST"
	FeatureMessageUpdate  string = "FEATURE_MESSAGE_UPDATE"
	FeatureMessageDelete  string = "FEATURE_MESSAGE_DELETE"
	FeatureReactionNotify string = "FEATURE_REACTION_NOTIFY"
)

// Bot type interface which every Bot has to implement
type Bot interface {
	// Run the Bot (blocking)
	Run(ctx context.Context) error

	AddPlugin(plugin BotPlugin)

	GetInfo() BotInfo
}

// BotPlugin is needed to connect a Plugin to a Bot
type BotPlugin interface {
	plugin.Hooks

	SetBotPluginID(botID string, pluginID string)

	SetAPI(api plugin.API) error
}

// PluginInfo contains info about one plugin
type PluginInfo struct {
	Plugin string `json:"plugin"`
	Active bool   `json:"active"`
}

// BotInfo contains info about one bot
type BotInfo struct {
	BotID    string       `json:"botId"`
	Platform string       `json:"platform"`
	Healthy  bool         `json:"healthy"`
	Plugins  []PluginInfo `json:"plugins"`
}

// BotImpl gives default implementations and basic functionalities for a bot.
// Embed this struct to get a default implementation of the plugin API for the bot.
type BotImpl struct {
	ProvidedFeatures map[string]bool

	Dispatcher *commanddispatcher.CommandDispatcher
	Storage    storage.Storage
}

// HasFeature returns true if the bot serving the API implements the feature.
func (b *BotImpl) HasFeature(feature string) bool {
	return b.ProvidedFeatures[feature]
}

// GetStorage returns the storage or nil if none is provided by the platform
func (b *BotImpl) GetStorage() storage.Storage { return b.Storage }

// RegisterCommand registers a custom slash or ! command, depending on what the bot supports.
func (b *BotImpl) RegisterCommand(p plugin.Hooks, command string) error {
	b.Dispatcher.Register(command, p)
	return nil
}

// UnRegisterCommand unregisters a command previously registered via RegisterCommand.
func (b *BotImpl) UnRegisterCommand(command string) error {
	b.Dispatcher.Unregister(command)
	return nil
}
