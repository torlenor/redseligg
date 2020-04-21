package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/utils"
)

const (
	defaultLoggingLevel = "info"
)

/**
 * version and compTime should be set while build using ldflags (see Makefile)
 */
var version string
var compTime string

var log *logrus.Entry

func setupLogging(loggingLevel string) {
	logging.Init()
	logging.SetLoggingLevel(loggingLevel)

	log = logging.Get("main")
}

func main() {

	fmt.Printf("Botter Version %s (%s)\n\n", version, compTime)

	var (
		loggingLevel = flag.String("l", defaultLoggingLevel, "Logging level (panic, fatal, error, warn/warning, info or debug)")
		v            = flag.Bool("v", false, "prints current version and exits")
	)

	flag.Parse()

	if *v {
		fmt.Printf("Version %s (%s)\n", version, compTime)
		os.Exit(0)
	}

	setupLogging(*loggingLevel)

	controlAPIConfig := config.API{
		Enabled: false,
	}

	utils.Version().Set(version)
	utils.Version().SetCompTime(compTime)

	server := NewServer(controlAPIConfig)

	go listenToSystemSignals(server)

	err := server.Run()

	code := server.ExitCode(err)

	os.Exit(code)
}

func listenToSystemSignals(server *Server) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case sig := <-signalChan:
			server.Shutdown(fmt.Sprintf("System signal: %s", sig))
		}
	}
}
