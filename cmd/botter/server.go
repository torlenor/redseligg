// Concept how to build such a server inspired from https://github.com/grafana/grafana/blob/master/pkg/cmd/grafana-server/server.go
// Modified for AbyleBotter

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	mongodbbotconfigprovider "github.com/torlenor/redseligg/botconfigprovider/mongodb"
	tomlbotconfigprovider "github.com/torlenor/redseligg/botconfigprovider/toml"

	"github.com/torlenor/redseligg/api"
	"github.com/torlenor/redseligg/config"
	"github.com/torlenor/redseligg/factories"
	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/pool"
	"github.com/torlenor/redseligg/providers"
)

// NewServer returns a new instance of Server
func NewServer(controlAPIConfig config.API) *Server {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	return &Server{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
		log:           logging.Get("server"),

		cfg: controlAPIConfig,
	}
}

type service interface {
	Run(context.Context) error
}

// Server is responsible for running the botter instance
type Server struct {
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	log                *logrus.Entry
	cfg                config.API
	shutdownReason     string
	shutdownInProgress bool

	controlAPI *api.API
	botPool    *pool.BotPool
}

func createBotProvider() (*providers.BotProvider, error) {
	cfgSource, exists := os.LookupEnv("BOTTER_BOT_CFG_SOURCE")
	if !exists {
		return nil, fmt.Errorf("Error determining bot config type, environment variable BOTTER_BOT_CFG_SOURCE not set")
	}
	cfgSource = strings.ToUpper(cfgSource)

	botFactory := &factories.BotFactory{}
	pluginFactory := &factories.PluginFactory{}

	switch cfgSource {
	case "TOML":
		tomlFile, exists := os.LookupEnv("BOTTER_BOT_CFG_TOML_FILE")
		if !exists {
			tomlFile = "/cfg/bots.toml"
		}
		cfgs, err := tomlbotconfigprovider.ParseTomlBotConfigFromFile(tomlFile)
		if err != nil {
			return nil, fmt.Errorf("Error parsing the toml bot config %s (check env variable BOTTER_BOT_CFG_TOML_FILE)", err)
		}

		botProvider, err := providers.NewBotProvider(cfgs, botFactory, pluginFactory)
		if err != nil {
			return nil, fmt.Errorf("Error creating bot provider: %s", err)
		}
		return botProvider, nil
	case "MONGODB":
		fallthrough
	case "MONGO":
		url, exists := os.LookupEnv("BOTTER_BOT_CFG_MONGO_URL")
		if !exists {
			return nil, fmt.Errorf("Error in setting up MongoDB bot config: MongoDB URL not set, check environment variable BOTTER_BOT_CFG_MONGO_URL")
		}
		db, exists := os.LookupEnv("BOTTER_BOT_CFG_MONGO_DB")
		if !exists {
			return nil, fmt.Errorf("Error in setting up MongoDB bot config: MongoDB database not set, check environment variable BOTTER_BOT_CFG_MONGO_DB")
		}
		cfgs, err := mongodbbotconfigprovider.NewBackend(url, db)
		if err != nil {
			return nil, fmt.Errorf("Error in setting up MongoDB bot config: %s", err)
		}
		cfgs.Connect()

		botProvider, err := providers.NewBotProvider(cfgs, botFactory, pluginFactory)
		if err != nil {
			return nil, fmt.Errorf("Error creating bot provider: %s", err)
		}
		return botProvider, nil
	default:
		return nil, fmt.Errorf("Unknown BOTTER_BOT_CFG_SOURCE specified")
	}
}

// Run initializes and starts services. This will block until all services have
// exited. To initiate shutdown, call the Shutdown method in another goroutine.
func (s *Server) Run() (err error) {
	services := []service{}
	if s.cfg.Enabled {
		s.controlAPI, err = api.NewAPI(s.cfg, "/v1")
		if err != nil {
			return fmt.Errorf("Error creating the API: %s", err.Error())
		}
		s.controlAPI.Init()
		services = append(services, s.controlAPI)
	}

	botProvider, err := createBotProvider()
	if err != nil {
		return fmt.Errorf("Error creating the BotProvider: %s", err.Error())
	}
	s.botPool, err = pool.NewBotPool(s.controlAPI, botProvider)
	if err != nil {
		return fmt.Errorf("Error creating the BotPool: %s", err.Error())
	}
	services = append(services, s.botPool)

	// Start background services
	for _, svc := range services {
		descriptor := svc
		s.childRoutines.Go(func() error {
			// Don't start new services when server is shutting down.
			if s.shutdownInProgress {
				return nil
			}

			if err := descriptor.Run(s.context); err != nil {
				if err != context.Canceled {
					s.log.Errorln("Stopped ", "reason", err)
				} else {
					s.log.Infoln("Stopped ", "reason", err)
				}
			}

			// Mark that we are in shutdown mode
			// So more services are not started
			s.shutdownInProgress = true
			return nil
		})
	}

	// Start all enabled bots
	for _, id := range botProvider.GetAllEnabledBotIDs() {
		err := s.botPool.AddViaID(id)
		if err != nil {
			s.log.Errorf("Error creating bot with ID %s: %s", id, err)
		}
	}

	defer func() {
		s.log.Debug("Waiting on services...")
		if waitErr := s.childRoutines.Wait(); waitErr != nil && !xerrors.Is(waitErr, context.Canceled) {
			s.log.Errorln("A service failed", "err", waitErr)
			if err == nil {
				err = waitErr
			}
		}
	}()

	return
}

// Shutdown the server instance
func (s *Server) Shutdown(reason string) {
	s.log.Infoln("Shutdown started", "reason", reason)
	s.shutdownReason = reason
	s.shutdownInProgress = true

	s.shutdownFn()

	if err := s.childRoutines.Wait(); err != nil && !xerrors.Is(err, context.Canceled) {
		s.log.Errorln("Failed waiting for services to shutdown", "err", err)
	}
}

// ExitCode returns an exit code for a given error
func (s *Server) ExitCode(reason error) int {
	code := 1

	if reason == context.Canceled && s.shutdownReason != "" {
		reason = fmt.Errorf(s.shutdownReason)
		code = 0
	}

	s.log.Errorln("Server shutdown", "reason", reason)

	return code
}
