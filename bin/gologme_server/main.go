package main

import (
	"github.com/erasche/gologme/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	server.ServeFromPath("file.db", ":8080")
}
