package main

import (
	"github.com/tommsawyer/itbooks/scraper"
	"github.com/urfave/cli/v2"
)

var scrape = &cli.Command{
	Name:  "scrape",
	Usage: "run scrapers",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "postgres-uri",
			Usage:   "uri to postgres",
			Value:   "postgres://itbooks:secret@localhost:5432/itbooks?sslmode=disable",
			EnvVars: []string{"POSTGRES_URI"},
		},
		&cli.StringSliceFlag{
			Name:    "publishers",
			Usage:   "publishers that needed to be scraped. All publishers if empty",
			Aliases: []string{"p"},
			EnvVars: []string{"PUBLISHERS"},
		},
	},
	Before: connectToPostgres,
	Action: func(c *cli.Context) error {
		ctx := c.Context

		publishers := c.StringSlice("publishers")
		if len(publishers) == 0 {
			scraper.RunAll(ctx)
			return nil
		}

		scraper.Run(ctx, publishers...)
		return nil
	},
}
