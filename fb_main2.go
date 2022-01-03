package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"log"
	"strings"
	"time"
)

type MarketplaceListing struct {
	title      string
	location   string
	price      string
	sold       bool
	imageUrl   string
	listingUrl string
}

func (l MarketplaceListing) String() string {
	return fmt.Sprintf("title: %v | location: %v | price: %v | sold: %v | imageUrl: %v, | listingUrl: %v",
		l.title, l.location, l.price, l.sold, l.imageUrl, l.listingUrl)
}

func main() {
	// create context
	//ctx, cancel := chromedp.NewContext(context.Background(), chromedp.Headless())
	//defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// to see what happen
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.WindowSize(1920, 1080),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// run task list
	//var res string
	var pageHtml string
	pageHtml = ""
	fmt.Println("Navigating and scraping...")
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.facebook.com/marketplace/109658472394519/search/?query=keyboard`),
		chromedp.WaitVisible(`div[role=main]`, chromedp.NodeVisible),
		//chromedp.Sleep(5*time.Second),
		//chromedp.ScrollIntoView(`#ssrb_root_end`, chromedp.BySearch),
		//chromedp.Sleep(1*time.Second),
		//chromedp.ScrollIntoView(`#ssrb_root_end`, chromedp.BySearch),
		//chromedp.Sleep(1*time.Second),
		//chromedp.Evaluate("document.querySelector('div[role=main]').scrollIntoViewIfNeeded(true)", nil),

		chromedp.KeyEvent(kb.End),
		chromedp.Sleep(2*time.Second),
		chromedp.KeyEvent(kb.End),
		chromedp.Sleep(2*time.Second),
		chromedp.KeyEvent(kb.End),
		chromedp.Sleep(2*time.Second),
		chromedp.KeyEvent(kb.End),
		chromedp.Sleep(2*time.Second),
		chromedp.KeyEvent(kb.End),
		chromedp.Sleep(2*time.Second),

		//chromedp.SendKeys(`div[role=main]`, kb.End, chromedp.NodeVisible),
		//chromedp.Value(`#input1`, val1, chromedp.ByID),
		//chromedp.KeyEvent(kb.End)
		chromedp.InnerHTML(`div[role=main]`, &pageHtml, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsing HTML...")

	content, err := goquery.NewDocumentFromReader(strings.NewReader(pageHtml))
	if err != nil {
		panic(err)
	}

	//fmt.Println(pageHtml)

	posts := content.Find("a[role=link]")
	log.Printf("[facebook] found %v listings", posts.Length())
	posts.Each(func(i int, s *goquery.Selection) {
		var listing MarketplaceListing
		listing.sold = false
		href, hrefExists := s.Attr("href")
		if hrefExists {
			listing.listingUrl = "https://facebook.com" + href
		}
		image, foundImage := s.Find("img").First().Attr("src")
		if foundImage {
			listing.imageUrl = image
		}

		// To get the info under the image of each listing, we use the CSS selector "span[dir=auto]"
		// This will find each piece of text and works well with simple listings, but ones that have reduced prices, are
		// sold, or both will return each of those chunks of text individually.
		//
		// To parse this all, we start at the bottom of the list of results since the location and listing title are
		// always the last and second to last items. After that, we individually check the remaining chunks and assign
		// them as appropriate.

		listingInfo := s.Find("span[dir=auto]")
		listingInfoItems := make([]*goquery.Selection, listingInfo.Length())
		listingInfo.Each(func(j int, t *goquery.Selection) {
			listingInfoItems[j] = t
		})
		listing.location = listingInfoItems[len(listingInfoItems)-1].Text()
		listing.title = strings.TrimSpace(listingInfoItems[len(listingInfoItems)-2].Text())

		for a := 0; a < len(listingInfoItems)-2; a++ {
			item := strings.TrimSpace(listingInfoItems[a].Text())
			if item == "·" { // Character put in between "Sold" and the price
				continue
			}
			if strings.ToLower(item) == "sold" {
				listing.sold = true
				continue
			}
			if listing.price == "" {
				listing.price = item
			}
		}

		//fmt.Println(listing)
	})

}
