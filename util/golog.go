package gologme

import (
	"database/sql"
    //"fmt"
	"log"
	"time"
	//"github.com/erasche/gologme/ulogme"
)

type Golog struct {
	Db *sql.DB
}

func (t *Golog) LogToDb(uid int, windowlogs []WindowLogs, keylogs []KeyLogs, wll int) {
	tx, err := t.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	wl_stmt, err := tx.Prepare("insert into windowLogs (uid, time, name) values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	kl_stmt, err := tx.Prepare("insert into keyLogs (uid, time, count) values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer wl_stmt.Close()
	defer kl_stmt.Close()

	log.Printf("%d window logs %d key logs from [%d]\n", wll, len(keylogs), uid)

	for _, w := range keylogs {
		_, err = kl_stmt.Exec(uid, w.Time.Unix(), w.Count)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, w := range windowlogs {
		_, err = wl_stmt.Exec(uid, w.Time.Unix(), w.Name)
		if err != nil {
			log.Fatal(err)
		}
		// Dunno why windowLogs comes through two too big, so whatever.
		if i >= wll-1 {
			break
		}
	}

	tx.Commit()
}

func (t *Golog) ensureAuth(user string, key string) (int, error) {
	// Pretty assuredly not safe from timing attacks.
	stmt, err := t.Db.Prepare("select id from users where username = ? and api_key = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var uid int
	err = stmt.QueryRow(user, key).Scan(&uid)
	if err != nil {
		return -1, err
	}

	return uid, nil
}

func (t *Golog) Log(args RpcArgs) int {
	uid, err := t.ensureAuth(args.User, args.ApiKey)
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

func (t *Golog) SetupDb(db *sql.DB) {
	_, err := db.Exec(DB_SCHEMA)
	if err != nil {
		log.Fatal(err)
	}
	t.Db = db
}

type SEvent struct {
	T int    `json:"t"`
	S string `json:"s"`
}

type IEvent struct {
	T int `json:"t"`
	S int `json:"s"`
}

type SEventT struct {
	RealTime time.Time
	Title    string
}

type IEventT struct {
	RealTime time.Time
	Count    int
}

type NoteEvent struct {
	RealTime time.Time
	Type     int
	Contents string
}

type EventLog struct {
	Blog           string    `json:"blog"`
	Note_events    []*SEvent `json:"notes_events"`
	Window_events  []*SEvent `json:"window_events"`
	Keyfreq_events []*IEvent `json:"keyfreq_events"`
}

func (t *Golog) exportWindowLogsByRange(t0 int64, t1 int64) []*SEvent {
	stmt, err := t.Db.Prepare("select time, name from windowLogs where time >= ? and time < ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*SEvent, 0)
	for rows.Next() {
		var (
			t int
			s string
		)
		rows.Scan(&t, &s)

		logs = append(
			logs,
			&SEvent{
				T: t,
				S: s,
			},
		)
	}
	return logs
}

func (t *Golog) exportKeyLogsByRange(t0 int64, t1 int64) []*IEvent {
	stmt, err := t.Db.Prepare("select time, count from keyLogs where time >= ? and time < ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*IEvent, 0)
	for rows.Next() {
		var (
			t int
			s int
		)
		rows.Scan(&t, &s)

		logs = append(
			logs,
			&IEvent{
				T: t,
				S: s,
			},
		)
	}
	return logs
}

func (t *Golog) exportBlog(t0 int64, t1 int64) []*SEvent {
	stmt, err := t.Db.Prepare("select time, type, contents from notes where time >= ? and time < ? and type = ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1, BLOG_TYPE)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*SEvent, 0)
	for rows.Next() {
		var (
			Time     int
			Type     int
			Contents string
		)
		rows.Scan(&Time, &Type, &Contents)
		logs = append(
			logs,
			&SEvent{
				T: Time,
				S: Contents,
			},
		)
	}
	return logs
}

func (t *Golog) exportNotes(t0 int64, t1 int64) []*SEvent {
	stmt, err := t.Db.Prepare("select time, type, contents from notes where time >= ? and time < ? and type = ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1, NOTE_TYPE)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*SEvent, 0)
	for rows.Next() {
		var (
			Time     int
			Type     int
			Contents string
		)
		rows.Scan(&Time, &Type, &Contents)
		logs = append(
			logs,
			&SEvent{
				T: Time,
				S: Contents,
			},
		)
	}
	return logs
}

func (t *Golog) ExportEventsByDate(tm time.Time) *EventLog {
	t0 := Ulogme7amTime(tm)
	t1 := Ulogme7amTime(Tomorrow(tm))

	blog := t.exportBlog(t0, t1)
	var blogstr string
	if len(blog) > 0 {
		blogstr = blog[0].S
	} else {
		blogstr = ""
	}

	return &EventLog{
		Window_events:  t.exportWindowLogsByRange(t0, t1),
		Keyfreq_events: t.exportKeyLogsByRange(t0, t1),
		Note_events:    t.exportNotes(t0, t1),
		Blog:           blogstr,
	}
}
