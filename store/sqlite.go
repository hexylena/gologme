package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	gologme "github.com/erasche/gologme/types"
	_ "github.com/mattn/go-sqlite3" // sqlite module
)

// SqliteSQLDataStore struct
type SqliteSQLDataStore struct {
	DSN string
	DB  *sql.DB
}

// SetupDb runs any migrations needed
func (ds *SqliteSQLDataStore) SetupDb() {
	_, err := ds.DB.Exec(DB_SCHEMA)
	if err != nil {
		log.Fatal(err)
	}
}

// LogToDb logs a set of windowLogs and keyLogs to the DB
func (ds *SqliteSQLDataStore) LogToDb(uid int, windowlogs []*gologme.WindowLogs, keylogs []*gologme.KeyLogs) {
	tx, err := ds.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	wlStmt, err := tx.Prepare("insert into windowLogs (uid, time, name) values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	klStmt, err := tx.Prepare("insert into keyLogs (uid, time, count) values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer wlStmt.Close()
	defer klStmt.Close()

	wll := len(windowlogs)
	log.Printf("%d window logs %d key logs from [%d]\n", wll, len(keylogs), uid)

	for _, w := range keylogs {
		_, err = klStmt.Exec(uid, w.Time.Unix(), w.Count)
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, w := range windowlogs {
		_, err = wlStmt.Exec(uid, w.Time.Unix(), w.Name)
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

// CheckAuth of the user+key, returning -1 or the user's ID
func (ds *SqliteSQLDataStore) CheckAuth(user string, key string) (int, error) {
	// Pretty assuredly not safe from timing attacks.
	stmt, err := ds.DB.Prepare("select id from users where username = ? and api_key = ?")
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

// Name of the DS implementatino
func (ds *SqliteSQLDataStore) Name() string {
	return "SqliteSQLDataStore"
}

// MaxDate returns latest log entry
func (ds *SqliteSQLDataStore) MaxDate() int {
	var mtime int
	err := ds.DB.QueryRow("SELECT time FROM windowLogs ORDER BY time DESC LIMIT 1").Scan(&mtime)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No data available")
	case err != nil:
		log.Fatal(err)
	}
	return mtime
}

// MinDate returns first log entry
func (ds *SqliteSQLDataStore) MinDate() int {
	var mtime int
	err := ds.DB.QueryRow("SELECT time FROM windowLogs ORDER BY time ASC LIMIT 1").Scan(&mtime)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No data available")
	case err != nil:
		log.Fatal(err)
	}
	return mtime
}

// FindUserNameByID returns a username for a given ID
func (ds *SqliteSQLDataStore) FindUserNameByID(id int) (string, error) {
	var username string
	err := ds.DB.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return username, nil
}

func (ds *SqliteSQLDataStore) exportWindowLogsByRange(t0 int64, t1 int64) []*gologme.SEvent {
	stmt, err := ds.DB.Prepare("select time, name from windowLogs where time >= ? and time < ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*gologme.SEvent, 0)
	for rows.Next() {
		var (
			t int
			s string
		)
		rows.Scan(&t, &s)

		logs = append(
			logs,
			&gologme.SEvent{
				T: t,
				S: s,
			},
		)
	}
	return logs
}

func (ds *SqliteSQLDataStore) exportKeyLogsByRange(t0 int64, t1 int64) []*gologme.IEvent {
	stmt, err := ds.DB.Prepare("select time, count from keyLogs where time >= ? and time < ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*gologme.IEvent, 0)
	for rows.Next() {
		var (
			t int
			s int
		)
		rows.Scan(&t, &s)

		logs = append(
			logs,
			&gologme.IEvent{
				T: t,
				S: s,
			},
		)
	}
	return logs
}

func (ds *SqliteSQLDataStore) exportBlog(t0 int64, t1 int64) []*gologme.SEvent {
	stmt, err := ds.DB.Prepare("select time, type, contents from notes where time >= ? and time < ? and type = ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1, gologme.BLOG_TYPE)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*gologme.SEvent, 0)
	for rows.Next() {
		var (
			Time     int
			Type     int
			Contents string
		)
		rows.Scan(&Time, &Type, &Contents)
		logs = append(
			logs,
			&gologme.SEvent{
				T: Time,
				S: Contents,
			},
		)
	}
	return logs
}

func (ds *SqliteSQLDataStore) exportNotes(t0 int64, t1 int64) []*gologme.SEvent {
	stmt, err := ds.DB.Prepare("select time, type, contents from notes where time >= ? and time < ? and type = ? order by id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(t0, t1, gologme.NOTE_TYPE)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]*gologme.SEvent, 0)
	for rows.Next() {
		var (
			Time     int
			Type     int
			Contents string
		)
		rows.Scan(&Time, &Type, &Contents)
		logs = append(
			logs,
			&gologme.SEvent{
				T: Time,
				S: Contents,
			},
		)
	}
	return logs
}

// ExportEventsByDate extracts events for a given day
func (ds *SqliteSQLDataStore) ExportEventsByDate(tm time.Time) *gologme.EventLog {
	t0 := Ulogme7amTime(tm)
	t1 := Ulogme7amTime(Tomorrow(tm))

	blog := ds.exportBlog(t0, t1)
	var blogstr string
	if len(blog) > 0 {
		blogstr = blog[0].S
	} else {
		blogstr = ""
	}

	return &gologme.EventLog{
		Window_events:  ds.exportWindowLogsByRange(t0, t1),
		Keyfreq_events: ds.exportKeyLogsByRange(t0, t1),
		Note_events:    ds.exportNotes(t0, t1),
		Blog:           blogstr,
	}
}

// NewSqliteSQLDataStore builds a new sqlite3 DS
func NewSqliteSQLDataStore(conf map[string]string) (DataStore, error) {
	var dsn string
	if val, ok := conf["DATASTORE_URL"]; ok {
		dsn = val
	} else {
		return nil, errors.New(fmt.Sprintf("%s is required for the sqlite datastore", "DATASTORE_URL"))
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
		return nil, ErrFailedToConnect
	}

	return &SqliteSQLDataStore{
		DSN: dsn,
		DB:  db,
	}, nil
}
