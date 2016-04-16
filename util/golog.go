package gologme

import (
	"fmt"
	"log"
	"time"

	"github.com/erasche/gologme/store"
	gologme_types "github.com/erasche/gologme/types"
)

type Golog struct {
	DS store.DataStore
}

func (t *Golog) LogToDb(uid int, windowlogs []*gologme_types.WindowLogs, keylogs []*gologme_types.KeyLogs) {
	t.DS.LogToDb(uid, windowlogs, keylogs)
}

func (t *Golog) ExportEventsByDate(tm time.Time) *gologme_types.EventLog {
	fmt.Printf("%#v\n", t)
	return t.DS.ExportEventsByDate(tm)
}

func (t *Golog) Log(args *gologme_types.DataLogRequest) int {
	log.Printf("golog.Log\n")
	uid, err := t.DS.CheckAuth(args.User, args.ApiKey)
	if err != nil {
		log.Fatal(err)
		return 1
	} else {
		log.Printf("%s authenticated successfully as uid %d\n", args.User, uid)
	}

	t.LogToDb(
		uid,
		args.Windows,
		args.KeyLogs,
	)
	return 0
}

func NewGolog(fn string) *Golog {
	datastore, err := store.CreateDataStore(map[string]string{
		"DATASTORE":      "sqlite3",
		"DATASTORE_PATH": fn,
	})
	if err != nil {
		log.Fatal(err)
	}
	x := &Golog{
		DS: datastore,
	}
	return x
}
