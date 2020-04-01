package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/utils"
)

const (
	defaultLoggingLevel  = "info"
	defaultPort          = "9081"
	defaultListenAddress = "0.0.0.0"
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

func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	status := fmt.Sprintf(`{"status":"OK", "version":"%s"}`, version)

	io.WriteString(w, string(status))
}

func main() {

	var (
		loggingLevel  = flag.String("l", defaultLoggingLevel, "Logging level (panic, fatal, error, warn/warning, info or debug)")
		port          = flag.String("p", defaultPort, "Port for the Control API")
		listenAddress = flag.String("c", defaultListenAddress, "Listen address for the Control API")
		v             = flag.Bool("v", false, "prints current version and exits")
	)

	flag.Parse()

	if *v {
		fmt.Printf("Version %s (%s)\n", version, compTime)
		os.Exit(0)
	}

	setupLogging(*loggingLevel)

	controlAPIConfig := config.API{
		Enabled: true,
		Port:    *port,
		IP:      *listenAddress,
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
