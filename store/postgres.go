package store

import (
	"database/sql"
	"errors"
	"fmt"
	gologme "github.com/erasche/gologme/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

//The first implementation.
type PostgreSQLDataStore struct {
	DSN string
	DB  *sqlx.DB
}

func (ds *PostgreSQLDataStore) SetupDb() {
	_, err := ds.DB.Exec(DB_SCHEMA)
	if err != nil {
		log.Fatal(err)
	}
}

func (pds *PostgreSQLDataStore) LogToDb(uid int, windowlogs []gologme.WindowLogs, keylogs []gologme.KeyLogs, wll int) {
	tx, err := pds.DB.Begin()
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

func (pds *PostgreSQLDataStore) CheckAuth(user string, key string) (int, error) {
	// Pretty assuredly not safe from timing attacks.
	stmt, err := pds.DB.Prepare("select id from users where username = ? and api_key = ?")
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

func (pds *PostgreSQLDataStore) Name() string {
	return "PostgreSQLDataStore"
}

func (pds *PostgreSQLDataStore) FindUserNameById(id int) (string, error) {
	var username string
	res, err := pds.DB.Query("SELECT username FROM users WHERE id=$1", id)
	res.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", UserNotFoundError
		}
		return "", err
	}
	return username, nil
}

func (pds *PostgreSQLDataStore) ExportEventsByDate(tm time.Time) *gologme.EventLog {
	return nil
}

func NewPostgreSQLDataStore(conf map[string]string) (DataStore, error) {
	var dsn string
	if val, ok := conf["DATASTORE_POSTGRES_DSN"]; ok {
		dsn = val
	} else {
		return nil, errors.New(fmt.Sprintf("%s is required for the postgres datastore", "DATASTORE_POSTGRES_DSN"))
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
		return nil, FailedToConnect
	}

	return &PostgreSQLDataStore{
		DSN: dsn,
		DB:  db,
	}, nil
}
