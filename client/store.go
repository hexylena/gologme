package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	gologme "github.com/erasche/gologme/types"
	_ "github.com/mattn/go-sqlite3"
)

type receiver struct {
	ServerAddress string
}

func (r *receiver) send(wl []*gologme.WindowLogs, kl []*gologme.KeyLogs) {
	args := &gologme.DataLogRequest{
		User:    "admin",
		ApiKey:  "deadbeefcafe",
		Windows: wl,
		KeyLogs: kl,
	}

	// Marshal into str
	data, err := json.Marshal(args)
	if err != nil {
		fmt.Println(err)
	}

	//// Post to our server URL
	req, err := http.NewRequest(
		"POST",
		r.ServerAddress,
		strings.NewReader(string(data)),
	)
	hc := http.Client{}
	_, err = hc.Do(req)

	if err != nil {
		fmt.Println(err)
	}
}
