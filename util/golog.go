package gologme

import (
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

func (t *Golog) CreateBlog(username string, api_key string, date time.Time, message string) error {
	uid, err := t.DS.CheckAuth(username, api_key)
	if err != nil {
		return err
	}
	t.DS.CreateBlog(uid, date, message)
	return nil
}

func (t *Golog) CreateNote(username string, api_key string, date time.Time, message string) error {
	uid, err := t.DS.CheckAuth(username, api_key)
	if err != nil {
		return err
	}
	t.DS.CreateNote(uid, date, message)
	return nil
}

func (t *Golog) ExportEventsByDate(tm time.Time) *gologme_types.EventLog {
	return t.DS.ExportEventsByDate(tm)
}

func (t *Golog) ExportWindowLogsByRange(t1, t2 int64) []*gologme_types.SEvent {
	return t.DS.ExportWindowLogsByRange(t1, t2)
}

func (t *Golog) RecordedDataRange() (time.Time, time.Time) {
	minTime := time.Unix(int64(t.DS.MinDate()), 0)
	maxTime := time.Unix(int64(t.DS.MaxDate()), 0)
	return minTime, maxTime
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

func NewGolog(typ, fn string) *Golog {
	datastore, err := store.CreateDataStore(map[string]string{
		"DATASTORE":     typ,
		"DATASTORE_URL": fn,
	})
	if err != nil {
		log.Fatal(err)
	}
	x := &Golog{
		DS: datastore,
	}
	return x
}
