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

func connectPlugins(cfg config.Config, bot botinterface.Bot) {
	echoPlugin := plugins.CreateEchoPlugin(bot.GetReceiveMessageChannel(), bot.GetSendMessageChannel())
	echoPlugin.SetOnlyOnWhisper(true)
	echoPlugin.Start()
}

func start(bots []botinterface.Bot, done chan struct{}) {
	for _, bot := range bots {
		bot.Start(done)
	}
}

func main() {
	log.Println("AbyleBotter (" + version + ") is STARTING")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var cfg config.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		fmt.Println(err)
		return
	}

	var bots []botinterface.Bot
	if cfg.Bots.Discord.Enabled {
		bots = append(bots, discordBotCreator(cfg))
		connectPlugins(cfg, bots[len(bots)-1])
	} else if cfg.Bots.Matrix.Enabled {
		bots = append(bots, matrixBotCreator(cfg))
		connectPlugins(cfg, bots[len(bots)-1])
	}

	if len(bots) == 0 {
		log.Fatal("No Bot enabled. Check config file: ", configPath)
	}

	done := make(chan struct{})
	start(bots, done)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			for _, bot := range bots {
				bot.Stop()
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			log.Println("AbyleBotter gracefully shut down")
			break
		}
	}
}
