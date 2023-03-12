package scraper

import (
	"context"
	"fmt"
	"log"
	"sync"

	"golang.org/x/exp/maps"
)

// Publisher is publisher that can be scraped.
// It should send all parsed books to provided books channel.
type Publisher interface {
	Parse(context.Context, chan<- Book) error
}

// PublisherFunc is adapter to use funcs as Publishers.
type PublisherFunc func(context.Context, chan<- Book) error

func (p PublisherFunc) Parse(ctx context.Context, books chan<- Book) error {
	return p(ctx, books)
}

var publishers = map[string]Publisher{
	"piter":    PublisherFunc(piter),
	"dmkpress": PublisherFunc(dmkpress),
}

// Book represents parsed book.
type Book struct {
	URL         string
	ImageURL    string
	ISBN        string
	Title       string
	Authors     []string
	Description string
	Publisher   string
	Details     map[string]string
}

// Run will scrape provided publishers.
func Run(ctx context.Context, names ...string) <-chan Book {
	publishersToScrape := make([]Publisher, 0, len(names))
	for _, name := range names {
		publishersToScrape = append(publishersToScrape, publishers[name])
	}

	return run(ctx, publishersToScrape...)
}

// RunAll runs all scrapers.
func RunAll(ctx context.Context) <-chan Book {
	return run(ctx, maps.Values(publishers)...)
}

// Test will just log scraper events.
func Test(ctx context.Context, name string) {
	publisher := publishers[name]
	books := make(chan Book)

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

func run(ctx context.Context, publishers ...Publisher) <-chan Book {
	books := make(chan Book)

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

	return books
}
