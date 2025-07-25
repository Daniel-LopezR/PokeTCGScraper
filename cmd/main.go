package main

import (
	s "poke-tcg-scraper/internal/scraper"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	wantedListText := "Wanted List"
	wantedList := []string{
		"Swablu (s12a 202)",
		"Drapion V (s12a 227)",
		"Raihan (s12a 237)",
		"Grant (s12a 238)",
		"Cheren's Care (s12a 241)",
		"Roxanne (s12a 242)",
		"Melony (s12a 244)",
		"Volo (s12a 245)",
		"Friends in Hisui (s12a 249)",
		"Boss's Orders (s12a 250)",
		}
	rl.InitWindow(1000, 600, "PokeTCGScraper")
	defer rl.CloseWindow()
	//rl.SetConfigFlags(rl.FlagWindowUndecorated)
	//rl.ToggleFullscreen()
	rl.SetTargetFPS(60)
	
	BASE_X := 100
	BASE_Y := 0
	scraper := &s.Scraper{
		WantedList: wantedList,
	}
	scraperScene := &s.ScraperScene{
		WantedList: wantedList,
	}
	
	// TODO: Modify scraper to notify the observer for every card change
	// TODO: Change how scraper operates, it should go to every card and look every seller not from a card look at every seller for every card
	scraper.Register(scraperScene)
	
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)
		//rl.DrawText(wantedListText, int32(rl.GetScreenWidth()/2)-rl.MeasureText(roundText, 20)/2-rl.MeasureText(player1PointsText, 20)-20, 10, 20, rl.RayWhite)

		rl.DrawText(wantedListText, int32(BASE_X), int32(BASE_Y), 20, rl.Gold)
		for i, wc := range wantedList {
			wcX := BASE_X
			wcY := BASE_Y + ((i+1) * 40)
			rl.DrawText(wc, int32(wcX), int32(wcY), 20, rl.RayWhite)
		}
		rl.EndDrawing()
	}
}
