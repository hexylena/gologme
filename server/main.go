package main

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

func main() {
	db, err := sql.Open("sqlite3", "file.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	golog = new(gologme.Golog)
	golog.SetupDb(db)

	server := NewServer()
	router := mux.NewRouter()
	router.Handle("/rpc", server)
	// Has to happen after rpc router is registered
	RegisterRoutes(router)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
