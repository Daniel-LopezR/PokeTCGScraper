package main

import (
	s "poke-tcg-scraper/internal/scraper"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	wantedList := []string{
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Swablu-V2-s12a202?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Drapion-V-V2-s12a227?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Raihan-V2-s12a237?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Grant-V2-s12a238?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Cynthias-Ambition-V2-s12a239?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Cherens-Care-V2-s12a241?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Roxanne-V2-s12a242?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Melony-V2-s12a244?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Volo-V2-s12a245?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Friends-in-Hisui-V2-s12a249?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/VSTAR-Universe/Bosss-Orders-Cyrus-V2-s12a250?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Venusaur-ex-V1-sv2a003?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Beedrill-sv2a015?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Golem-ex-V1-sv2a076?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Mr-Mime-V1-sv2a122?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Ivysaur-V2-sv2a167?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Squirtle-V2-sv2a170?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Pikachu-V2-sv2a173?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Nidoking-V2-sv2a174?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Poliwhirl-V2-sv2a176?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Mr-Mime-V2-sv2a179?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Omanyte-V2-sv2a180?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Venusaur-ex-V2-sv2a184?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Charizard-ex-V2-sv2a185?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Blastoise-ex-V2-sv2a186?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Wigglytuff-ex-V2-sv2a189?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Alakazam-ex-V2-sv2a190?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Kangaskhan-ex-V2-sv2a192?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Jynx-ex-V2-sv2a193?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Zapdos-ex-V2-sv2a194?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Mew-ex-V2-sv2a195?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Erikas-Invitation-V2-sv2a196?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Bills-Transfer-V2-sv2a199?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Venusaur-ex-V3-sv2a200?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Charizard-ex-V3-sv2a201?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Blastoise-ex-V3-sv2a202?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Zapdos-ex-V3-sv2a204?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Erikas-Invitation-V3-sv2a206?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Giovannis-Charisma-V3-sv2a207?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Mew-ex-V4-sv2a208?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Switch-sv2a209?language=7&minCondition=2",
		"https://www.cardmarket.com/en/Pokemon/Products/Singles/Pokemon-Card-151/Basic-Psychic-Energy-sv2a210?language=7&minCondition=2",
	}
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1000, 600, "PokeTCGScraper")
	defer rl.CloseWindow()
	//rl.SetConfigFlags(rl.FlagWindowUndecorated)
	//rl.ToggleFullscreen()
	rl.SetTargetFPS(60)

	scraper := &s.Scraper{
		WantedList: wantedList,
		Sellers: make(map[string]*s.Seller),
	}
	//scraper.InitializeBrowser("brave")
	scraperScene := &s.ScraperScene{
		WantedList: wantedList,
		CurrentCard: &s.Card{
			Name: "None",
		},
		Sellers: make(map[string]*s.Seller),
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
