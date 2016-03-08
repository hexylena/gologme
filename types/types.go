package types

import (
	"time"
)

type DataLogRequest struct {
	User    string
	ApiKey  string
	Windows []*WindowLogs
	KeyLogs []*KeyLogs
}

type WindowLogs struct {
	Name string
	Time time.Time
}

type KeyLogs struct {
	Time  time.Time
	Count int
}

const (
	BLOG_TYPE int = iota
	NOTE_TYPE
)

const LOCKED_SCREEN string = "__LOCKEDSCREEN"

type SEvent struct {
	T int    `json:"t"`
	S string `json:"s"`
}

type IEvent struct {
	T int `json:"t"`
	S int `json:"s"`
}

type SEventT struct {
	RealTime time.Time
	Title    string
}

type IEventT struct {
	RealTime time.Time
	Count    int
}

type NoteEvent struct {
	RealTime time.Time
	Type     int
	Contents string
}

type EventLog struct {
	Blog           string    `json:"blog"`
	Note_events    []*SEvent `json:"notes_events"`
	Window_events  []*SEvent `json:"window_events"`
	Keyfreq_events []*IEvent `json:"keyfreq_events"`
}
