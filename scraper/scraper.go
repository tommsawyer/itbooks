package scraper

import (
	"context"
	"log"
	"sync"

	"golang.org/x/exp/maps"
)

// Site is book publisher that can be scraped.
// It should send all parsed books to provided books channel.
type Site interface {
	Scrape(context.Context, chan<- Book) error
}

// SiteFunc is adapter to use funcs as site scrapers.
type SiteFunc func(context.Context, chan<- Book) error

func (f SiteFunc) Scrape(ctx context.Context, books chan<- Book) error {
	return f(ctx, books)
}

var sites = map[string]Site{
	"piter":    SiteFunc(scrapePiter),
	"dmkpress": SiteFunc(scrapeDMKPress),
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

// Scrape will scrape provided sites.
func Scrape(ctx context.Context, names ...string) <-chan Book {
	sitesToScrape := make([]Site, 0, len(names))
	for _, name := range names {
		sitesToScrape = append(sitesToScrape, sites[name])
	}

	return run(ctx, sitesToScrape...)
}

// ScrapeAll runs all scrapers.
func ScrapeAll(ctx context.Context) <-chan Book {
	return run(ctx, maps.Values(sites)...)
}

func run(ctx context.Context, sites ...Site) <-chan Book {
	books := make(chan Book)

	var wg sync.WaitGroup
	wg.Add(len(sites))
	for _, site := range sites {
		go func(site Site) {
			defer wg.Done()

			if err := site.Scrape(ctx, books); err != nil {
				log.Print(err)
			}
		}(site)
	}

	go func() {
		wg.Wait()
		close(books)
	}()

	return books
}
