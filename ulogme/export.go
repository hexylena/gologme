package main

import (
	"log"
	"time"
)



func exportWindows(t *Golog, uid int) []StringLog {
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

	logs := make([]StringLog, 0)
	for rows.Next() {
		var (
			ltime int64
			name  string
		)

		rows.Scan(&ltime, &name)
		realTime := time.Unix(ltime, 0)
		logs = append(logs, StringLog{
			RealTime: realTime,
			Title:    name,
		})
	}
	return logs
}

func exportKeys(t *Golog, uid int) []IntLog {
	stmt, err := t.Db.Prepare("select time, count from keyLogs where uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]IntLog, 0)
	for rows.Next() {
		var (
			ltime int64
			count int
		)

		rows.Scan(&ltime, &count)
		realTime := time.Unix(ltime, 0)
		logs = append(logs, IntLog{
			RealTime: realTime,
			Count:    count,
		})
	}
	return logs
}
