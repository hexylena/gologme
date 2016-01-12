package gologme

import (
	"database/sql"
	"log"
	"time"
)

type RpcArgs struct {
	User             string
	ApiKey           string
	Windows          []WindowLogs
	KeyLogs          []KeyLogs
	WindowLogsLength int
}

type WindowLogs struct {
	Name string
	Time time.Time
}

type KeyLogs struct {
	Time  time.Time
	Count int
}

const (
	BLOG_TYPE int = iota
	NOTE_TYPE
)

const LOCKED_SCREEN string = "__LOCKEDSCREEN"

const DB_SCHEMA string = `
create table if not exists users (
    id integer not null primary key autoincrement,
    username text,
    api_key text
);

create table if not exists windowLogs (
    id integer not null primary key autoincrement,
    uid integer,
    time integer,
    name text,
    foreign key (uid) references users(id)
);

create table if not exists keyLogs (
    id integer not null primary key autoincrement,
    uid integer,
    time integer,
    count integer,
    foreign key (uid) references users(id)
);

create table if not exists notes (
    id integer not null primary key autoincrement,
    uid integer,
    time integer,
    type integer,
    contents text,
    foreign key (uid) references users(id)
);
`

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

func (t *Golog) Log(args RpcArgs, result *int) error {
	uid, err := t.ensureAuth(args.User, args.ApiKey)
	if err != nil {
		log.Fatal(err)
		*result = 1
		return nil
	} else {
		log.Printf("%s authenticated successfully as uid %d\n", args.User, uid)
	}

	t.LogToDb(
		uid,
		args.Windows,
		args.KeyLogs,
		args.WindowLogsLength,
	)

	*result = 0
	return nil
}

func (t *Golog) SetupDb(db *sql.DB) {
	_, err := db.Exec(DB_SCHEMA)
	if err != nil {
		log.Fatal(err)
	}
	t.Db = db
}
