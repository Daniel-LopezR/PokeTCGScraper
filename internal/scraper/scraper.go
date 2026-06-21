package scraper

import (
	"fmt"
	"log"
	"maps"
	"poke-tcg-scraper/internal"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Scraper struct {
	Launcher   *launcher.Launcher
	Browser    *rod.Browser
	observers []internal.Observer
	// This should be in the config/state
	WantedList []string
	Sellers    map[string]*Seller
}

func (s *Scraper) Register(o internal.Observer) {
	s.observers = append(s.observers, o)
}

func (s *Scraper) Deregister(o internal.Observer) {
	//s.observers = removeFromSlice(s.observers, o)
}

func (s *Scraper) NotifyAll(message internal.Message) {
	for _, observer := range s.observers {
		observer.Update(message)
	}
}

func removeFromSlice(observers []internal.Observer, observerToRemove internal.Observer) []internal.Observer {
	observersLength := len(observers)
	for i, observer := range observers {
		if observerToRemove.GetID() == observer.GetID() {
			observers[observersLength-1], observers[i] = observers[i], observers[observersLength-1]
			return observers[:observersLength-1]
		}
	}
	return observers
}

const CARDMARKET_URL = "https://www.cardmarket.com"

type Card struct {
	Name        string
	Url         string
	ImageName   string
	ImageUrl    *string
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
	CardsAvailable []Card
	TotalPrice	   float32
}

// Create a page if needed
func GetCreatePage(browser *rod.Browser) func() *rod.Page {
	create := func() *rod.Page {
		// We use MustIncognito to isolate pages with each other
		return browser.MustIncognito().MustPage()
	}
	return create
}

func (s *Scraper) GetOrCreateSeller(row *rod.Element) *Seller {
	sellerInfo := row.MustElement(".col-seller")

	sellerLinkElement := sellerInfo.MustElement("a[href]")
	sellerName := sellerLinkElement.MustText()
	for _, sel := range s.Sellers {
		if sel.Name == sellerName {
			return sel
		}
	}
	sel := &Seller{
		Name: sellerName,
	}
	sel.Url = CARDMARKET_URL + *sellerLinkElement.MustAttribute("href")

	sellerRegionAttr := *sellerInfo.MustElement("span.icon").MustAttribute("aria-label")
	sel.Region = strings.TrimSpace(strings.Split(sellerRegionAttr, ":")[1])

	if sellerInfo.MustHas("span.fonticon-users-powerseller") {
		sel.Category = "Powerseller"
	} else if sellerInfo.MustHas("span.fonticon-users-professional") {
		sel.Category = "Professional"
	} else {
		sel.Category = "Hobby"
	}

	log.Printf("Added new Seller%+v\n", s)
	return sel
}

func (s *Seller) LoadSellerCardRow(loadSeller bool, sCard *Card, row *rod.Element) {

	// cardInfo := row.MustElement(".col-seller")
	// cardLinkElement := cardInfo.MustElement("a[href]")
	// cardName := cardLinkElement.MustText()
	// if cardName == sCard.Name {
	// 	// Load Card Url
	// 	sCard.Url = CARDMARKET_URL + *row.MustElement(".col-seller").MustElement("a").MustAttribute("href")
	// 	log.Println("Updated Card Url!")
	// } else {
	// 	return
	// }
	
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
		cardPrice := float32(price)
		sCard.Price = cardPrice
		s.TotalPrice += cardPrice
	}
	log.Printf("Added new Card(%d x %s(%s) - %0.2f €) to Seller %s\n", sCard.Quantity, sCard.Name, sCard.Condition, sCard.Price, s.Name)

	s.CardsAvailable = append(s.CardsAvailable, *sCard)
}

