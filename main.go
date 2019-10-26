package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/api"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/pool"
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
