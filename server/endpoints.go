package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/erasche/gologme/types"
	util "github.com/erasche/gologme/util"
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

type eventListJson struct {
	FName string `json:"fname"`
	T0    int    `json:"t0"`
	T1    int    `json:"t1"`
}

func ExportList(w http.ResponseWriter, r *http.Request) {
	min_time, max_time := golog.RecordedDataRange()
	// add a day to max_time
	max_time.Add(time.Hour * 48)
	// Convert to u7am
	imin_time := util.Ulogme7amTime(min_time)
	imax_time := util.Ulogme7amTime(max_time)
	//{"fname": "/api/events/1448197200", "t0": 1448197200, "t1": 1448283600}

	eventList := make([]*eventListJson, 0)
	var aDayInSeconds = int64((time.Hour * 24).Seconds())
	var i int64
	for i = imin_time; i < imax_time; i += aDayInSeconds {
		eventList = append(eventList, &eventListJson{
			FName: fmt.Sprintf("/api/events/%d", int(i)),
			T0:    int(i),
			T1:    int(i + aDayInSeconds),
		})
	}
	js, _ := json.MarshalIndent(eventList, "", "  ")
	w.Write(js)
}

func DataUpload(w http.ResponseWriter, r *http.Request) {
	//var logs
	fmt.Println("endpoints.DataUpload")
	decoder := json.NewDecoder(r.Body)
	logData := new(types.DataLogRequest)
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
