package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func Events(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := strconv.ParseInt(vars["date"], 10, 64)
	if err != nil {
		// handle this
	}

	tm := time.Unix(i, 0)
	eventData := golog.ExportEventsByDate(tm)
	js, err := json.MarshalIndent(eventData, "", "  ")
	if err != nil {
		// handle
	}
	w.Write(js)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
