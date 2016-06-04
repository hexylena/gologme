package server

import (
	"fmt"
	"log"
	"net/http"

	gologme "github.com/erasche/gologme/util"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var golog *gologme.Golog

func ServeFromGolog(g *gologme.Golog, url string) {
	golog = g

	router := mux.NewRouter()
	RegisterRoutes(router)

	// Listen and serve
	fmt.Printf("Listening on %s ...\n", url)
	log.Fatal(http.ListenAndServe(url, handlers.CORS()(router)))
}

func ServeFromPath(dbType string, dbPath string, url string) {
	g := gologme.NewGolog(dbType, dbPath)
	ServeFromGolog(g, url)
}
