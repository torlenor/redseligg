package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/torlenor/AbyleBotter/botinterface"
	"github.com/torlenor/AbyleBotter/discord"
	"github.com/torlenor/AbyleBotter/plugins"
)

func main() {
	log.Println("AbyleBotter is STARTING")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	token := os.Getenv("DISCORD_BOT_TOKEN")

	var bot botinterface.Bot = discord.CreateDiscordBot(token)
	bot.Start(done)

	echoPlugin := plugins.CreateEchoPlugin(bot.GetReceiveMessageChannel(), bot.GetSendMessageChannel())
	echoPlugin.SetOnlyOnWhisper(true)
	echoPlugin.Start()

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
