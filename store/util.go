package store

import (
	"log"
	"time"
)

func Ulogme7amTime(t time.Time) int64 {
	// TODO: check timezones
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatal(err)
	}

	am7 := time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, location)
	if t.Hour() < 7 {
		am7 = Yesterday(am7)
	}
	return am7.Unix()
}

func Yesterday(t time.Time) time.Time {
	return t.AddDate(0, 0, -1)
}

func Tomorrow(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}
