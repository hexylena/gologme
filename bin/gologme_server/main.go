package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/codegangsta/cli"
	"github.com/erasche/gologme/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "gologme"
	app.Usage = "local logging server"
	user, err := user.Current()
	var dbUrl string
	if err != nil {
		dbUrl = "gologme.db"
	} else {
		dbUrl = path.Join(user.HomeDir, ".gologme.db")
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "port",
			Value:  8080,
			Usage:  "port to listen on",
			EnvVar: "DB_PORT",
		},

		cli.StringFlag{
			Name:   "dbType",
			Usage:  "Database Type (sqlite3, postgres)",
			Value:  "sqlite3",
			EnvVar: "DB_TYPE",
		},
		cli.StringFlag{
			Name:   "dbUrl",
			Usage:  "Database URL",
			Value:  dbUrl,
			EnvVar: "DB_URL",
		},
	}

	app.Action = func(c *cli.Context) {
		server.ServeFromPath(
			c.String("dbType"),
			c.String("dbUrl"),
			fmt.Sprintf(":%d", c.Int("port")),
		)
	}
	app.Run(os.Args)
}
