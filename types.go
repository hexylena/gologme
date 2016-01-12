package gologme

import (
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
