package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/api"
	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/pool"

	"github.com/torlenor/abylebotter/plugins/echoplugin"
	"github.com/torlenor/abylebotter/plugins/httppingplugin"
	"github.com/torlenor/abylebotter/plugins/randomplugin"
	"github.com/torlenor/abylebotter/plugins/sendmessagesplugin"
)

const (
	defaultConfigPath   = "./abylebotter.toml"
	defaultLoggingLevel = "info"
)

/**
 * Version should be set while build using ldflags (see Makefile)
 */
var version string

var configPath string
var loggingLevel string

var interrupt chan os.Signal

var log *logrus.Entry

var botPool pool.Pool

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "Path to toml config file")
	flag.StringVar(&loggingLevel, "l", defaultLoggingLevel, "Logging level (panic, fatal, error, warn/warning, info or debug)")
	flag.Parse()
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
	if cfg.SendMessage.Enabled {
		sendMessagesPlugin, err := sendmessagesplugin.CreateSendMessagesPlugin()
		sendMessagesPlugin.RegisterToRestAPI()
		if err != nil {
			log.Errorln("Could not create SendMessagesPlugin: ", err)
			return err
		}
		bot.AddPlugin(&sendMessagesPlugin)
		sendMessagesPlugin.Start()
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
		createPlugins(cfg.Bots.Slack.Plugins, bot)
		botPool.Add(bot)
	}

	return nil
}

func startAbyleBotter() {
	log.Println("Starting the bots")

	botPool.StartAll()

	for {
		select {
		case <-interrupt:
			botPool.StopAll()
			log.Println("AbyleBotter gracefully shut down")
			return
		}
	}
}

func setupLogging() {
	logging.Init()
	logging.SetLoggingLevel(loggingLevel)

	log = logging.Get("main")
}

func main() {
	setupLogging()

	log.Println("AbyleBotter (" + version + ") is STARTING")

	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	cfg, err := config.ParseFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = createBots(cfg)
	if err != nil {
		log.Fatalln("Error initializing the bots and plugins:" + err.Error() + "Quitting...")
	}

	if botPool.Len() == 0 {
		log.Fatal("No Bot enabled. Check config file: ", configPath)
	}

	log.Infoln("AbyleBotter: Number of configured bots:", botPool.Len())

	// Start API
	go api.Start(cfg.General.API)

	startAbyleBotter()
}
