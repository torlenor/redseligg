package main

import (
	"log"
	"os"

	"./discord"
	"./matrix"
)

func discordBotCreator(done chan struct{}) *discord.Bot {
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")

	bot := discord.CreateDiscordBot(discordToken)
	// if err != nil {
	// 	log.Println("DiscordBot: ERROR: ", err)
	// }
	bot.Start(done)

	return bot
}

func matrixBotCreator(done chan struct{}) *matrix.Bot {
	matrixUsername := os.Getenv("MATRIX_BOT_USERNAME")
	matrixPassword := os.Getenv("MATRIX_BOT_PASSWORD")
	matrixToken := os.Getenv("MATRIX_BOT_TOKEN")

	bot, err := matrix.CreateMatrixBot("https://matrix.abyle.org", matrixUsername, matrixPassword, matrixToken)
	if err != nil {
		log.Println("MatrixBot: ERROR: ", err)
	}
	bot.Start(done)

	return bot
}
