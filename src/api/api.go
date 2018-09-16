/**
* api.go - REST API base implementation
*
 */

package api

import (
	"net/http"

	"config"
	"io"
	"logging"

	"github.com/gorilla/mux"
)

var router *mux.Router
var log = logging.Get("api")

func init() {
	router = mux.NewRouter()
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func attachPublic(rtr *mux.Router) {
	AttachModuleGet("/ping", handlePing)
}

// AttachModuleGet registers a new GET handler for the API
func AttachModuleGet(path string, f func(http.ResponseWriter, *http.Request)) {
	log.Infoln("Registering GET handler:", path)
	router.HandleFunc(path, f).Methods("GET")
}

// AttachModulePost registers a new POST handler for the API
func AttachModulePost(path string, f func(http.ResponseWriter, *http.Request)) {
	log.Infoln("Registering POST handler:", path)
	router.HandleFunc(path, f).Methods("POST")
}

// Start the REST API
func Start(cfg config.API) {
	if cfg.Enabled == true {

		attachPublic(router)

		var listenAddress string
		if len(cfg.IP) > 0 && len(cfg.Port) > 0 {
			listenAddress = cfg.IP + ":" + cfg.Port
		} else if len(cfg.Port) > 0 {
			listenAddress = ":" + cfg.Port
		} else {
			log.Fatal("REST API activated but no valid configuration found. At least port has to specified in config file!")
		}

		log.Infof("REST API running on %s", listenAddress)
		log.Fatal(http.ListenAndServe(listenAddress, router))
	} else {
		log.Info("NOT starting REST API")
	}
}
