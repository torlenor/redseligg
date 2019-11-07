// Concept how to build such a server inspired from https://github.com/grafana/grafana/blob/master/pkg/cmd/grafana-server/server.go
// Modified for AbyleBotter

package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/api"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/pool"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
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

// Run initializes and starts services. This will block until all services have
// exited. To initiate shutdown, call the Shutdown method in another goroutine.
func (s *Server) Run() (err error) {

	s.controlAPI, err = api.NewAPI(s.cfg)
	if err != nil {
		return fmt.Errorf("Error creating the API: %s", err.Error())
	}
	s.controlAPI.Init()
	s.controlAPI.AttachModuleGet("/status", statusEndpoint)
	services := []service{s.controlAPI}

	s.botPool = pool.NewBotPool(s.controlAPI)
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
