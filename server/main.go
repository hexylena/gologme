// Example get-active-window reads the _NET_ACTIVE_WINDOW property of the root
// window and uses the result (a window id) to get the name of the window.
package main

import (
	"database/sql"
	"fmt"
	"github.com/erasche/gologme"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Golog struct {
	Db *sql.DB
}

func (t *Golog) Log(args gologme.RpcArgs, result *gologme.Result) error {
	tx, err := t.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into windowLogs (time, name) values (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	log.Printf("%d logs from [%d]\n", args.Length, args.User)
	for i, w := range args.Windows {
		_, err = stmt.Exec(w.Time.Unix(), w.Name)
		if err != nil {
			log.Fatal(err)
		}

		if i >= args.Length-1 {
			break
		}
	}

	tx.Commit()
	*result = 0
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "file.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbCreation := `
	create table if not exists windowLogs (
		id integer not null primary key,
		time integer,
		name text
	)
	`

	_, err = db.Exec(dbCreation)
	if err != nil {
		log.Fatal(err)
	}

	golog := new(Golog)
	golog.Db = db

	rpc.Register(golog)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	fmt.Println("Listening...")

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Error serving:", err)
	}
}
