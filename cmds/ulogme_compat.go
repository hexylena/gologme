package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

var dayStart int = 7

type Golog struct {
	Db *sql.DB
}

func ulogme7amTime(t time.Time) int64 {
	// TODO: check timezones
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatal(err)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, location).Unix()
}

type exportLog struct {
	realTime time.Time
	title    string
}

func writeFile(category string, ulogtime int64, windows []exportLog) {
	f, err := os.Create(fmt.Sprintf("logs/%s_%d.txt", category, ulogtime))
	if err != nil {
		log.Fatal(err)
	}

	for _, wl := range windows {
		f.WriteString(fmt.Sprintf("%d %s\n", wl.realTime.Unix(), wl.title))
	}
}

func (t *Golog) exportWindows(uid int) {
	stmt, err := t.Db.Prepare("select time, name from windowLogs where uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]exportLog, 0)
	var ulogtime int64 = 0
	for rows.Next() {
		var (
			ltime int64
			name  string
		)

		rows.Scan(&ltime, &name)
		realTime := time.Unix(ltime, 0)
		ulogtime = ulogme7amTime(realTime)
		logs = append(logs, exportLog{
			realTime: realTime,
			title:    name,
		})
	}

	if len(logs) > 0 {
		writeFile("window", ulogtime, logs)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "file.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	golog := new(Golog)
	golog.Db = db

	golog.exportWindows(1)
}
