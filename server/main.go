// Example get-active-window reads the _NET_ACTIVE_WINDOW property of the root
// window and uses the result (a window id) to get the name of the window.
package main

import (
	"fmt"
	"github.com/erasche/gologme"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Golog int

func (t *Golog) Log(args gologme.RpcArgs, result *gologme.Result) error {
	return Log(args, result)
}

func Log(args gologme.RpcArgs, result *gologme.Result) error {
	log.Printf("%d logs from [%d]\n", args.Length, args.User)
	for i, w := range args.Windows {
		log.Printf("\t%d %s\n", i, w)
		if i >= args.Length - 1{
			break
		}
	}
	*result = 0
	return nil
}

func main() {
	golog := new(Golog)
	rpc.Register(golog)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":8080")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("Listening...")
	err := http.Serve(l, nil)
	if err != nil {
		log.Fatal("Error serving:", err)
	}
}
