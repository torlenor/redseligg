package main

import (
	"fmt"

	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugins/echoplugin"
	"github.com/torlenor/abylebotter/plugins/httppingplugin"
	"github.com/torlenor/abylebotter/plugins/randomplugin"

	"github.com/torlenor/abylebotter/discord"
	"github.com/torlenor/abylebotter/matrix"
	"github.com/torlenor/abylebotter/mattermost"
	"github.com/torlenor/abylebotter/platform/slack"

	"github.com/torlenor/abylebotter/ws"

	newechoplugin "github.com/torlenor/abylebotter/plugin/echoplugin"
	newhttppingplugin "github.com/torlenor/abylebotter/plugin/httppingplugin"
	newrollplugin "github.com/torlenor/abylebotter/plugin/rollplugin"
)

func discordBotCreator(config config.Config) (*discord.Bot, error) {
	bot, err := discord.CreateDiscordBot(config.Bots.Discord)
	if err != nil {
		return nil, fmt.Errorf("Error creating DiscordBot: %s", err)
	}

	return bot, nil
}

func matrixBotCreator(config config.Config) (*matrix.Bot, error) {
	bot, err := matrix.CreateMatrixBot(config.Bots.Matrix)
	if err != nil {
		return nil, fmt.Errorf("Error creating MatrixBot: %s", err)
	}

	return bot, nil
}

func mattermostBotCreator(config config.MattermostConfig) (*mattermost.Bot, error) {
	bot, err := mattermost.CreateMattermostBot(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating MattermostBot: %s", err)
	}

	return bot, nil
}

func slackBotCreator(config config.SlackConfig) (*slack.Bot, error) {
	bot, err := slack.CreateSlackBot(config, ws.NewClient())
	if err != nil {
		return nil, fmt.Errorf("Error creating SlackBot: %s", err)
	}

	return bot, nil
}

func createPlugins(cfg config.Plugins, bot botinterface.Bot) error {
	if cfg.Echo.Enabled {
		echoPlugin, err := echoplugin.CreateEchoPlugin()
		if err != nil {
			log.Errorln("Could not create EchoPlugin: ", err)
			return err
		}
		echoPlugin.SetOnlyOnWhisper(cfg.Echo.OnlyWhispers)
		bot.AddPlugin(&echoPlugin)
		echoPlugin.Start()
	}
	if cfg.HTTPPing.Enabled {
		plugin, err := httppingplugin.CreateHTTPPingPlugin()
		if err != nil {
			log.Errorln("Could not create HTTPPingPlugin: ", err)
			return err
		}
		bot.AddPlugin(&plugin)
		plugin.Start()
	}
	if cfg.Random.Enabled {
		plugin, err := randomplugin.CreateRandomPlugin()
		if err != nil {
			log.Errorln("Could not create RandomPlugin: ", err)
			return err
		}
		bot.AddPlugin(&plugin)
		plugin.Start()
	}

	return nil
}

func createPlatformPlugins(cfg config.Plugins, bot platform.Bot) error {
	if cfg.Echo.Enabled {
		p := &newechoplugin.EchoPlugin{}
		p.SetOnlyOnWhisper(cfg.Echo.OnlyWhispers)
		bot.AddPlugin(p)
	}
	if cfg.Random.Enabled {
		plugin, err := newrollplugin.New()
		if err != nil {
			log.Errorln("Could not create RollPlugin: ", err)
			return err
		}
		bot.AddPlugin(&plugin)
	}
	if cfg.HTTPPing.Enabled {
		p := &newhttppingplugin.HTTPPingPlugin{}
		bot.AddPlugin(p)
	}

	return nil
}

func createBots(cfg config.Config) error {
	if cfg.Bots.Discord.Enabled {
		bot, err := discordBotCreator(cfg)
		if err != nil {
			return fmt.Errorf("Could not create Discord Bot: %s", err)
		}
		createPlugins(cfg.Bots.Discord.Plugins, bot)
		botPool.Add(bot)
	} else if cfg.Bots.Matrix.Enabled {
		bot, err := matrixBotCreator(cfg)
		if err != nil {
			return fmt.Errorf("Could not create Matrix Bot: %s", err)
		}
		createPlugins(cfg.Bots.Matrix.Plugins, bot)
		botPool.Add(bot)
	} else if cfg.Bots.Mattermost.Enabled {
		bot, err := mattermostBotCreator(cfg.Bots.Mattermost)
		if err != nil {
			return fmt.Errorf("Could not create Mattermost Bot: %s", err)
		}
		createPlugins(cfg.Bots.Mattermost.Plugins, bot)
		botPool.Add(bot)
	} else if cfg.Bots.Slack.Enabled {
		bot, err := slackBotCreator(cfg.Bots.Slack)
		if err != nil {
			return fmt.Errorf("Could not create Slack Bot: %s", err)
		}
		createPlatformPlugins(cfg.Bots.Slack.Plugins, bot)
		botPool.Add(bot)
	}

	return nil
}
