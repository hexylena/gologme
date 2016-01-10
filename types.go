package gologme

import (
	"time"
)

type Result int

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
	Time time.Time
	Name string
}
