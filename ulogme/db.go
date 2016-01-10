package main

import (
	//"strings"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"time"
)

func importWindows(t *Golog, uid int, logDir string) {
	tx, err := t.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into windowLogs (uid, time, name) values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	files, err := filepath.Glob(logDir + "/window*.txt")
	for _, file := range files {
		handle, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		scanner := bufio.NewScanner(handle)
		for scanner.Scan() {
			t := scanner.Text()
			if len(t) > 12 {
				unixtime := t[0:10]
				title := t[11:]
				_, err := stmt.Exec(uid, unixtime, title)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		break
	}
	tx.Commit()
}

func exportWindows(t *Golog, uid int) []ExportLog {
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

	logs := make([]ExportLog, 0)
	for rows.Next() {
		var (
			ltime int64
			name  string
		)

		rows.Scan(&ltime, &name)
		realTime := time.Unix(ltime, 0)
		logs = append(logs, ExportLog{
			RealTime: realTime,
			Title:    name,
		})
	}
	return logs
}
