package gologme

import (
	"time"
)

type Result int

type RpcArgs struct {
	User    int
	Windows []WindowLogs
	Length  int
}

type WindowLogs struct {
	Name string
	Time time.Time
}

type KeyLogs struct {
	Time time.Time
    Name string
}
