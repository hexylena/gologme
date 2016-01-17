package main

import (
	"log"
	"net/http"
)

type RpcLog int

type Args struct {
	A, B int
}

func (t *RpcLog) TestRpc(r *http.Request, args *Args, result *int) error {
	log.Printf("%s\n", golog)

	*result = 0
	return nil
}
