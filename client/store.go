package main

import (
	"database/sql"
	"fmt"
	gologme "github.com/erasche/gologme/util"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/rpc"
	"os/user"
	"path"
)

func send(wl []gologme.WindowLogs, wi int, kl []gologme.KeyLogs, standalone bool) {
	if standalone {
		send_local(wl, kl, wi)
	} else {
		send_remote(wl, kl, wi)
	}
}

func send_local(wl []gologme.WindowLogs, kl []gologme.KeyLogs, wi int) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)

	}
	fn := path.Join(user.HomeDir, ".gologme.db")
	fmt.Printf("Storing to %s\n", fn)
	db, err := sql.Open("sqlite3", fn)
	// TODO: Ensure admin user?
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	golog := new(gologme.Golog)
	golog.SetupDb(db)
	golog.LogToDb(1, wl, kl, wi)
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
