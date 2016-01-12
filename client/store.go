package main

import (
	"fmt"
	"github.com/erasche/gologme"
	"net/rpc"
)

func send(wl []gologme.WindowLogs, wi int, kl []gologme.KeyLogs) {
	send_remote(wl, kl, wi)
}


func send_remote(wl []gologme.WindowLogs, kl []gologme.KeyLogs, wi int) {
	client, err := rpc.DialHTTP("tcp", ":10000")
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
