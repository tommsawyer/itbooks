package main

import (
	"github.com/tommsawyer/itbooks/scraper"
	"github.com/urfave/cli/v2"
)

var test = &cli.Command{
	Name:  "test",
	Usage: "test scraper",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "publisher",
			Usage:   "publisher that needed to be scraped",
			Aliases: []string{"p"},
			EnvVars: []string{"PUBLISHER"},
		},
	},
	Action: func(c *cli.Context) error {
		scraper.Test(c.Context, c.String("publisher"))
		return nil
	},
}
