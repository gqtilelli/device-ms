package handler

import (
	"net/http"

	"github.com/device-ms/controller"
	"github.com/gorilla/mux"
)

type deviceHandler struct {
	*mux.Router
	service controller.ServiceController
}

func (handler deviceHandler) addRoute(router *mux.Router, path, method string, f func(http.ResponseWriter, *http.Request)) {
	router.Path(path).Methods(method).HandlerFunc(f)
}

func addRoutes(router *mux.Router, handler deviceHandler) {
	handler.addRoute(router, "/{id}/name", http.MethodPut, handler.updateDeviceName)
	handler.addRoute(router, "/{id}/brand", http.MethodPut, handler.updateDeviceBrand)
	handler.addRoute(router, "/{id}", http.MethodGet, handler.getDevice)
	handler.addRoute(router, "/{id}", http.MethodPut, handler.updateDevice)
	handler.addRoute(router, "/{id}", http.MethodDelete, handler.deleteDevice)
	handler.addRoute(router, "", http.MethodPost, handler.createDevice)
	handler.addRoute(router, "", http.MethodGet, handler.getDevices)
}

func newDevice(service controller.ServiceController) deviceHandler {
	router := mux.NewRouter().PathPrefix(URLPath).Subrouter()
	handler := deviceHandler{
		Router:  router,
		service: service,
	}
	addRoutes(router, handler)
	return handler
}
