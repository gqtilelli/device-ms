package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/device-ms/controller"
	"github.com/gorilla/mux"
)

const (
	// URLPath Device resource base url
	URLPath = "/device"
)

type (
	// Router represents gorilla mux driver.
	Router struct {
		*mux.Router
	}
)

// HealthzHandler is a common healthz handler
func HealthzHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintln(w, "ok")
	if err != nil {
		return
	}
}

const (
	// DefaultRouteNotFoundErrMsg contains the default message of the route not found error.
	DefaultRouteNotFoundErrMsg = "resource not found"
)

// RouteNotFoundResponse represents route not found (404) response model.
type RouteNotFoundResponse struct {
	Error string `json:"error"`
	Path  string `json:"path"`
	URL   string `json:"url"`
}

func HandleNotFound(w http.ResponseWriter, req *http.Request) {
	resp := RouteNotFoundResponse{
		Error: DefaultRouteNotFoundErrMsg,
		Path:  req.URL.Path,
		URL:   req.URL.String(),
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("could not encode json return: %v, resp: %v", err, resp)
	}
}

// NewDeviceRouter creates a router for this microservice.
func NewDeviceRouter(service controller.ServiceController) Router {
	router := Router{
		Router: mux.NewRouter(),
	}
	router.HandleFunc("/heartbeat", HealthzHandler)
	router.PathPrefix(URLPath).Handler(newDevice(service))
	router.NotFoundHandler = http.HandlerFunc(HandleNotFound)

	return router
}
