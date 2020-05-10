package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Ting Tong")
	systray.SetTooltip("Go Ting Tong")
	mTimeLabel := systray.AddMenuItem("0 mins", "0 mins")
	systray.AddSeparator()
	mFiveMin := systray.AddMenuItem("Take a 5 minutes break", "Come on, life's short")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	counter := newCounter()
	go counter.start(1 * time.Second)

	state := "working"

	for {
		select {
		case count := <-counter.count:
			if state == "working" {
				title := fmt.Sprintf("Working for %d mins", count)
				mTimeLabel.SetTitle(title)
				if (count % 60) == 0 {
					for i := 0; i < (count / 60); i++ {
						fmt.Println("TAKE A BREAK!")
					}
				}
			} else if state == "resting" {
				title := fmt.Sprintf("Resting for %d mins", count)
				mTimeLabel.SetTitle(title)
			}
		case <-mTimeLabel.ClickedCh:
			mTimeLabel.SetTitle("Working for 0 mins")
			state = "working"
			counter.reset <- true
		case <-mFiveMin.ClickedCh:
			mTimeLabel.SetTitle("Resting for 0 mins")
			state = "resting"
			counter.reset <- true
		case <-mQuitOrig.ClickedCh:
			counter.done <- true
			systray.Quit()
		}
	}
}

func onExit() {
	fmt.Println("On exit")
}
