package scraper

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/tommsawyer/itbooks/postgres"
	"github.com/tommsawyer/itbooks/postgres/postgrestest"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	flag.Parse()

	uri, stopPostgres, err := postgrestest.RunContainer(ctx, testing.Verbose())
	if err != nil {
		log.Fatalf("cannot run testing postgres: %v", err)
	}

	if err := postgres.Connect(ctx, uri); err != nil {
		log.Fatalf("cannot connect to testing postgres: %v", err)
	}

	code := m.Run()
	stopPostgres()
	os.Exit(code)
}

func TestScraperSavesTheBook(t *testing.T) {
	ctx, is, rollback := postgres.RunAndRollback(t)
	defer rollback()

	publisher := &mockPublisher{
		book: postgres.UpsertBookParams{
			URL:         "url",
			ISBN:        "isbn",
			Title:       "title",
			Authors:     []string{"author"},
			Image:       "image",
			Description: "description",
			Publisher:   "publisher",
			Properties:  map[string]string{"test": "test"},
		},
	}

	run(ctx, publisher)

	books, err := postgres.FindBooks(ctx)
	is.NoErr(err)

	is.Equal(len(books), 1)

	book := books[0]
	is.Equal(book.ISBN.String, publisher.book.ISBN)
	is.Equal(book.URL.String, publisher.book.URL)
	is.Equal(book.Title.String, publisher.book.Title)
	is.Equal(book.Authors.Elements[0].String, publisher.book.Authors[0])
	is.Equal(book.Image.String, publisher.book.Image)
	is.Equal(book.Description.String, publisher.book.Description)
	is.Equal(book.Publisher.String, publisher.book.Publisher)
	is.Equal(book.Properties, publisher.book.Properties)
}

type mockPublisher struct {
	book postgres.UpsertBookParams
}

func (p *mockPublisher) Parse(ctx context.Context, books chan<- postgres.UpsertBookParams) error {
	books <- p.book
	return nil
}
