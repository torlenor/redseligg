package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/logging"
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

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "Path to toml config file")
	flag.StringVar(&loggingLevel, "l", defaultLoggingLevel, "Logging level (panic, fatal, error, warn/warning, info or debug)")
	flag.Parse()
}

func setupLogging() {
	logging.Init()
	logging.SetLoggingLevel(loggingLevel)

	log = logging.Get("main")
}

func main() {
	setupLogging()

	log.Println("BotterInstance (" + version + ") is STARTING")

	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-interrupt:
			log.Println("BotterInstance gracefully shut down")
			return
		}
	}
}
