package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	gologme "github.com/erasche/gologme/types"
	"github.com/gorilla/mux"
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

func DataUpload(w http.ResponseWriter, r *http.Request) {
	//var logs
	fmt.Println("endpoints.DataUpload")
	decoder := json.NewDecoder(r.Body)
	logData := new(gologme.DataLogRequest)
	err := decoder.Decode(&logData)

	if err != nil {
		//log.Error(fmt.Sprintf("Error unmarshalling posted data %s", err))
		http.Error(w, "Invalid Route Data", http.StatusBadRequest)
		return
	}
	fmt.Println("endpoints.DataUpload 2")

	golog.Log(logData)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
