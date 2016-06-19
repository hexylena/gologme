package server

import (
	"fmt"
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

var (
	xAuthorization = http.CanonicalHeaderKey("Authorization")
)

func authenticationHandler(fn func(w http.ResponseWriter, r *http.Request, uid int)) http.HandlerFunc {
	//, username string, APIkey string
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("asdf")
		api_key := r.Header.Get(xAuthorization)
		if len(api_key) == 0 {
			return
		}
		uid, err := golog.Authenticate(api_key)
		if err != nil {
			return
		}
		fn(w, r, uid)
	}
}

var routes = Routes{
	Route{"Events", "GET", "/api/events/{date:[0-9-]+}", authenticationHandler(Events)},
	Route{"RecentWindows", "GET", "/api/events/recent", authenticationHandler(RecentWindows)},
	Route{"KeyEvents", "GET", "/api/events/key/{date:[0-9-]+}", authenticationHandler(KeyEvents)},
	Route{"WinEvents", "GET", "/api/events/win/{date:[0-9-]+}", authenticationHandler(WinEvents)},
	Route{"AddNote", "POST", "/api/notes", authenticationHandler(AddNote)},
	Route{"AddBlog", "POST", "/api/blog", authenticationHandler(AddBlog)},
	Route{"DataUpload", "POST", "/logs", DataUpload},
	Route{"ExportList", "GET", "/export_list.json", authenticationHandler(ExportList)},
}

// RegisterRoutes operates over `Routes` and registers all of them
func RegisterRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	router.PathPrefix("/").Handler(http.FileServer(assetFS()))
	return router
}
