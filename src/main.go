package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"botinterface"
	"config"
	"logging"
	"plugins"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

const (
	defaultConfigPath = "./abylebotter.toml"
)

/**
 * Version should be set while build using ldflags (see Makefile)
 */
var version string
var configPath string

var bots = struct {
	m map[string]botinterface.Bot
}{m: make(map[string]botinterface.Bot)}

var interrupt chan os.Signal

var log *logrus.Entry

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "Path to toml config file")
	flag.Parse()
}

func connectPlugins(cfg config.Config, bot botinterface.Bot) {

}

func start(done chan struct{}) {
	for _, bot := range bots.m {
		bot.Start(done)
	}
}

func createBots(cfg config.Config) {
	if cfg.Bots.Discord.Enabled {
		bots.m["discord"] = discordBotCreator(cfg)
		if cfg.Bots.Discord.Plugins.Echo.Enabled {
			echoPlugin := plugins.CreateEchoPlugin()
			echoPlugin.SetOnlyOnWhisper(true)
			bots.m["discord"].AddPlugin(&echoPlugin)
			echoPlugin.Start()
		}
	} else if cfg.Bots.Matrix.Enabled {
		bots.m["matrix"] = matrixBotCreator(cfg)
		if cfg.Bots.Matrix.Plugins.Echo.Enabled {
			echoPlugin := plugins.CreateEchoPlugin()
			echoPlugin.SetOnlyOnWhisper(true)
			bots.m["matrix"].AddPlugin(&echoPlugin)
			echoPlugin.Start()
		}
	} else if cfg.Bots.Fake.Enabled {
		bots.m["fake"] = fakeBotCreator(cfg)
		if cfg.Bots.Fake.Plugins.Echo.Enabled {
			echoPlugin := plugins.CreateEchoPlugin()
			echoPlugin.SetOnlyOnWhisper(true)
			bots.m["fake"].AddPlugin(&echoPlugin)
			echoPlugin.Start()
		}
	}
}

func startAbyleBotter() {
	log.Println("Starting the bots")

	done := make(chan struct{})
	start(done)

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				for botname, bot := range bots.m {
					if bot.Status().Fatal {
						log.Println("Status of bot", botname, " it FATAL, trying to recover...")
						bot.Stop()
						bot.Start(done)
					}
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			for _, bot := range bots.m {
				bot.Stop()
			}
			ticker.Stop()
			log.Println("AbyleBotter gracefully shut down")
			break
		}
	}
}

func main() {
	logging.Init()

	log = logging.Get("main")

	log.Println("AbyleBotter (" + version + ") is STARTING")

	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var cfg config.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		fmt.Println(err)
		return
	}

	createBots(cfg)

	if len(bots.m) == 0 {
		log.Fatal("No Bot enabled. Check config file: ", configPath)
	}

	log.Println("AbyleBotter: Number of configured bots:", len(bots.m))

	startAbyleBotter()

}
