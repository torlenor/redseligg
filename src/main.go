package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"./botinterface"
	"./discord"
	"./matrix"
	"./plugins"
)

/**
 * Version should be set while build using ldflags (see Makefile)
 */
var version string

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

func connectPlugins(bot botinterface.Bot) {
	echoPlugin := plugins.CreateEchoPlugin(bot.GetReceiveMessageChannel(), bot.GetSendMessageChannel())
	echoPlugin.SetOnlyOnWhisper(true)
	echoPlugin.Start()
}

func main() {
	log.Println("AbyleBotter (" + version + ") is STARTING")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	discordBot := os.Getenv("DISCORD_BOT")
	matrixBot := os.Getenv("MATRIX_BOT")

	var bot botinterface.Bot

	if len(discordBot) > 0 {
		bot = discordBotCreator(done)
	} else if len(matrixBot) > 0 {
		bot = matrixBotCreator(done)
	} else {
		log.Fatal("No Bot chosen to start. Set DISCORD_BOT or MATRIX_BOT env variables")
	}

	connectPlugins(bot)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			bot.Stop()
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			log.Println("AbyleBotter gracefully shut down")
			break
		}
	}
}
