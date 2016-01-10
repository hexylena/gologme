package main

import (
    "fmt"
	"log"
	"net/rpc"
	"github.com/erasche/gologme"
)

func send(wl []gologme.WindowLogs, wi int){
    send_local(wl, wi)
    //send_remote(wl, wi)
}

func send_local(wl []gologme.WindowLogs, wi int){
    for i, w := range(wl){
        fmt.Printf("%s\n", w)
        if i >= wi-1 {
            break
        }
    }
}

func send_remote(wl []gologme.WindowLogs, wi int) {
	client, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatal("Error in dialing", err)
	}
	args := &gologme.RpcArgs{
		User:    0,
		Windows: wl,
		Length:  wi,
	}
	var result gologme.Result
	err = client.Call("Golog.Log", args, &result)
	if err != nil {
		log.Fatal("Error calling RPC method", err)
	}
}

