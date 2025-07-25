package scraper

import (
	"log"
	"poke-tcg-scraper/internal"
)

type ScraperScene struct {
	id          string
	CurrentCard *Card
	WantedList  []string
}

func (ss *ScraperScene) Update(m internal.Message) {
	switch m.Topic {
	case "card":
		ss.CurrentCard = m.Data.(*Card)
	default:
		log.Printf("Don't care about this event -> %s", m.Topic)
	}
}

func (ss *ScraperScene) GetID() string {
	return ss.id
}

// TODO: Fill the draw function with all the visual stuff
func (ss *ScraperScene) Draw() {

}
