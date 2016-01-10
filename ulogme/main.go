package main

import (
	"database/sql"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", "file.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	golog := new(Golog)
	golog.Db = db

	app := cli.NewApp()
	app.Name = "ulogmeExport"
	app.Usage = "export logs to ulogme"
	app.Commands = []cli.Command{
		{
			Name:    "import",
			Aliases: []string{"i"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "logDir",
					Usage: "Directory containing ulogme's logs",
				},
				cli.IntFlag{
					Name:  "uid",
					Value: 1,
					Usage: "User ID from database",
				},
			},
			Usage: "import logs",
			Action: func(c *cli.Context) {
			},
		},
		{
			Name:    "export",
			Aliases: []string{"e"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "logDir",
					Usage: "Directory to output ulogme's logs",
				},
				cli.IntFlag{
					Name:  "uid",
					Value: 1,
					Usage: "User ID from database",
				},
			},
			Usage: "import logs",
			Action: func(c *cli.Context) {
				ulogtime, logs := exportWindows(
					golog,
					c.Int("uid"),
				)
				if len(logs) > 0 {
					WriteFile(
						"window",
						ulogtime,
						logs,
						c.String("logDir"),
					)
				}
			},
		},
	}
	app.Run(os.Args)
}
