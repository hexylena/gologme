package server

import (
	"fmt"
	"log"
	"net/http"

	gologme "github.com/erasche/gologme/util"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var golog *gologme.Golog

func ServeFromGolog(g *gologme.Golog, url string) {
	golog = g
	router := mux.NewRouter()
	// Has to happen after rpc router is registered
	RegisterRoutes(router)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(url, router))
}

func ServeFromPath(dbPath string, url string) {
	g := gologme.NewGolog(dbPath)
	ServeFromGolog(g, url)
}
