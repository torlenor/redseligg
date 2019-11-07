package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
)

// Router is the interface to the router used for distributing to the endpoints
type router interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// API represents a Rest API instance of a ALoLStats instance
type API struct {
	config config.API
	router router
	log    *logrus.Entry
	prefix string

	server *http.Server
	quit   chan interface{}
}

// NewAPI creates a new Rest API instance
func NewAPI(cfg config.API, prefix string) (*API, error) {
	a := &API{
		config: cfg,
		router: mux.NewRouter(),
		log:    logging.Get("RestAPI"),
		prefix: prefix,
	}

	if len(a.config.Port) == 0 {
		return nil, fmt.Errorf("REST API activated but no valid configuration found. At least port has to specified")
	}

	return a, nil
}

// NewAPICustom creates a new Rest API instance with a custom router
func NewAPICustom(cfg config.API, prefix string, router router) (*API, error) {
	a := &API{
		config: cfg,
		router: router,
		log:    logging.Get("RestAPI"),
		prefix: prefix,
	}

	if len(a.config.Port) == 0 {
		return nil, fmt.Errorf("REST API activated but no valid configuration found. At least port has to specified")
	}

	return a, nil
}

// AttachModuleGet registers a new GET handler for the API
func (a *API) AttachModuleGet(path string, f func(http.ResponseWriter, *http.Request)) {
	a.log.Infoln("Registering GET handler:", a.prefix+path)
	a.router.HandleFunc(a.prefix+path, f).Methods("GET")
}

// AttachModulePost registers a new POST handler for the API
func (a *API) AttachModulePost(path string, f func(http.ResponseWriter, *http.Request)) {
	a.log.Infoln("Registering POST handler:", a.prefix+path)
	a.router.HandleFunc(a.prefix+path, f).Methods("POST")
}

// AttachModulePut registers a new PUT handler for the API
func (a *API) AttachModulePut(path string, f func(http.ResponseWriter, *http.Request)) {
	a.log.Infoln("Registering PUT handler:", a.prefix+path)
	a.router.HandleFunc(a.prefix+path, f).Methods("PUT")
}

// AttachModuleDelete registers a new DELETE handler for the API
func (a *API) AttachModuleDelete(path string, f func(http.ResponseWriter, *http.Request)) {
	a.log.Infoln("Registering DELETE handler:", a.prefix+path)
	a.router.HandleFunc(a.prefix+path, f).Methods("DELETE")
}

func (a *API) run() {
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Fatalf("Could not start http server: %v\n", err)
	}
}

// Init the API instance
func (a *API) Init() error {
	var listenAddress string
	if len(a.config.IP) > 0 && len(a.config.Port) > 0 {
		listenAddress = a.config.IP + ":" + a.config.Port
	} else if len(a.config.Port) > 0 {
		listenAddress = ":" + a.config.Port
	} else {
		return fmt.Errorf("REST API activated but no valid configuration found. At least port has to specified")
	}

	a.log.Infof("Configured REST API on %s", listenAddress)

	a.server = &http.Server{
		Addr:    listenAddress,
		Handler: handlers.CORS()(a.router),
	}

	return nil
}

// Run the REST API (blocking)
func (a *API) Run(ctx context.Context) error {
	if a.server == nil {
		a.Init()
	}
	go a.run()

	a.log.Printf("REST API server started")

	<-ctx.Done()

	a.log.Printf("REST API server shutdown")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	var err error
	if err = a.server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("REST API server shutdown Failed:%+s", err)
	}

	a.log.Printf("REST API server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}
