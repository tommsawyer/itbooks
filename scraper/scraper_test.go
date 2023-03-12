package scraper

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestScraperParsesOneBook(t *testing.T) {
	is := is.New(t)
	publisher := &mockPublisher{
		book: Book{
			URL:         "url",
			ISBN:        "isbn",
			Title:       "title",
			Authors:     []string{"author"},
			ImageURL:    "image",
			Description: "description",
			Publisher:   "publisher",
			Details:     map[string]string{"test": "test"},
		},
	}

	parsedBook := <-run(context.Background(), publisher)

	is.Equal(parsedBook, publisher.book)
}

type mockPublisher struct {
	book Book
}

func (p *mockPublisher) Parse(ctx context.Context, books chan<- Book) error {
	books <- p.book
	return nil
}
