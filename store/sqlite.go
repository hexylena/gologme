package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"errors"
	"fmt"
	"github.com/erasche/gologme/util"
	"log"
)

//The first implementation.
type SqliteSQLDataStore struct {
	DSN string
	DB  *sql.DB
}

func (ds *SqliteSQLDataStore) LogToDb(uid int, windowlogs []gologme.WindowLogs, keylogs []gologme.KeyLogs, wll int) {
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

func (ds *SqliteSQLDataStore) Name() string {
	return "SqliteSQLDataStore"
}

func (ds *SqliteSQLDataStore) FindUserNameById(id int64) (string, error) {
	var username string
	res, err := ds.DB.Query("SELECT username FROM users WHERE id=$1", id)
	res.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", UserNotFoundError
		}
		return "", err
	}
	return username, nil
}

func NewSqliteSQLDataStore(conf map[string]string) (DataStore, error) {
	var dsn string
	if dsn, ok := conf["DATASTORE_PATH"]; !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the sqlite datastore", "DATASTORE_PATH"))
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
		return nil, FailedToConnect
	}

	return &SqliteSQLDataStore{
		DSN: dsn,
		DB:  db,
	}, nil
}

