package scraper

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeEksmo(ctx context.Context, books chan<- Book) error {
	const startPage = "https://eksmo.ru/professionalnaia-literatura/kompyuternaya-literatura/"
	collector := colly.NewCollector()

	// book links on list page
	collector.OnHTML(".book_fast-view[data-last-category=\"Компьютерная литература\"] .book__link", func(h *colly.HTMLElement) {
		if err := h.Request.Visit(h.Attr("href")); err != nil && !errors.Is(err, colly.ErrAlreadyVisited) {
			log.Printf("[eksmo.com] cannot visit book link on products page: %v", err)
		}
	})

	// pagination elements
	// obtaining only newest books
	collector.OnHTML(".pagenav__list a", func(h *colly.HTMLElement) {
		if err := h.Request.Visit(h.Attr("href")); err != nil && !errors.Is(err, colly.ErrAlreadyVisited) {
			log.Printf("[eksmo.com] cannot visit pagination link: %v", err)
		}
	})

	// book page
	collector.OnHTML(".book-page__card-cont", func(h *colly.HTMLElement) {
		authors := strings.Split(h.ChildText(".book-page__card-author a"), ",")
		for i := range authors {
			authors[i] = strings.TrimSpace(authors[i])
		}
		img := h.ChildAttr(".book-page__cover-link", "href")

		year := ""
		if cht := h.ChildText(".book-page__card-props div:nth-child(9) span"); cht == "Дата выхода:" {
			publishDate := strings.Split(h.ChildText(".book-page__card-props div:nth-child(9)"), " ")
			year = publishDate[len(publishDate)-1]
		}

		books <- Book{
			ISBN:        h.ChildText(".book-page__copy-isbn .copy__val"),
			URL:         h.Request.URL.String(),
			Title:       h.ChildText(".book-page__card-title"),
			Authors:     authors,
			ImageURL:    img,
			Description: h.ChildText(".spoiler__text p"),
			Details: map[string]string{
				"year": year,
			},
			Publisher: "Эксмо",
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("[eksmo.com] parsing " + r.URL.String())
	})

	return collector.Visit(startPage)
}
