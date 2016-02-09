package gologme

import (
	"github.com/erasche/gologme/store"
	gologme_types "github.com/erasche/gologme/types"
	"log"
	"os/user"
	"path"
	"time"
)

type Golog struct {
	DS store.DataStore
}

func (t *Golog) LogToDb(uid int, windowlogs []gologme_types.WindowLogs, keylogs []gologme_types.KeyLogs, wll int) {
	t.DS.LogToDb(uid, windowlogs, keylogs, wll)
}

func (t *Golog) ExportEventsByDate(tm time.Time) *gologme_types.EventLog {
	return t.DS.ExportEventsByDate(tm)
}

func (t *Golog) Log(args gologme_types.RpcArgs) int {
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
		args.WindowLogsLength,
	)
	return 0
}

func NewGolog(fn string) *Golog {
	if len(fn) == 0 {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		fn = path.Join(user.HomeDir, ".gologme.db")
	}

	datastore, err := store.CreateDataStore(map[string]string{
		"DATASTORE":      "sqlite3",
		"DATASTORE_PATH": fn,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Golog{
		DS: datastore,
	}
}
