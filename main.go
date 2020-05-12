package main

import (
	"fmt"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/getlantern/systray"
	"github.com/markbates/pkger"
)

func main() {
	systray.Run(onReady, onExit)
}

func playGong(times int) {
	f, err := pkger.Open("/audio/gong.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	done := make(chan bool)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	for i := 0; i < times; i++ {
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))
		<-done
		streamer.Seek(0)
	}
}

func onReady() {
	systray.SetTitle("Ting Tong")
	systray.SetTooltip("Go Ting Tong")
	mTimeLabel := systray.AddMenuItem("0 mins", "0 mins")
	systray.AddSeparator()
	mBreak := systray.AddMenuItem("Start break timer", "Come on, life's short")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	counter := newCounter()
	go counter.start(1 * time.Minute)

	state := "working"

	for {
		select {
		case count := <-counter.count:
			if state == "working" {
				title := fmt.Sprintf("Working for %d mins", count)
				mTimeLabel.SetTitle(title)
				if (count % 60) == 0 {
					go playGong(count / 60)
				}
			} else if state == "resting" {
				title := fmt.Sprintf("Resting for %d mins", count)
				mTimeLabel.SetTitle(title)
			}
		case <-mTimeLabel.ClickedCh:
			mTimeLabel.SetTitle("Working for 0 mins")
			state = "working"
			counter.reset <- true
		case <-mBreak.ClickedCh:
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
