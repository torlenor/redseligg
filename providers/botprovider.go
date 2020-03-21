package providers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
)

type configProvider interface {
	GetBotConfig(id string) (botconfig.BotConfig, error)
}

type botFactory interface {
	CreateBot(platform string, config botconfig.BotConfig) (platform.Bot, error)
}

type pluginFactory interface {
	CreatePlugin(plugin string, pluginConfig botconfig.PluginConfig) (platform.BotPlugin, error)
}

// BotProvider creates configured bots ready to run
type BotProvider struct {
	log *logrus.Entry

	botConfigs    configProvider
	botFactory    botFactory
	pluginFactory pluginFactory
}

// NewBotProvider creates a new BotProvider
func NewBotProvider(botConfigProvider configProvider, botFactory botFactory, pluginFactory pluginFactory) (*BotProvider, error) {
	bp := BotProvider{
		log: logging.Get("BotProvider"),
	}

	bp.botConfigs = botConfigProvider
	bp.botFactory = botFactory
	bp.pluginFactory = pluginFactory

	bp.log.Debug("New BotProvider created")

	return &bp, nil
}

func (b *BotProvider) createPlatformPlugins(plugins map[string]botconfig.PluginConfig, bot platform.Bot) error {
	var lastError error

	for _, plugin := range plugins {
		p, err := b.pluginFactory.CreatePlugin(plugin.Type, plugin)
		if err != nil {
			lastError = err
			continue
		}
		bot.AddPlugin(p)
	}

	if lastError != nil {
		return fmt.Errorf("Could not create all plugins, last error was: %s", lastError)
	}

	return nil
}

// GetBot creates the bot with the given id
func (b *BotProvider) GetBot(id string) (platform.Bot, error) {
	var botConfig botconfig.BotConfig
	var err error
	if botConfig, err = b.botConfigs.GetBotConfig(id); err != nil {
		return nil, fmt.Errorf("Bot ID %s not known: %s", id, err)
	}

	var bot platform.Bot

	bot, err = b.botFactory.CreateBot(botConfig.Type, botConfig)
	if err != nil {
		return nil, fmt.Errorf("Error creating bot with id %s: %s", id, err)
	}
	err = b.createPlatformPlugins(botConfig.Plugins, bot)
	if err != nil {
		b.log.Warnf("Error adding plugins to the bot with id %s: %s", id, err)
	}

	return bot, nil
}
