package main

import (
	"fmt"

	"github.com/tommsawyer/itbooks/postgres"
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

		var books <-chan scraper.Book
		publishers := c.StringSlice("publishers")
		if len(publishers) == 0 {
			books = scraper.RunAll(ctx)
		} else {
			books = scraper.Run(ctx, publishers...)
		}

		for book := range books {
			if _, err := postgres.UpsertBook(ctx, postgres.UpsertBookParams{
				ISBN:        book.ISBN,
				URL:         book.URL,
				Title:       book.Title,
				Image:       book.ImageURL,
				Description: book.Description,
				Authors:     book.Authors,
				Publisher:   book.Publisher,
				Properties:  book.Details,
			}); err != nil {
				return fmt.Errorf("cannot save book: %w", err)
			}
		}

		return nil
	},
}
