package main

import (
	"github.com/codegangsta/cli"
	"github.com/erasche/gologme/client"
	"os"
	"os/user"
	"path"
)

func main() {
	app := cli.NewApp()
	app.Name = "gologme"
	app.Usage = "local logging client"
	user, err := user.Current()
	var dbPath string
	if err != nil {
		dbPath = "gologme.db"
	} else {
		dbPath = path.Join(user.HomeDir, ".gologme.db")
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "buffSize",
			Value: 32,
			Usage: "size of buffer before sending logs",
		},
		cli.IntFlag{
			Name:  "windowLogGranularity",
			Value: 2000,
			Usage: "How often to poll window title in ms",
		},
		cli.IntFlag{
			Name:  "keyLogGranularity",
			Value: 2000,
			Usage: "How often to aggregate caught keypresses in ms",
		},
		cli.BoolFlag{
			Name:  "standalone",
			Usage: "Run in non-networked, standalone mode",
		},
		cli.StringFlag{
			Name:  "dbPath",
			Usage: "Path to the database",
			Value: dbPath,
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("standalone") {
			go client.Serve(
				dbPath,
				":8080",
			)
		}

		client.Golog(
			c.Int("buffSize"),
			c.Int("windowLogGranularity"),
			c.Int("keyLogGranularity"),
			c.Bool("standalone"),
		)

	}
	app.Run(os.Args)
}
