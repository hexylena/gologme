package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	gologme "github.com/erasche/gologme/types"
	_ "github.com/lib/pq"
)

//The first implementation.
type PostgreSQLDataStore struct {
	DSN string
	DB  *sql.DB
}

func (ds *PostgreSQLDataStore) SetupDb() {
	_, err := ds.DB.Exec(
		strings.Replace(DB_SCHEMA, "id integer not null primary key autoincrement", "id serial not null primary key", -1),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (ds *PostgreSQLDataStore) LogToDb(uid int, windowlogs []*gologme.WindowLogs, keylogs []*gologme.KeyLogs) {
	tx, err := ds.DB.Begin()
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

	wll := len(windowlogs)
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

func (ds *PostgreSQLDataStore) CheckAuth(user string, key string) (int, error) {
	// Pretty assuredly not safe from timing attacks.
	var id int
	err := ds.DB.QueryRow("SELECT id FROM users WHERE username = ? AND api_key = ?", user, key).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1, UserNotFoundError
		}
		return -1, err
	}
	return id, nil
}

func (ds *PostgreSQLDataStore) Name() string {
	return "PostgreSQLDataStore"
}

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

func (ds *PostgreSQLDataStore) FindUserNameById(id int) (string, error) {
	var username string
	err := ds.DB.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", UserNotFoundError
		}
		return "", err
	}
	return username, nil
}

func (ds *PostgreSQLDataStore) exportWindowLogsByRange(t0 int64, t1 int64) []*gologme.SEvent {
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

func (ds *PostgreSQLDataStore) exportKeyLogsByRange(t0 int64, t1 int64) []*gologme.IEvent {
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

func (ds *PostgreSQLDataStore) exportBlog(t0 int64, t1 int64) []*gologme.SEvent {
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

func (ds *PostgreSQLDataStore) exportNotes(t0 int64, t1 int64) []*gologme.SEvent {
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
		Window_events:  ds.exportWindowLogsByRange(t0, t1),
		Keyfreq_events: ds.exportKeyLogsByRange(t0, t1),
		Note_events:    ds.exportNotes(t0, t1),
		Blog:           blogstr,
	}
}

func NewPostgreSQLDataStore(conf map[string]string) (DataStore, error) {
	var dsn string
	if val, ok := conf["DATASTORE_URL"]; ok {
		dsn = val
	} else {
		return nil, errors.New(fmt.Sprintf("%s is required for the postgres datastore", "DATASTORE_URL"))
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
		return nil, FailedToConnect
	}

	return &PostgreSQLDataStore{
		DSN: dsn,
		DB:  db,
	}, nil
}
