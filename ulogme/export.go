package main

import (
	"log"
	"time"

	gologme "github.com/erasche/gologme/util"
)

func exportNotes(t *gologme.Golog, uid int) []gologme.NoteEvent {
	stmt, err := t.Db.Prepare("select time, type, contents from notes where uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	logs := make([]gologme.NoteEvent, 0)
	for rows.Next() {
		var (
			ltime int64
			ntype int
			name  string
		)

		rows.Scan(&ltime, &ntype, &name)
		realTime := time.Unix(ltime, 0)
		logs = append(logs, gologme.NoteEvent{
			RealTime: realTime,
			Type:     ntype,
			Contents: name,
		})
	}
	return logs
}

func exportWindows(t *gologme.Golog, uid int) []gologme.SEventT {
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

	logs := make([]gologme.SEventT, 0)
	for rows.Next() {
		var (
			ltime int64
			name  string
		)

		rows.Scan(&ltime, &name)
		realTime := time.Unix(ltime, 0)
		logs = append(logs, gologme.SEventT{
			RealTime: realTime,
			Title:    name,
		})
	}
	return logs
}

func exportKeys(t *gologme.Golog, uid int) []gologme.IEventT {
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

	logs := make([]gologme.IEventT, 0)
	for rows.Next() {
		var (
			ltime int64
			count int
		)

		rows.Scan(&ltime, &count)
		realTime := time.Unix(ltime, 0)
		logs = append(logs, gologme.IEventT{
			RealTime: realTime,
			Count:    count,
		})
	}
	return logs
}
