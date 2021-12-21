package scrapers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"log"
	"net/http"
	"strings"
	"time"
)

func ScrapeFreeGeek() ([]*feeds.Item, error) {
	var feedItems []*feeds.Item

	// Request the HTML page.
	res, err := http.Get("https://www.freegeektwincities.org/cables")
	if err != nil {
		log.Fatalf("[freegeek] error getting page: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("[freegeek] bad status code: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	products := doc.Find(".ProductList-grid .ProductList-item")
	log.Printf("[freegeek] Found %v products", products.Length())
	products.Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find(".ProductList-title").Text()
		price := fmt.Sprintf("$%v", strings.TrimSpace(s.Find(".product-price").Text()))
		image, foundImage := s.Find(".ProductList-image--primary").First().Attr("data-src")
		if foundImage {
			image = strings.TrimSpace(image)
		}
		link, foundLink := s.Find(".ProductList-item-link").First().Attr("href")
		if foundLink {
			link = "https://www.freegeektwincities.org" + strings.TrimSpace(link)
		}
		feedItem := &feeds.Item{
			Title:       strings.TrimSpace(title),
			Link:        &feeds.Link{Href: link},
			Description: fmt.Sprintf("<div><ul><li>Title: %v</li><li>Price: %v</li></ul><img src='%v' /></div>", title, price, image),
			Author:      &feeds.Author{Name: "Free Geek", Email: ""},
			Created:     time.Now(),
		}
		feedItems = append(feedItems, feedItem)
	})

	return feedItems, nil
}
