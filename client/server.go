package client

import (
	"github.com/erasche/gologme/server"
)

func Serve(path string, listenAddr string) {
	server.ServeFromPath(path, listenAddr)
}
