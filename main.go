package main

import (
	"log"
	"strings"

	"github.com/go-rod/rod"
)

func main() {
	showMore := true
	url := "https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Elesas-Sparkle-V2-s12a246"

	browser := rod.New().
		MustConnect()

	defer browser.MustClose()

	// TODO: Press "loadMoreButton" button as long as it's present
	page := browser.MustPage(url).MustWaitStable()

	title := page.MustElement("h1").MustText()
	subject := strings.SplitAfter(title, ")")
	cardName := subject[0]
	setName := strings.TrimSpace(subject[1])
	log.Printf("Looking Card -> %s", cardName)
	log.Printf("Looking Set -> %s", setName)

	for showMore {
		showMoreButton := page.MustElement("#loadMoreButton")
		attrs := showMoreButton.MustDescribe().Attributes
		log.Printf("Attributes %v", attrs)
		disabled, err := showMoreButton.Attribute("disabled")
		if err != nil {
			panic(err)
		}
		if disabled == nil {
			log.Println("Showing more sellers...")
			showMoreButton.MustWaitStable().MustClick().MustWaitInvisible()
		} else {
			log.Printf("Show More Results - disabled attr(%s)", *disabled)
			log.Println("No more sellers to show!")
			showMore = false
		}
	}

	log.Println("Looking sellers...")
	sellers := page.MustElements("div.row.g-0.article-row")
	for _, seller := range sellers {
		// TODO: Use MustElementByJs using jquery for easier querys
		sellerInfo := seller.MustElement(".col-seller")
		sellerName := sellerInfo.MustElement("a[href]").MustText()
		sellerProductInfo := seller.MustElement(".col-product")
		sellerProductCondition := sellerProductInfo.MustElement("a.article-condition").MustText()
		sellerOffer := seller.MustElement(".col-offer")
		sellerProductPrice := sellerOffer.MustElement(".price-container").MustText()
		sellerProductQuantity := sellerOffer.MustElement(".item-count").MustText()

		log.Printf("Seller %s -> %s x %s(%s) - %s", sellerName, sellerProductQuantity, cardName, sellerProductCondition, sellerProductPrice)
	}

	// resp, err := http.Get(url)
	// if err != nil {
	// 	panic(err)
	// }
	// cookies := resp.Cookies()
	// log.Printf("Cookies %+v", cookies)

	// c := colly.NewCollector()
	// err = c.SetCookies(url, cookies)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// c.OnResponse(func(f *colly.Response) {
	// 	log.Printf("Response -> %s", string(f.Body))
	// })
	// c.OnHTML("h1", func(h *colly.HTMLElement) {
	// 	log.Printf("Looking Card -> %s", h.Text)
	// })
	// c.OnHTML("div.row.g-0.article-row", func(h *colly.HTMLElement) {
	// 	log.Printf("Looking Seller -> %+v", h)
	// })
	// c.OnRequest(func(r *colly.Request) {
	// 	log.Printf("Scrapping URL -> %s", r.URL)
	// })
	// c.OnError(func(r *colly.Response, err error) {
	// 	log.Println("Something went wrong:", err)
	// })

	// c.Visit(url)

}