func (s *Scraper) Scrap() {
	now := time.Now()
	log.Println("Creating browser...")
	err := s.InitializeBrowser("chrome")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer s.Browser.MustClose()

	for _, wURL := range s.WantedList {
		showMore := true
		log.Println("Loading page...")
		page := s.Browser.MustPage(wURL).MustWaitDOMStable()
		defer page.MustClose()

		card := &Card{
			Url: wURL,
		}

		title := page.MustElement("h1").MustText()
		subject := strings.SplitAfter(title, ")")
		card.Name = subject[0]
		setName := strings.TrimSpace(subject[1])
		log.Printf("Looking Card -> %s", card.Name)
		log.Printf("Looking Set -> %s", setName)

		//Image
		imageElement := page.MustElement("#image").MustElement("img")
		imageUrl, err := imageElement.Attribute("src")
		if err != nil {
			panic(err)
		}
		card.ImageUrl = imageUrl
		
		imageName := ""
		cardNameSections := strings.Split(card.Name, "(")
		if len(cardNameSections) > 0 {
			idSec := cardNameSections[1]
			baseImageName := strings.Replace(
				strings.Replace(
					strings.Replace(idSec, "(", "", 1),
					")", "", 1,
				),
				" ", "_", 1,
			)

			imageExtension := ".jpg"
			if imageUrl != nil {
				imageExtension = (*imageUrl)[strings.LastIndex(*imageUrl, "."):]
			}
			imageName = baseImageName + imageExtension
		}
		card.ImageName = imageName

		s.NotifyAll(internal.Message{
			Topic: "card",
			Data:  card,
		})

		buttonIdx := 0
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
					//showMoreButton.MustWaitStable()
					showMoreButton.MustClick()
					_,err := showMoreButton.Timeout(3 * time.Second).WaitInteractable()
					if err != nil {
						maxResReached, _ := page.Timeout(1 * time.Second).Element("div#MaxResultsReachedNotice")
						if maxResReached != nil {
							log.Println("Max ammount of sellers reached...")
							showMore = false
							continue
						}
						// This error should mean that the button is now invisible
						log.Println("Couldn't wait button to be stable")
						continue
					}
				} else {
					log.Printf("Show More Results - disabled attr(%s)", *disabled)
					log.Println("No more sellers to show!")
					showMore = false
				}
			}
			buttonIdx++
		}
		log.Println("Looking sellers...")
		sellerElements := page.MustElements("div.row.g-0.article-row")

		for _, sEl := range sellerElements {
			// Seller name
			sellerInfo := sEl.MustElement(".col-seller")
			sellerName := sellerInfo.MustElement("a").MustText()

			if _, ok := s.Sellers[sellerName]; !ok {
				newSeller := &Seller{
					Name: sellerName,
				}
				s.Sellers[sellerName] = newSeller
			}
			s.Sellers[sellerName].LoadSellerCardRow(
				true,
				card,
				sEl,
			)

			//seller.Round(browser, &wantedList)
		}
		sellersToSort := slices.Collect(maps.Values(s.Sellers))
		sort.Slice(sellersToSort, func(i, j int) bool {
			if len(sellersToSort[i].CardsAvailable) == len(sellersToSort[j].CardsAvailable) {
				return sellersToSort[i].TotalPrice > sellersToSort[j].TotalPrice
			}
			return len(sellersToSort[i].CardsAvailable) > len(sellersToSort[j].CardsAvailable)
		})
		sortedSellerMap := make(map[string]*Seller, len(s.Sellers))
		for _,sSorted := range sellersToSort {
			sortedSellerMap[sSorted.Name] = sSorted
		}
		s.Sellers = sortedSellerMap

		// Notify Scene with the sorted map
		s.NotifyAll(internal.Message{
			Topic: "sellers",
			Data:  s.Sellers,
		})
	}

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


	log.Println("Finished looking for potential sellers:")
	limit := 50
	count := 0
	for sellerName, sellerInfo := range s.Sellers {
		if count > limit {
			break
		}
		log.Printf("--- %s %d ---\n", sellerName, count+1)
		log.Printf("URL:\t%s\n", sellerInfo.Url)
		log.Printf("Cards for sale (%d/%d) at %.2f€:\n", len(sellerInfo.CardsAvailable), len(s.WantedList), sellerInfo.TotalPrice)
		for _, ca := range sellerInfo.CardsAvailable {
			log.Printf("\tCard(%s(%s) x %d - %0.2f €)\n", ca.Name, ca.Condition, ca.Quantity, ca.Price)
		}
		count++
	}
	log.Println("Finished after ", time.Since(now))
}
