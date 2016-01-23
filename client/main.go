package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gologme"
	app.Usage = "local logging client"
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
	}

	app.Action = func(c *cli.Context) {
		golog(
			c.Int("buffSize"),
			c.Int("windowLogGranularity"),
			c.Int("keyLogGranularity"),
			c.Bool("standalone"),
		)
	}
	app.Run(os.Args)
}
