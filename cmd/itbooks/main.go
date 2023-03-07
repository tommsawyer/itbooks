package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	itbooks := &cli.App{
		Name:     "itbooks",
		Usage:    "TODO",
		Commands: []*cli.Command{scrape, publish, test},
	}

	if err := itbooks.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
