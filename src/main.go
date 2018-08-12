package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"botinterface"
	"config"
	"plugins"

	"github.com/BurntSushi/toml"
)

const (
	defaultConfigPath = "./abylebotter.toml"
)

/**
 * Version should be set while build using ldflags (see Makefile)
 */
var version string
var configPath string

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "Path to toml config file")

	flag.Parse()
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

	var config config.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("\n", config.Bots.Discord.Plugins.Echo.Enabled, "\n")

	fmt.Print("\n", config.Bots.Discord.Enabled, "\n")

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
