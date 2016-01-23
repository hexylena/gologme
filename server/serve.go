package server

import (
	"database/sql"
	"fmt"
	gologme "github.com/erasche/gologme/util"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var golog *gologme.Golog

func ServeFromDb(db *sql.DB, url string){
	golog = new(gologme.Golog)
	golog.SetupDb(db)

	server := NewServer()
	router := mux.NewRouter()
	router.Handle("/rpc", server)
	// Has to happen after rpc router is registered
	RegisterRoutes(router)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(url, router))
}

func ServeFromPath(dbPath string, url string){
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	go ServeFromDb(db, url)
}
