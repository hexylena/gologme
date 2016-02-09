package server

import (
	"fmt"
	gologme "github.com/erasche/gologme/util"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var golog *gologme.Golog

func ServeFromGolog(g *gologme.Golog, url string) {
	golog = g
	server := NewServer()
	router := mux.NewRouter()
	router.Handle("/rpc", server)
	// Has to happen after rpc router is registered
	RegisterRoutes(router)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(url, router))
}

func ServeFromPath(dbPath string, url string) {
	g := gologme.NewGolog(dbPath)
	fmt.Printf("%#v\n", g)
	ServeFromGolog(g, url)
}
