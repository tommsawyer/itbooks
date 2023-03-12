package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tommsawyer/itbooks/postgres"
	"github.com/tommsawyer/itbooks/telegram"
	"github.com/urfave/cli/v2"
)

func main() {
	itbooks := &cli.App{
		Name:     "itbooks",
		Usage:    "TODO",
		Commands: []*cli.Command{scrape, publish},
	}

	if err := itbooks.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func combine(hooks ...cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		for _, hook := range hooks {
			if err := hook(ctx); err != nil {
				return err
			}
		}

		return nil
	}
}

func connectToPostgres(ctx *cli.Context) error {
	err := postgres.Connect(ctx.Context, ctx.String("postgres-uri"))
	if err != nil {
		return fmt.Errorf("cannot connect to postgres: %w", err)
	}

	return nil
}

func authorizeInTelegram(ctx *cli.Context) error {
	if err := telegram.Authorize(ctx.Context, ctx.String("telegram-token")); err != nil {
		return fmt.Errorf("cannot authorize in telegram: %w", err)
	}

	return nil
}
