/**
* api.go - REST API base implementation
*
 */

package api

import (
	"io"
	"logging"
	"net/http"

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
func Start() {
	log.Info("Starting up REST API")

	attachPublic(router)

	log.Fatal(http.ListenAndServe(":8000", router))
}
