package main

import (
	"log"

	"github.com/tommsawyer/itbooks/postgres"
	"github.com/tommsawyer/itbooks/telegram"
	"github.com/urfave/cli/v2"
)

var publish = &cli.Command{
	Name:  "publish",
	Usage: "publishes book to provided telegram channel",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "postgres-uri",
			Usage:   "uri to postgres",
			Value:   "postgres://itbooks:secret@localhost:5432/itbooks?sslmode=disable",
			Aliases: []string{},
			EnvVars: []string{"POSTGRES_URI"},
		},
		&cli.StringFlag{
			Name:    "telegram-token",
			Usage:   "token for telegram bot",
			Value:   "",
			Aliases: []string{"t"},
			EnvVars: []string{"TELEGRAM_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "telegram-group-id",
			Usage:   "id of telegram group",
			Value:   "",
			Aliases: []string{"g"},
			EnvVars: []string{"TELEGRAM_GROUP"},
		},
		&cli.StringFlag{
			Name:    "isbn",
			Usage:   "isbn of book to publish",
			Value:   "",
			EnvVars: []string{"ISBN"},
		},
	},
	Before: combine(connectToPostgres, authorizeInTelegram),
	Action: func(c *cli.Context) error {
		ctx := c.Context

		var b *postgres.Book
		isbn := c.String("isbn")
		if isbn == "" {
			unpublished, err := postgres.FindUnpublishedBooks(ctx)
			if err != nil {
				return err
			}

			if len(unpublished) == 0 {
				log.Println("no unpublished books, skipping...")
				return nil
			}

			b = unpublished[0]
		} else {
			isbnBook, err := postgres.GetBookByISBN(ctx, isbn)
			if err != nil {
				return err
			}

			b = isbnBook
		}

		if err := telegram.Send(ctx, c.String("telegram-group-id"), telegram.Message{
			ImageURL: b.Image.String,
			Title:    b.Title.String,
			Subtitle: b.Publisher.String,
			Link:     b.Image.String,
			Text:     b.Description.String,
		}); err != nil {
			return err
		}

		return postgres.SetBookPublished(ctx, b.ID, true)
	},
}
