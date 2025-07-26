package main

import (
	s "poke-tcg-scraper/internal/scraper"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
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
	rl.InitWindow(1000, 600, "PokeTCGScraper")
	defer rl.CloseWindow()
	//rl.SetConfigFlags(rl.FlagWindowUndecorated)
	//rl.ToggleFullscreen()
	rl.SetTargetFPS(60)

	scraper := &s.Scraper{
		WantedList: wantedList,
	}
	scraperScene := &s.ScraperScene{
		WantedList: wantedList,
		CurrentCard: &s.Card{
			Name: "None",
		},
	}

	// TODO: Modify scraper to notify the observer for every card change
	// TODO: Change how scraper operates, it should go to every card and look every seller not from a card look at every seller for every card
	scraper.Register(scraperScene)
	go scraper.Scrap()
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)
		scraperScene.Draw()
		rl.EndDrawing()
	}
}
