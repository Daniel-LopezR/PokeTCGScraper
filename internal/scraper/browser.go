package scraper

import (
	"context"
	"os/exec"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// type Browser struct {
// 	Type string
// }

func (s *Scraper) InitializeBrowser(bType string) (error) {
	path,_ := exec.LookPath(bType)
	s.Launcher = launcher.New().Bin(path).Headless(true).Leakless(false)
	s.Browser = rod.New().Context(context.Background()).ControlURL(s.Launcher.MustLaunch()).MustConnect().MustIncognito()
	return nil
}
