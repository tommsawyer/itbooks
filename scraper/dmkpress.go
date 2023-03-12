package scraper

import (
	"context"
	"errors"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func dmkpress(ctx context.Context, books chan<- Book) error {
	const startPage = "https://dmkpress.com/catalog/computer/?&filter%5Bavailable%5D=on&filter%5Bprice%5D%5Bfrom%5D=1&filter%5Bprice%5D%5Bto%5D=2999&filter%5Brelease_date%5D%5Bfrom%5D=1041379200&filter%5Brelease_date%5D%5Bto%5D=1767225600&filter%5Btranslator%5D=&filter%5Bformat%5D=&filter%5Bbumaga%5D=&filter%5Boblozhka%5D=&order_filter%5Brelease_date%5D=1"
	collector := colly.NewCollector()

	authorRegExp := regexp.MustCompile("(?is)Автор:\n(?P<Author>.*?)Дата выхода")

	// links on book list
	collector.OnHTML("#new-products .item-name a", func(h *colly.HTMLElement) {
		if err := h.Request.Visit(h.Attr("href")); err != nil && !errors.Is(err, colly.ErrAlreadyVisited) {
			log.Printf("[dmkpress.com] cannot visit book page: %v", err)
		}
	})

	// pagination
	collector.OnHTML(".pages.pull-right a", func(h *colly.HTMLElement) {
		page, err := strconv.Atoi(h.Text)
		if err != nil {
			return
		}
		// we don't need too old books
		if page > 5 {
			return
		}

		if !strings.Contains(h.Attr("href"), "filter") {
			return
		}

		if err := h.Request.Visit(h.Attr("href")); err != nil && !errors.Is(err, colly.ErrAlreadyVisited) {
			log.Printf("[dmkpress.com] cannot visit pagination page: %v", err)
		}
	})

	// book page
	collector.OnHTML("div[itemscope]", func(h *colly.HTMLElement) {
		if h.Attr("itemtype") != "http://schema.org/Product" {
			return
		}

		match := authorRegExp.FindStringSubmatch(h.Text)
		var authors []string
		if len(match) > 0 {
			authors = strings.Split(strings.TrimSpace(match[1]), ",")
			if strings.Contains(authors[len(authors)-1], "Перевод") {
				authors[len(authors)-1] = strings.Split(authors[len(authors)-1], "Перевод")[0]
			}
			for i := range authors {
				authors[i] = strings.TrimSpace(authors[i])
			}
		}

		books <- Book{
			ISBN:        path.Base(h.Request.URL.String()),
			URL:         h.Request.URL.String(),
			Title:       h.ChildText("span[itemprop=name]"),
			ImageURL:    h.Request.AbsoluteURL(h.ChildAttr(".card-img", "src")),
			Description: h.ChildText("#description"),
			Authors:     authors,
			Publisher:   "ДМК-Пресс",
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("[dmkpress.com] parsing " + r.URL.String())
	})

	return collector.Visit(startPage)
}
