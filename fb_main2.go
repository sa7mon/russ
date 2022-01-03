package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
)

type MarketplaceListing struct {
	title      string
	location   string
	price      string
	imageUrl   string
	listingUrl string
}

func (l MarketplaceListing) String() string {
	return fmt.Sprintf("title: %v | location: %v | price: %v | imageUrl: %v, | listingUrl: %v",
		l.title, l.location, l.price, l.imageUrl, l.listingUrl)
}

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	//var res string
	var pageHtml string
	pageHtml = ""
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.facebook.com/marketplace/109658472394519/search/?query=keyboard`),
		chromedp.InnerHTML(`div[role=main]`, &pageHtml, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	content, err := goquery.NewDocumentFromReader(strings.NewReader(pageHtml))
	if err != nil {
		panic(err)
	}

	//fmt.Println(pageHtml)

	posts := content.Find("a[role=link]")
	log.Printf("[facebook] found %v listings", posts.Length())
	posts.Each(func(i int, s *goquery.Selection) {
		var listing MarketplaceListing
		href, hrefExists := s.Attr("href")
		if hrefExists {
			listing.listingUrl = href
		}
		image, foundImage := s.Find("img").First().Attr("src")
		if foundImage {
			listing.imageUrl = image
		}
		listingInfo := s.Find("span[dir=auto]")
		listingInfo.Each(func(j int, t *goquery.Selection) {
			switch j {
			case 0:
				listing.price = t.Text()
				break
			case 1:
				listing.title = t.Text()
				break
			case 2:
				listing.location = t.Text()
				break
			default:
				break
			}
		})

		fmt.Println(listing)
	})

}
