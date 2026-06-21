package scraper

import (
	"errors"
	"fmt"
	"image/color"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"poke-tcg-scraper/internal"
	"syscall"

	rl "github.com/gen2brain/raylib-go/raylib"
	//rlg "github.com/gen2brain/raylib-go/raygui"
)

type ScraperScene struct {
	id          	 	string
	CurrentCard 	 	*Card
	CurrentCardTexture  *rl.Texture2D
	WantedList  	 	[]string
	Sellers				map[string]*Seller
}

func (ss *ScraperScene) Update(m internal.Message) {
	switch m.Topic {
	case "card":
		ss.CurrentCard = m.Data.(*Card)
		if ss.CurrentCard.ImageUrl != nil {
			//ss.CurrentCardTexture = saveImageFromURL(ss.CurrentCard.ImageName, *ss.CurrentCard.ImageUrl)
		}
	case "sellers":
		println("UPDATING SELLERS!!!")
		ss.Sellers = maps.Clone(m.Data.(map[string]*Seller))
	default:
		log.Printf("Don't care about this event -> %s", m.Topic)
	}
}

func saveImageFromURL(name string, url string) *rl.Texture2D {
	var fileData []byte
	imagePath := "./data/" + name
	f, err := os.OpenFile(imagePath, os.O_RDONLY, os.ModeAppend.Perm())
	if err != nil {
		if errors.Is(err, syscall.ERROR_FILE_NOT_FOUND) {
			newF, err := os.Create(imagePath)
			if err != nil {
				log.Fatal(err)
			}
			defer newF.Close()
			println("IMAGE URL: ", url)
			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			fileData, err = io.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			_,err = newF.Write(fileData)
			if err != nil {
				log.Fatal(err)
			}
			err = newF.Sync()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("UNEXPECTED ERROR OPENING FILE: ", err)
		}
	} else {
		defer f.Close()
		fileData, err = io.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
	}
	image := rl.LoadImage(imagePath)
	if image != nil {
		texture := rl.LoadTextureFromImage(image)
		return &texture
	}
	return nil
}

func (ss *ScraperScene) GetID() string {
	return ss.id
}

// TODO: Fill the draw function with all the visual stuff
func (ss *ScraperScene) Draw() {
	BASE_X := int32(100)
	BASE_Y := int32(100)
	wantedListText := fmt.Sprintf("Wanted List - Total Sellers(%d)", len(ss.Sellers)) 

	//rl.DrawText(wantedListText, int32(rl.GetScreenWidth()/2)-rl.MeasureText(roundText, 20)/2-rl.MeasureText(player1PointsText, 20)-20, 10, 20, rl.RayWhite)

	rl.DrawText(wantedListText, int32(BASE_X), int32(BASE_Y), 20, rl.Gold)
	count := 0
	for sellerName, sellerInfo := range ss.Sellers {
		wcX := BASE_X
		wcY := BASE_Y + int32((count+1) * 40)
		rl.DrawText(fmt.Sprintf("%d) %s - %.2f€(%d)", count + 1, sellerName, sellerInfo.TotalPrice, len(sellerInfo.CardsAvailable)), int32(wcX), int32(wcY), 20, rl.RayWhite)
		count++
	}

	if ss.CurrentCardTexture != nil {
		rl.DrawTexture(*ss.CurrentCardTexture, BASE_X, BASE_Y + 40, color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		})
	}
	
	rl.DrawText(fmt.Sprintf("Current Card: %s", ss.CurrentCard.Name), int32(BASE_X * 5 + 10), int32(BASE_Y), 20, rl.Gold)
}
