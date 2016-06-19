package store

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	gologme "github.com/erasche/gologme/types"
	_ "github.com/lib/pq" // postgres module
)

//PostgreSQLDataStore struct
type PostgreSQLDataStore struct {
	DSN string
	DB  *sql.DB
}

// SetupDb runs any migrations needed
func (ds *PostgreSQLDataStore) SetupDb() {
	_, err := ds.DB.Exec(
		strings.Replace(DB_SCHEMA, "id integer not null primary key autoincrement", "id serial not null primary key", -1),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// LogToDb logs a set of windowLogs and keyLogs to the DB
func (ds *PostgreSQLDataStore) LogToDb(uid int, windowlogs []*gologme.WindowLogs, keylogs []*gologme.KeyLogs) {
	tx, err := ds.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	wlStmt, err := tx.Prepare("insert into windowLogs (uid, time, name) values ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	klStmt, err := tx.Prepare("insert into keyLogs (uid, time, count) values ($1, $2, $3)")
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
func (ds *PostgreSQLDataStore) CheckAuth(user string, key string) (int, error) {
	// Pretty assuredly not safe from timing attacks.
	var id int
	err := ds.DB.QueryRow("SELECT id FROM users WHERE username = $1 AND api_key = $2", user, key).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1, ErrUserNotFound
		}
		return -1, err
	}
	return id, nil
}

func (ds *PostgreSQLDataStore) CreateNote(uid int, date time.Time, message string) {
	tx, err := ds.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	noteInsert, err := tx.Prepare("insert into notes (uid, time, type, contents) values ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}

	defer noteInsert.Close()

	_, err = noteInsert.Exec(uid, date.Unix(), gologme.NOTE_TYPE, message)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	return
}

func (ds *PostgreSQLDataStore) CreateBlog(uid int, date time.Time, message string) {
	tx, err := ds.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	noteInsert, err := tx.Prepare("insert into notes (uid, time, type, contents) values ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}

	defer noteInsert.Close()

	_, err = noteInsert.Exec(uid, date.Unix(), gologme.BLOG_TYPE, message)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	return
}

// Name of the DS implementatino
func (ds *PostgreSQLDataStore) Name() string {
	return "PostgreSQLDataStore"
}

// MaxDate returns latest log entry
func (ds *PostgreSQLDataStore) MaxDate() int {
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
func (ds *PostgreSQLDataStore) MinDate() int {
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
func (ds *PostgreSQLDataStore) FindUserNameByID(id int) (string, error) {
	var username string
	err := ds.DB.QueryRow("SELECT username FROM users WHERE id = $1", id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return username, nil
}

func (ds *PostgreSQLDataStore) ExportWindowLogsByRange(t0 int64, t1 int64) []*gologme.SEvent {
	stmt, err := ds.DB.Prepare("select time, name from windowLogs where time >= $1 and time < $2 order by id")
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

func (ds *PostgreSQLDataStore) exportKeyLogsByRange(t0 int64, t1 int64) []*gologme.IEvent {
	stmt, err := ds.DB.Prepare("select time, count from keyLogs where time >= $1 and time < $2 order by id")
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

func (ds *PostgreSQLDataStore) exportBlog(t0 int64, t1 int64) []*gologme.SEvent {
	stmt, err := ds.DB.Prepare("select time, type, contents from notes where time >= $1 and time < $2 and type = $3 order by id")
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

func (ds *PostgreSQLDataStore) exportNotes(t0 int64, t1 int64) []*gologme.SEvent {
	stmt, err := ds.DB.Prepare("select time, type, contents from notes where time >= $1 and time < $2 and type = $3 order by id")
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
func (ds *PostgreSQLDataStore) ExportEventsByDate(tm time.Time) *gologme.EventLog {
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
		Window_events:  ds.ExportWindowLogsByRange(t0, t1),
		Keyfreq_events: ds.exportKeyLogsByRange(t0, t1),
		Note_events:    ds.exportNotes(t0, t1),
		Blog:           blogstr,
	}
}

// NewPostgreSQLDataStore builds a new PG DS
func NewPostgreSQLDataStore(conf map[string]string) (DataStore, error) {
	var dsn string
	if val, ok := conf["DATASTORE_URL"]; ok {
		dsn = val
	} else {
		return nil, fmt.Errorf("%s is required for the postgres datastore", "DATASTORE_URL")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
		return nil, ErrFailedToConnect
	}

	return &PostgreSQLDataStore{
		DSN: dsn,
		DB:  db,
	}, nil
}
