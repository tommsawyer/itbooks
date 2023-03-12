package scraper

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func piter(ctx context.Context, books chan<- Book) error {
	const startPage = "https://www.piter.com/collection/kompyutery-i-internet?page_size=100&order=descending_age&q=&options%5B169105%5D%5B%5D=1717868"
	collector := colly.NewCollector()

	// book links on list page
	collector.OnHTML(".products-list a", func(h *colly.HTMLElement) {
		if err := h.Request.Visit(h.Attr("href")); err != nil && !errors.Is(err, colly.ErrAlreadyVisited) {
			log.Printf("[piter.com] cannot visit book link on products page: %v", err)
		}
	})

	// pagination elements
	collector.OnHTML(".pagination a", func(h *colly.HTMLElement) {
		if err := h.Request.Visit(h.Attr("href")); err != nil && !errors.Is(err, colly.ErrAlreadyVisited) {
			log.Printf("[piter.com] cannot visit pagination link: %v", err)
		}
	})

	// book page
	collector.OnHTML(".product-block", func(h *colly.HTMLElement) {
		authors := strings.Split(h.ChildText(".author"), ",")
		for i := range authors {
			authors[i] = strings.TrimSpace(authors[i])
		}
		img := h.DOM.Find(".coverProduct").AttrOr("src", "")

		books <- Book{
			ISBN:        h.ChildText("li:nth-child(7) .grid-7"),
			URL:         h.Request.URL.String(),
			Title:       h.ChildText(".product-info h1"),
			Authors:     authors,
			ImageURL:    img,
			Description: h.DOM.Parent().Find("#tab-1").Text(),
			Details: map[string]string{
				"year": h.ChildText("li:nth-child(2) .grid-7"),
			},
			Publisher: "Питер",
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("[piter.com] parsing " + r.URL.String())
	})

	return collector.Visit(startPage)
}
