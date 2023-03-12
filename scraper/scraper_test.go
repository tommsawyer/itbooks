package scraper

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestScraperParsesOneBook(t *testing.T) {
	is := is.New(t)
	site := &mockSite{
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

	parsedBook := <-run(context.Background(), site)

	is.Equal(parsedBook, site.book)
}

type mockSite struct {
	book Book
}

func (p *mockSite) Scrape(ctx context.Context, books chan<- Book) error {
	books <- p.book
	return nil
}
