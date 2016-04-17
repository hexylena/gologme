package main

import (
	"os"
	"os/user"
	"path"
	"time"

	"github.com/codegangsta/cli"
	"github.com/erasche/gologme/client"
	gologme "github.com/erasche/gologme/types"
)

func main() {
	app := cli.NewApp()
	app.Name = "gologme"
	app.Usage = "local logging client"
	user, err := user.Current()
	var defaultDbPath string
	if err != nil {
		defaultDbPath = "gologme.db"
	} else {
		defaultDbPath = path.Join(user.HomeDir, ".gologme.db")
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "standalone",
			Usage: "Run in non-networked, standalone mode",
		},
		cli.StringFlag{
			Name:  "dbPath",
			Usage: "Path to the database",
			Value: defaultDbPath,
		},
		cli.StringFlag{
			Name:  "serverAddr",
			Usage: "Address to send logs to, defaults to localhost for --standalone mode.",
			Value: "http://127.0.0.1:10000",
		},
	}

	app.Action = func(c *cli.Context) {
		receiver := &client.Receiver{
			ServerAddress: c.String("serverAddr"),
		}
		wlogs := []*gologme.WindowLogs{
			&gologme.WindowLogs{Name: "testing-window", Time: time.Now()},
			&gologme.WindowLogs{Name: "testing-window2", Time: time.Now()},
		}
		klogs := []*gologme.KeyLogs{}
		receiver.Send(wlogs, klogs)
	}
	app.Run(os.Args)
}
