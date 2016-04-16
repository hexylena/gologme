package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"Events", "GET", "/api/events/{date:[0-9]+}", Events},
	Route{"DataUpload", "POST", "/logs", DataUpload},
}

func RegisterRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	router.PathPrefix("/").Handler(http.FileServer(assetFS()))

	return router
}
