package main

import (
	"errors"
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

func start(done chan struct{}) {
	for _, bot := range bots.m {
		bot.Start(done)
	}
}

func createBots(cfg config.Config) error {
	if cfg.Bots.Discord.Enabled {
		bots.m["discord"] = discordBotCreator(cfg)
		if bots.m["discord"] == nil {
			return errors.New("Could not create Discord Bot")
		}
		if cfg.Bots.Discord.Plugins.Echo.Enabled {
			echoPlugin, err := plugins.CreateEchoPlugin()
			if err != nil {
				log.Errorln("Could not create EchoPlugin: ", err)
				return err
			}
			echoPlugin.SetOnlyOnWhisper(true)
			bots.m["discord"].AddPlugin(&echoPlugin)
			echoPlugin.Start()
		}
		if cfg.Bots.Discord.Plugins.SendMessage.Enabled {
			sendMessagesPlugin, err := plugins.CreateSendMessagesPlugin()
			if err != nil {
				log.Errorln("Could not create SendMessagesPlugin: ", err)
				return err
			}
			bots.m["discord"].AddPlugin(&sendMessagesPlugin)
			sendMessagesPlugin.Start()
		}

	} else if cfg.Bots.Matrix.Enabled {
		bots.m["matrix"] = matrixBotCreator(cfg)
		if bots.m["matrix"] == nil {
			return errors.New("Could not create Matrix Bot")
		}
		if cfg.Bots.Matrix.Plugins.Echo.Enabled {
			echoPlugin, err := plugins.CreateEchoPlugin()
			if err != nil {
				log.Errorln("Could not create EchoPlugin: ", err)
				return err
			}
			echoPlugin.SetOnlyOnWhisper(true)
			bots.m["matrix"].AddPlugin(&echoPlugin)
			echoPlugin.Start()
		}
		if cfg.Bots.Matrix.Plugins.SendMessage.Enabled {
			sendMessagesPlugin, err := plugins.CreateSendMessagesPlugin()
			if err != nil {
				log.Errorln("Could not create SendMessagesPlugin: ", err)
				return err
			}
			bots.m["matrix"].AddPlugin(&sendMessagesPlugin)
			sendMessagesPlugin.Start()
		}
	} else if cfg.Bots.Fake.Enabled {
		bots.m["fake"] = fakeBotCreator(cfg)
		if bots.m["fake"] == nil {
			return errors.New("Could not create Fake Bot")
		}
		if cfg.Bots.Fake.Plugins.Echo.Enabled {
			echoPlugin, err := plugins.CreateEchoPlugin()
			if err != nil {
				log.Errorln("Could not create EchoPlugin: ", err)
				return err
			}
			echoPlugin.SetOnlyOnWhisper(true)
			bots.m["fake"].AddPlugin(&echoPlugin)
			echoPlugin.Start()
		}
		if cfg.Bots.Fake.Plugins.SendMessage.Enabled {
			sendMessagesPlugin, err := plugins.CreateSendMessagesPlugin()
			if err != nil {
				log.Errorln("Could not create SendMessagesPlugin: ", err)
				return err
			}
			bots.m["fake"].AddPlugin(&sendMessagesPlugin)
			sendMessagesPlugin.Start()
		}
	}

	return nil
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

	err := createBots(cfg)
	if err != nil {
		log.Fatalln("Error initializing the bots and plugins:" + err.Error() + "Quitting...")
	}

	if len(bots.m) == 0 {
		log.Fatal("No Bot enabled. Check config file: ", configPath)
	}

	log.Println("AbyleBotter: Number of configured bots:", len(bots.m))

	startAbyleBotter()

}
