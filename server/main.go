// Example get-active-window reads the _NET_ACTIVE_WINDOW property of the root
// window and uses the result (a window id) to get the name of the window.
package main

import (
	"database/sql"
	"fmt"
	"github.com/erasche/gologme"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	db, err := sql.Open("sqlite3", "file.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	golog := new(gologme.Golog)
	golog.SetupDb(db)
	rpc.Register(golog)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	fmt.Println("Listening...")
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Error serving:", err)
	}
}
