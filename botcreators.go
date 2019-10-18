package main

import (
	"fmt"

	"github.com/torlenor/abylebotter/config"

	"github.com/torlenor/abylebotter/discord"
	"github.com/torlenor/abylebotter/matrix"
	"github.com/torlenor/abylebotter/mattermost"
	"github.com/torlenor/abylebotter/slack"

	"github.com/torlenor/abylebotter/ws"
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
