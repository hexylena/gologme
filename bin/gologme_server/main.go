package main

import (
	_ "github.com/mattn/go-sqlite3"
    "github.com/erasche/gologme/server"
)

func main() {
	server.ServeFromPath("file.db", ":8080")
}
