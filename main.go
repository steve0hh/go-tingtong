package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func timer() (countChan <-chan int, resetChan chan<- bool, doneChan chan<- bool) {
	count := 0
	c := make(chan int)
	reset := make(chan bool)
	done := make(chan bool)
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
			case <-ticker.C:
				count++
				c <- count
			case <-reset:
				count = 0
			}
		}
	}()
	return c, reset, done
}

func onReady() {
	systray.SetTitle("Ting Tong")
	systray.SetTooltip("Go Ting Tong")
	mTimeLabel := systray.AddMenuItem("0 mins", "0 mins")
	systray.AddSeparator()
	mFiveMin := systray.AddMenuItem("5 Minutes", "Take a break for 5 minutes")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	countChan, resetCountChan, doneChan := timer()

	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		doneChan <- true
		systray.Quit()
		fmt.Println("Finished quiting")
	}()

	go func() {
		for {
			select {
			case count := <-countChan:
				label := fmt.Sprintf("%v mins", count)
				mTimeLabel.SetTitle(label)
				if count > 60 {
					fmt.Println("Time for a break")
				}
			case <-mFiveMin.ClickedCh:
				mTimeLabel.SetTitle("0 mins")
				resetCountChan <- true
			}
		}
	}()
}

func onExit() {
	fmt.Println("On exit")
}
