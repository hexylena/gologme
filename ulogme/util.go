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

type StringLog struct {
	RealTime time.Time
	Title    string
}

type IntLog struct {
	RealTime time.Time
	Count    int
}

func WriteIntFile(category string, windows []IntLog, dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	var (
		f            *os.File
		lastUlogTime int64
		init         bool = false
	)

	for _, wl := range windows {
		ult := Ulogme7amTime(wl.RealTime)
		if !init || ult != lastUlogTime {
            _ = f.Close()
            fn := fmt.Sprintf("%s/%s_%d.txt", dir, category, ult)
            fmt.Printf("Writing to %s\n", fn)
			f, err = os.Create(fn)
			if err != nil {
				log.Fatal(err)
			}
			lastUlogTime = ult
			init = true
		}
		f.WriteString(fmt.Sprintf("%d %d\n", wl.RealTime.Unix(), wl.Count))
	}
}

func WriteStringFile(category string, windows []StringLog, dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	var (
		f            *os.File
		lastUlogTime int64
		init         bool = false
	)

	for _, wl := range windows {
		ult := Ulogme7amTime(wl.RealTime)
		if !init || ult != lastUlogTime {
            _ = f.Close()
            fn := fmt.Sprintf("%s/%s_%d.txt", dir, category, ult)
            fmt.Printf("Writing to %s\n", fn)
			f, err = os.Create(fn)
			if err != nil {
				log.Fatal(err)
			}
			lastUlogTime = ult
			init = true
		}
		f.WriteString(fmt.Sprintf("%d %s\n", wl.RealTime.Unix(), wl.Title))
	}
}

func Ulogme7amTime(t time.Time) int64 {
	// TODO: check timezones
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatal(err)
	}

    am7 := time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, location)
    if t.Hour() < 7 {
        am7 = am7.AddDate(0, 0, -1)
    }
    return am7.Unix()
}
