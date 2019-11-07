package providers

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/platform/mattermost"
	"github.com/torlenor/abylebotter/platform/slack"
	"github.com/torlenor/abylebotter/plugin/echoplugin"
	"github.com/torlenor/abylebotter/plugin/httppingplugin"
	"github.com/torlenor/abylebotter/plugin/rollplugin"
	"github.com/torlenor/abylebotter/ws"
)

type configProvider interface {
	GetBotConfig(id string) (BotConfig, error)
}

// BotProvider creates configured bots ready to run
type BotProvider struct {
	log *logrus.Entry

	// botConfigs TomlBotConfig
	botConfigs configProvider
}

// NewBotProvider creates a new BotProvider
func NewBotProvider(botConfigProvider configProvider) (*BotProvider, error) {
	bp := BotProvider{
		log: logging.Get("BotProvider"),
	}

	bp.botConfigs = botConfigProvider

	bp.log.Debug("New BotProvider created")

	return &bp, nil
}

func (b *BotProvider) createPlatformPlugins(plugins map[string]PluginConfig, bot platform.Bot) error {
	for _, plugin := range plugins {
		switch plugin.Type {
		case "echo":
			p := &echoplugin.EchoPlugin{}
			bot.AddPlugin(p)
		case "roll":
			plugin, err := rollplugin.New()
			if err != nil {
				return err
			}
			bot.AddPlugin(&plugin)
		case "httpping":
			p := &httppingplugin.HTTPPingPlugin{}
			bot.AddPlugin(p)
		default:
			b.log.Warnf("Unknown plugin type %s", plugin)
		}
	}

	return nil
}

// GetBot creates the bot with the given id
func (b *BotProvider) GetBot(id string) (platform.Bot, error) {
	var botConfig BotConfig
	var err error
	if botConfig, err = b.botConfigs.GetBotConfig(id); err != nil {
		return nil, fmt.Errorf("Bot ID %s not known: %s", id, err)
	}

	var bot platform.Bot

	switch botConfig.Type {
	case "slack":
		slackCfg, err := botConfig.AsSlackConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating slack bot: %s", err)
		}

		bot, err = slack.CreateSlackBot(slackCfg, ws.NewClient())
		if err != nil {
			return nil, fmt.Errorf("Error creating slack bot: %s", err)
		}
	case "mattermost":
		mmCfg, err := botConfig.AsMattermostConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating mattermost bot: %s", err)
		}

		bot, err = mattermost.CreateMattermostBot(mmCfg)
		if err != nil {
			return nil, fmt.Errorf("Error creating mattermost bot: %s", err)
		}
	default:
		return nil, fmt.Errorf("Unknown Bot type %s", botConfig.Type)
	}

	err = b.createPlatformPlugins(botConfig.Plugins, bot)
	if err != nil {
		return nil, fmt.Errorf("Error creating plugins: %s", err)
	}

	return bot, nil
}
