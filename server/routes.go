package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route struct
type Route struct {
	Name        string           // Route Name
	Method      string           // POST or GET
	Pattern     string           // URL
	HandlerFunc http.HandlerFunc // Handler function
}

// Routes list
type Routes []Route

var routes = Routes{
	Route{"Events", "GET", "/api/events/{date:[0-9-]+}", Events},
	Route{"DataUpload", "POST", "/logs", DataUpload},
	Route{"ExportList", "GET", "/export_list.json", ExportList},
}

// RegisterRoutes operates over `Routes` and registers all of them
func RegisterRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	router.PathPrefix("/").Handler(http.FileServer(assetFS()))
	return router
}
