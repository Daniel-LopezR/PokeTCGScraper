package main

import (
	"log"
	"net/url"
	"strconv"
	"strings"

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

func (s *Seller) Round(wantedList *[]string) {
	queryURL := "/Offers/Singles?condition=2&name="

	for _, chase := range *wantedList {
		encodedChase := url.QueryEscape(chase)
		fullURL := s.Url + queryURL + encodedChase

		browser := rod.New().MustConnect()
		defer browser.MustClose()

		sellerCardsPage := browser.MustPage(fullURL).MustWaitStable()

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
		// Load Card Url
		sCard.Url = CARDMARKET_URL + *row.MustElement(".col-seller").MustElement("a").MustAttribute("href")
		log.Println("Updated Card Url!")
	}

	// Seller product
	sellerProductInfo := row.MustElement(".col-product")

	sCard.Condition = sellerProductInfo.MustElement("a.article-condition").MustText()
	sCard.Language = *sellerProductInfo.MustElement("span.icon").MustAttribute("aria-label")
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
	price, err := strconv.ParseFloat(strings.Split(sellerProductPrice, " ")[0], 32)
	if err != nil {
		sCard.Price = 0
	} else {
		sCard.Price = float32(price)
	}
	log.Printf("Added new Card(%d x %s(%s) - %0.2f €) to Seller%s\n", sCard.Quantity, sCard.Name, sCard.Condition, sCard.Price, s.Name)

	s.CardsAvaialble = append(s.CardsAvaialble, *sCard)
}

func main() {
	sellers := []*Seller{}
	showMore := true
	// TODO: Loop wanted list instead of hardcoding url
	url := "https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Elesas-Sparkle-V2-s12a246"
	// TODO: import wantedList from e.g. csv, maybe cardmarket api(don't think i will do this tbh)
	wantedList := []string{"Zeraora V (s12a 040)", "Pikachu (s12a 205)", "Radiant Charizard (s12a 15)", "Radiant Greninja (s12a 33)", "Arceus V (s12a 126)", "Irida (s12a 236)", "Water Energy (s12a 253)"}

	browser := rod.New().
		MustConnect()

	defer browser.MustClose()

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
	sellerElements := page.MustElements("div.row.g-0.article-row")
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
		seller.Round(&wantedList)
	}
	log.Printf("Finished looking for potential sellers: \n%v\n", sellers)
}
