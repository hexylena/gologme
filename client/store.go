package client

import (
	"fmt"
	"net/rpc"

	gologme "github.com/erasche/gologme/types"
	_ "github.com/mattn/go-sqlite3"
)

type receiver struct {
    Standalone bool
    ServerAddress string
}

func (r *receiver) send(wl []gologme.WindowLogs, wi int, kl []gologme.KeyLogs) {
	client, err := rpc.DialHTTP("tcp", r.ServerAddress)
	if err != nil {
		fmt.Printf("Error in dialing, droping logs, %s\n", err)
		return
		// TODO: requeue
	}
	args := &gologme.RpcArgs{
		User:             "hxr",
		ApiKey:           "deadbeefcafe",
		Windows:          wl,
		KeyLogs:          kl,
		WindowLogsLength: wi,
	}
	var result int
	err = client.Call("Golog.Log", args, &result)
	if err != nil {
		fmt.Printf("Error in calling RPC method, droping logs, %s\n", err)
		return
		// TODO: retry
	}
}
