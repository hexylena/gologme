package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

type Golog struct {
	Db *sql.DB
}

type ExportLog struct {
	RealTime time.Time
	Title    string
}

func WriteFile(category string, ulogtime int64, windows []ExportLog, dir string) {
	err := os.MkdirAll(dir, 0777)
	f, err := os.Create(fmt.Sprintf("%s/%s_%d.txt", dir, category, ulogtime))
	if err != nil {
		log.Fatal(err)
	}

	for _, wl := range windows {
		f.WriteString(fmt.Sprintf("%d %s\n", wl.RealTime.Unix(), wl.Title))
	}
}

func Ulogme7amTime(t time.Time) int64 {
	// TODO: check timezones
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatal(err)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, location).Unix()
}
