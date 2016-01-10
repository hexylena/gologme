package main

import (
	"fmt"
	"github.com/erasche/gologme"
	"net/rpc"
)

func send(wl []gologme.WindowLogs, wi int, kl []gologme.KeyLogs) {
	//send_local(wl, wi, kl)
	//for i, e := range kl {
	//fmt.Printf("KL: %d - %s\n", i, e)
	//}
	send_remote(wl, kl, wi)
}

func send_local(wl []gologme.WindowLogs, kl []gologme.KeyLogs, wi int) {
	for i, w := range wl {
		fmt.Printf("WL: %s\n", w)
		if i >= wi-1 {
			break
		}
	}
}

func send_remote(wl []gologme.WindowLogs, kl []gologme.KeyLogs, wi int) {
	client, err := rpc.DialHTTP("tcp", ":8080")
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
	var result gologme.Result
	err = client.Call("Golog.Log", args, &result)
	if err != nil {
		fmt.Printf("Error in calling RPC method, droping logs, %s\n", err)
        return
		// TODO: retry
	}
}
