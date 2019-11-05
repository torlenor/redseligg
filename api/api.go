package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
)

// API represents a Rest API instance of a ALoLStats instance
type API struct {
	config config.API
	router *mux.Router
	log    *logrus.Entry
	prefix string

	server *http.Server
	quit   chan interface{}
}

// NewAPI creates a new Rest API instance
func NewAPI(cfg config.API) (*API, error) {
	a := &API{
		config: cfg,
		router: mux.NewRouter(),
		log:    logging.Get("RestAPI"),
		prefix: "/v1",
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
func (a *API) Init() {
	var listenAddress string
	if len(a.config.IP) > 0 && len(a.config.Port) > 0 {
		listenAddress = a.config.IP + ":" + a.config.Port
	} else if len(a.config.Port) > 0 {
		listenAddress = ":" + a.config.Port
	} else {
		a.log.Fatal("REST API activated but no valid configuration found. At least port has to specified in config file!")
	}

	a.log.Infof("Starting REST API on %s", listenAddress)

	a.server = &http.Server{
		Addr:    listenAddress,
		Handler: handlers.CORS()(a.router),
	}
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
