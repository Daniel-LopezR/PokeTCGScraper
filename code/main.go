package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)
const CARDMARKET_URL = "https://www.cardmarket.com"

type Scraper struct {
	Launcher   *launcher.Launcher
	Browser    *rod.Browser
	// This should be in the config/state
	WantedList []string
	Sellers    []*Seller
}
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

func GetCreatePage(browser *rod.Browser) func() *rod.Page {
	create := func() *rod.Page {
		// We use MustIncognito to isolate pages with each other
		return browser.MustIncognito().MustPage()
	}
	return create
}

func (s *Scraper) InitializeBrowser(bType string) (error) {
	path, _ := exec.LookPath(bType)
	s.Launcher = launcher.New().Bin(path).Headless(true)
	s.Browser = rod.New().ControlURL(s.Launcher.MustLaunch()).MustConnect().MustIncognito()
	return nil
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

	//log.Printf("Added new Seller%+v\n", s)
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
		sCard.Price = float32(price)
	}
	//log.Printf("Added new Card(%d x %s(%s) - %0.2f €) to Seller %s\n", sCard.Quantity, sCard.Name, sCard.Condition, sCard.Price, s.Name)

	s.CardsAvaialble = append(s.CardsAvaialble, *sCard)
}

func (s *Scraper) Scrap(writer io.Writer) {
	//log.Println("Creating browser...")
	err := s.InitializeBrowser("brave")
	if err != nil {
		//log.Println(err)
		panic(err)
	}
	defer s.Browser.MustClose()

	for _, wURL := range s.WantedList {
		showMore := true
		//log.Println("Loading page...")
		page := s.Browser.MustPage(wURL).MustWaitDOMStable()
		defer page.MustClose()

		title := page.MustElement("h1").MustText()
		subject := strings.SplitAfter(title, ")")
		cardName := subject[0]
		setName := strings.TrimSpace(subject[1])
		_ = setName
		//log.Printf("Looking Card -> %s", cardName)
		//log.Printf("Looking Set -> %s", setName)
		card := &Card{
			Name:        cardName,
			Url:         wURL,
			Description: "",
		}

		for showMore {
			//log.Println("Looking for more sellers")
			showMoreButton, err := page.Timeout(10 * time.Second).Element("#loadMoreButton")
			if err != nil {
				//log.Printf("Error Loading more sellers: %s", err.Error())
				showMore = false
			} else {
				attrs := showMoreButton.MustDescribe().Attributes
				_ = attrs
				//log.Printf("Attributes %v", attrs)
				disabled, err := showMoreButton.Attribute("disabled")
				if err != nil {
					panic(err)
				}
				if disabled == nil {
					//log.Println("Showing more sellers...")
					showMoreButton.MustWaitStable().MustClick().MustWaitInvisible()
				} else {
					//log.Printf("Show More Results - disabled attr(%s)", *disabled)
					//log.Println("No more sellers to show!")
					showMore = false
				}
			}
		}
		//log.Println("Looking sellers...")
		sellerElements := page.MustElements("div.row.g-0.article-row")

		for _, sEl := range sellerElements {
			seller := &Seller{}

			seller.LoadSellerCardRow(
				true,
				card,
				sEl,
			)

			s.Sellers = append(s.Sellers, seller)
			//seller.Round(browser, &wantedList)
		}
		
	}
}

func main(){
	writer := bufio.NewWriterSize(os.Stdout, 4096)
	reader := bufio.NewReader(os.Stdin)
	wantedList := []string{
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Cynthias-Ambition-V2-s12a239?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Swablu-V2-s12a202?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Drapion-V-V2-s12a227?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Raihan-V2-s12a237?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Grant-V2-s12a238?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Cherens-Care-V2-s12a241?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Roxanne-V2-s12a242?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Melony-V2-s12a244?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Volo-V2-s12a245?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Friends-in-Hisui-V2-s12a249?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Bosss-Orders-Cyrus-V2-s12a250?language=7&minCondition=2",
	}

	scraper := &Scraper{
		WantedList: wantedList,
	}

	for {
		readableBytes, _ := reader.Peek(512)
		if len(readableBytes) > 0 {
			//fmt.Printf("Peeked %d bytes(%v)\n", len(readableBytes), string(readableBytes))
			command, err := reader.ReadString('\n')
			if err != nil {
				command = err.Error()
			} else {
				switch command {
				case "SCRAP":
					scraper.Scrap(writer)
				}	
			}
			writer.WriteString(fmt.Sprintf("[Scraper] Recieved command from Parent: %s", command))
			writer.Flush()
		}
	}
}
