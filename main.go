package main

import (
	"flag"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/sa7mon/russ/data"
	"github.com/sa7mon/russ/scrapers"
	"log"
	"net/http"
	"time"
)

func main() {
	var scrapeInterval int
	flag.IntVar(&scrapeInterval, "interval", 60, "Minutes to wait between scrapes")
	flag.Parse()
	log.Printf("[main] app started. Scraping every %v minutes", scrapeInterval)

	data.GetManager().CurrentFeed = &feeds.Feed{
		Title:       "Free Geek",
		Link:        &feeds.Link{Href: ""},
		Description: "Free Geek Twin Cities - Cables/Peripherals",
		Author:      &feeds.Author{Name: "dan", Email: "dan@salmon.cat"},
		Created:     time.Now(),
		Items:       []*feeds.Item{},
	}

	go func(scrapeInterval int) {
		for {
			freeGeekItems, err := scrapers.ScrapeFreeGeek()
			if err != nil {
				log.Printf("[freegeek] got error when scraping: %v", err)
			} else {
				log.Printf("[freegeek] scrape successful")
				data.GetManager().CurrentFeed.Items = freeGeekItems
			}
			time.Sleep(time.Duration(scrapeInterval) * time.Minute)
		}
	}(scrapeInterval)

	r := mux.NewRouter()
	r.HandleFunc("/rss", RSSHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func RSSHandler(w http.ResponseWriter, r *http.Request) {
	manager := data.GetManager()
	rss, err := manager.CurrentFeed.ToRss()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/rss+xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rss))
}
