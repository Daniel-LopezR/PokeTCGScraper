package scraper

import (
	"os/exec"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// type Browser struct {
// 	Type string
// }

func (s *Scraper) InitializeBrowser(bType string) (error) {
	// bChan := make(chan string)
	// switch bType {
	// case "brave":
	// 	cmdHB := exec.Command(bType, "--headless", "--remote-debugging-port=9222", "--no-sandbox", "--disable-gpu")
	// 	cmdHBOut, err := cmdHB.StderrPipe()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if err = cmdHB.Start(); err != nil {
	// 		return err
	// 	}
	// 	go func() {
	// 		defer close(bChan)
	// 		scanner := bufio.NewScanner(cmdHBOut)
	// 		for scanner.Scan() {
	// 			line := scanner.Text()
	// 			if strings.Contains(line, "DevTools listening on") {
	// 				parts := strings.Fields(line)
	// 				if len(parts) > 3 {
	// 					bChan <- parts[3]
	// 				}
	// 			}
	// 		}
	// 		//cmdHB.Wait()
	// 		// Should do something with the open process
	// 	}()
	// default:
	// 	log.Println("None")
	// 	return fmt.Errorf("None")
	// }
	// wsURL := <-bChan
	// log.Println(wsURL)
	path, _ := exec.LookPath(bType)
	s.Launcher = launcher.New().Bin(path).Headless(true)
	s.Browser = rod.New().ControlURL(s.Launcher.MustLaunch()).MustConnect().MustIncognito()
	return nil
}
