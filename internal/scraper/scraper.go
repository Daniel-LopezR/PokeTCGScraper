package scraper

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

const CARDMARKET_URL = "https://www.cardmarket.com"

type Card struct {
	Name        string
	Url         string
	Condition   string
	Language    string
	Description string
	Quantity    int
	Price       float32
}

type Seller struct {
	//mu sync.Mutex
	Name           string
	Region         string
	Category       string
	Url            string
	CardsAvaialble []Card
}

// Create a page if needed
func GetCreatePage(browser *rod.Browser) func() *rod.Page {
	create := func() *rod.Page {
		// We use MustIncognito to isolate pages with each other
		return browser.MustIncognito().MustPage()
	}
	return create
}

func (s *Seller) Round(b *rod.Browser, wantedList *[]string) {
	//b := rod.New().MustConnect()
	queryURL := "/Offers/Singles?idExpansion=5198&idLanguage=7&condition=2"
	// a[data-direction="next"]
	for _, chase := range *wantedList {
		encodedChase := url.QueryEscape(chase)
		fullURL := s.Url + queryURL + encodedChase

		log.Printf("Loading page for chase card (%s)", chase)
		sellerCardsPage := b.MustIncognito().MustPage(fullURL).MustWaitDOMStable()
		defer sellerCardsPage.MustClose()

		cardsTable := sellerCardsPage.MustElement("div.table-body")

		if cardsTable.MustHas("div.article-row") {
			availableCards := cardsTable.MustElements("div.article-row")
			for _, avCard := range availableCards {
				s.LoadSellerCardRow(false, &Card{
					Name:        chase,
					Description: "",
				}, avCard)
			}
		}
	}
}

func (s *Seller) LoadSellerCardRow(loadSeller bool, sCard *Card, row *rod.Element) {
	if loadSeller {
		// Seller Info
		sellerInfo := row.MustElement(".col-seller")

		sellerLinkElement := sellerInfo.MustElement("a[href]")
		s.Name = sellerLinkElement.MustText()
		s.Url = CARDMARKET_URL + *sellerLinkElement.MustAttribute("href")

		sellerRegionAttr := *sellerInfo.MustElement("span.icon").MustAttribute("aria-label")
		s.Region = strings.TrimSpace(strings.Split(sellerRegionAttr, ":")[1])

		if sellerInfo.MustHas("span.fonticon-users-powerseller") {
			s.Category = "Powerseller"
		} else if sellerInfo.MustHas("span.fonticon-users-professional") {
			s.Category = "Professional"
		} else {
			s.Category = "Hobby"
		}

		log.Printf("Added new Seller%+v\n", s)
	} else {
		cardInfo := row.MustElement(".col-seller")
		cardLinkElement := cardInfo.MustElement("a[href]")
		cardName := cardLinkElement.MustText()
		if cardName == sCard.Name {
			// Load Card Url
			sCard.Url = CARDMARKET_URL + *row.MustElement(".col-seller").MustElement("a").MustAttribute("href")
			log.Println("Updated Card Url!")
		} else {
			return
		}
	}

	// Seller product
	sellerProductInfo := row.MustElement(".col-product")

	sCard.Condition = sellerProductInfo.MustElement("a.article-condition").MustText()
	languageP := sellerProductInfo.MustWaitStable().MustElement("span.icon").MustAttribute("aria-label")
	if languageP != nil {
		sCard.Language = *languageP
	}
	if ok, descriptionElement, _ := sellerProductInfo.Has("span.text-muted"); ok {
		sCard.Description = descriptionElement.MustText()
	}

	sellerOffer := row.MustElement(".col-offer")

	sellerProductQuantity := sellerOffer.MustElement(".item-count").MustText()
	quantity, err := strconv.Atoi(sellerProductQuantity)
	if err != nil {
		sCard.Quantity = 0
	} else {
		sCard.Quantity = quantity
	}

	sellerProductPrice := sellerOffer.MustElement(".price-container").MustText()
	unparsedPrice := strings.Split(strings.Split(sellerProductPrice, " ")[0], ",")
	unparsedPrice[0] = strings.ReplaceAll(
		unparsedPrice[0],
		".", "",
	)
	price, err := strconv.ParseFloat(
		fmt.Sprintf("%s.%s", unparsedPrice[0], unparsedPrice[1]),
		32,
	)

	if err != nil {
		sCard.Price = 0
	} else {
		sCard.Price = float32(price)
	}
	log.Printf("Added new Card(%d x %s(%s) - %0.2f €) to Seller %s\n", sCard.Quantity, sCard.Name, sCard.Condition, sCard.Price, s.Name)

	s.CardsAvaialble = append(s.CardsAvaialble, *sCard)
}

func scrap(wantedList []string) {
	sellers := []*Seller{}
	showMore := true
	// TODO: Loop wanted list instead of hardcoding url
	url := "https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Cynthias-Ambition-V2-s12a239?language=7&minCondition=2"
	// TODO: import wantedList from e.g. csv, maybe cardmarket api(don't think i will do this tbh)

	log.Println("Creating browser...")
	//u := launcher.New().
	//	Set("ozone-platform", "wayland").
	//	Set("headless").MustLaunch()
	u := "ws://127.0.0.1:9222/devtools/browser/da8466a0-e51f-4131-9629-d794b6e7db07"
	//u := launcher.New().Set("headless").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	log.Println("Loading page...")
	page := browser.MustIncognito().MustPage(url).MustWaitDOMStable()
	defer page.MustClose()

	title := page.MustElement("h1").MustText()
	subject := strings.SplitAfter(title, ")")
	cardName := subject[0]
	setName := strings.TrimSpace(subject[1])
	log.Printf("Looking Card -> %s", cardName)
	log.Printf("Looking Set -> %s", setName)

	for showMore {
		log.Println("Looking for more sellers")
		showMoreButton, err := page.Timeout(10 * time.Second).Element("#loadMoreButton")
		if err != nil {
			log.Printf("Error Loading more sellers: %s", err.Error())
			showMore = false
		} else {
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
	}

	log.Println("Looking sellers...")
	sellerElements := page.MustElements("div.row.g-0.article-row")

	// pool := rod.NewPagePool(10)
	// // Run jobs concurrently
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go func() {

	// 	defer wg.Done()
	// 	page := pool.MustGet(GetCreatePage(browser))
	// 	defer pool.Put(page)
	// }()
	// wg.Wait()
	// // cleanup pool
	// pool.Cleanup(func(p *rod.Page) { p.MustClose() })

	for _, sEl := range sellerElements {
		seller := &Seller{}

		seller.LoadSellerCardRow(
			true,
			&Card{
				Name:        cardName,
				Url:         url,
				Description: "",
			},
			sEl,
		)

		sellers = append(sellers, seller)
		seller.Round(browser, &wantedList)
	}

	sort.Slice(sellers, func(i, j int) bool {
		return len(sellers[i].CardsAvaialble) > len(sellers[j].CardsAvaialble)
	})
	log.Println("Finished looking for potential sellers:")
	for idx, s := range sellers {
		log.Printf("--- %s %d ---\n", s.Name, idx+1)
		log.Printf("URL:\t%s\n", s.Url)
		log.Printf("Cards for sale (%d/%d):\n", len(s.CardsAvaialble), len(wantedList))
		var total float32 = 0
		for _, ca := range s.CardsAvaialble {
			total += ca.Price
			log.Printf("\tCard(%s(%s) x %d - %0.2f €)\n", ca.Name, ca.Condition, ca.Quantity, ca.Price)
		}
	}
}
