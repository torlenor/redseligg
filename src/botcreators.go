package main

import (
	"fake"
	"mattermost"
	"os"

	"config"
	"discord"
	"matrix"
)

func discordBotCreator(config config.Config) *discord.Bot {
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	if len(discordToken) == 0 {
		discordToken = config.Bots.Discord.Token
	}

	bot, err := discord.CreateDiscordBot(config.Bots.Discord.ID, config.Bots.Discord.Secret, discordToken)
	if err != nil {
		log.Println("DiscordBot: ERROR: ", err)
	}

	return bot
}

func matrixBotCreator(config config.Config) *matrix.Bot {
	matrixServer := os.Getenv("MATRIX_SERVER")
	if len(matrixServer) == 0 {
		matrixServer = config.Bots.Matrix.Server
	}
	matrixUsername := os.Getenv("MATRIX_USERNAME")
	if len(matrixUsername) == 0 {
		matrixUsername = config.Bots.Matrix.Username
	}
	matrixPassword := os.Getenv("MATRIX_PASSWORD")
	if len(matrixPassword) == 0 {
		matrixPassword = config.Bots.Matrix.Password
	}
	matrixToken := os.Getenv("MATRIX_TOKEN")
	if len(matrixToken) == 0 {
		matrixToken = config.Bots.Matrix.Token
	}

	bot, err := matrix.CreateMatrixBot("https://matrix.abyle.org", matrixUsername, matrixPassword, matrixToken)
	if err != nil {
		log.Println("MatrixBot: ERROR: ", err)
	}

	return bot
}

func fakeBotCreator(config config.Config) *fake.Bot {
	bot, err := fake.CreateFakeBot()
	if err != nil {
		log.Println("FakeBot: ERROR: ", err)
	}

	return bot
}

func mattermostBotCreator(config config.MattermostConfig) *mattermost.Bot {
	bot, err := mattermost.CreateMattermostBot(config)
	if err != nil {
		log.Println("FakeBot: ERROR: ", err)
	}

	return bot
}
