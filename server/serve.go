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

func ServeFromGolog(golog *gologme.Golog, url string) {
	server := NewServer()
	router := mux.NewRouter()
	router.Handle("/rpc", server)
	// Has to happen after rpc router is registered
	RegisterRoutes(router)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(url, router))
}

func ServeFromPath(dbPath string, url string) {
	gologme := gologme.NewGolog(dbPath)
	x, _ := gologme.DS.FindUserNameById(1)
	println(x)
	ServeFromGolog(gologme, url)
}
