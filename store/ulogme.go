package main

import (
	"database/sql"
	"errors"
	"fmt"
	ulogme "github.com/erasche/gologme/ulogme"
	"github.com/erasche/gologme/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

//The first implementation.
type UlogmeDataStore struct {
	LogDir string
}

func (uds *UlogmeDataStore) LogToDb(uid int, windowlogs []gologme.WindowLogs, keylogs []gologme.KeyLogs, wll int) {
	for _, w := range keylogs {
		ulogme.WriteIntFile(
			"keyfreq",
			key_logs,
			"x",
		)
	}

	for i, w := range windowlogs {
		ulogme.WriteStringFile(
			"window",
			window_logs,
			"x",
		)
		// Dunno why windowLogs comes through two too big, so whatever.
		if i >= wll-1 {
			break
		}
	}
}

func (uds *UlogmeDataStore) CheckAuth(user string, key string) (int, error) {
	return 1, nil
}

func (uds *UlogmeDataStore) Name() string {
	return "UlogmeDataStore"
}

func NewUlogmeDataStore(conf map[string]string) (DataStore, error) {
	var dsn string
	if val, ok := conf["DATASTORE_ULOGME_DSN"]; ok {
		dsn = val
	} else {
		return nil, errors.New(fmt.Sprintf("%s is required for the postgres datastore", "DATASTORE_POSTGRES_DSN"))
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
		return nil, FailedToConnect
	}

	return &UlogmeDataStore{
        LogDir: dsn,
	}, nil
}
