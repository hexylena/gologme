package main

import (
    "fmt"
	"github.com/codegangsta/cli"
	"github.com/erasche/gologme/server"
	"os"
	"os/user"
	"path"
)

func main() {
	app := cli.NewApp()
	app.Name = "gologme"
	app.Usage = "local logging server"
	user, err := user.Current()
	var dbPath string
	if err != nil {
		dbPath = path.Join(user.HomeDir, ".gologme.db")
	} else {
		dbPath = "~/.gologme.db"
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "port to listen on",
		},
		cli.StringFlag{
			Name:  "dbPath",
			Usage: "Path to the database",
			Value: dbPath,
		},
	}

	app.Action = func(c *cli.Context) {
        server.ServeFromPath(
            c.String("dbPath"),
            fmt.Sprintf(":%d", c.Int("port")),
        )
	}
	app.Run(os.Args)
}
