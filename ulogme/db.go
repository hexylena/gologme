package main

import (
	"log"
	"time"
)

func exportWindows(t *Golog, uid int) (int64, []ExportLog) {
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
	var ulogtime int64 = 0
	for rows.Next() {
		var (
			ltime int64
			name  string
		)

		rows.Scan(&ltime, &name)
		realTime := time.Unix(ltime, 0)
		ulogtime = Ulogme7amTime(realTime)
		logs = append(logs, ExportLog{
			RealTime: realTime,
			Title:    name,
		})
	}
	return ulogtime, logs
}
