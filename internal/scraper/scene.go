package scraper

import (
	"fmt"
	"log"
	"poke-tcg-scraper/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
	//rlg "github.com/gen2brain/raylib-go/raygui"
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
	BASE_X := 100
	BASE_Y := 100
	wantedListText := "Wanted List"

	//rl.DrawText(wantedListText, int32(rl.GetScreenWidth()/2)-rl.MeasureText(roundText, 20)/2-rl.MeasureText(player1PointsText, 20)-20, 10, 20, rl.RayWhite)

	rl.DrawText(wantedListText, int32(BASE_X), int32(BASE_Y), 20, rl.Gold)
	for i, wc := range ss.WantedList {
		wcX := BASE_X
		wcY := BASE_Y + ((i+1) * 40)
		rl.DrawText(wc, int32(wcX), int32(wcY), 20, rl.RayWhite)
	}
	
	rl.DrawText(fmt.Sprintf("Current Card: %s", ss.CurrentCard.Name), int32(BASE_X * 3 + 10), int32(BASE_Y), 20, rl.Gold)
}
