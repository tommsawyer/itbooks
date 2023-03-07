package scraper

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/tommsawyer/itbooks/postgres"
	"golang.org/x/exp/maps"
)

// Publisher is publisher that can be scraped.
// It should send all parsed books to provided books channel.
type Publisher interface {
	Parse(context.Context, chan<- postgres.UpsertBookParams) error
}

// PublisherFunc is adapter to use funcs as Publishers.
type PublisherFunc func(context.Context, chan<- postgres.UpsertBookParams) error

func (p PublisherFunc) Parse(ctx context.Context, books chan<- postgres.UpsertBookParams) error {
	return p(ctx, books)
}

var publishers = map[string]Publisher{
	"piter":    PublisherFunc(piter),
	"dmkpress": PublisherFunc(dmkpress),
}

// Run will scrape provided publishers.
func Run(ctx context.Context, names ...string) {
	publishersToScrape := make([]Publisher, 0, len(names))
	for _, name := range names {
		publishersToScrape = append(publishersToScrape, publishers[name])
	}

	run(ctx, publishersToScrape...)
}

// RunAll runs all scrapers.
func RunAll(ctx context.Context) {
	run(ctx, maps.Values(publishers)...)
}

// Test will just log scraper events.
func Test(ctx context.Context, name string) {
	publisher := publishers[name]
	books := make(chan postgres.UpsertBookParams)

	go func() {
		defer close(books)
		if err := publisher.Parse(ctx, books); err != nil {
			log.Println(err)
		}
	}()

	for book := range books {
		log.Println("Scraped: ", fmt.Sprintf("%#v", book))
	}
}

func run(ctx context.Context, publishers ...Publisher) {
	books := make(chan postgres.UpsertBookParams)

	var wg sync.WaitGroup
	wg.Add(len(publishers))
	for _, publisher := range publishers {
		go func(publisher Publisher) {
			defer wg.Done()

			if err := publisher.Parse(ctx, books); err != nil {
				log.Print(err)
			}
		}(publisher)
	}

	go func() {
		wg.Wait()
		close(books)
	}()

	for book := range books {
		if _, err := postgres.UpsertBook(ctx, book); err != nil {
			log.Print(err)
		}
	}
}
