package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var pomodoroMinute int = 1
var totalMinutes int = 25

func main() {

	startTime := time.Now()
	fmt.Println("Pomodoro Started:", startTime)

	go pomodoroHeartbeat()
	time.Sleep(25 * time.Minute)
	finishPomodoro()
}

func pomodoroHeartbeat() {
	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		timeMod := pomodoroMinute%5 == 0
		if timeMod == true && pomodoroMinute != 25 {
			println("Remaining minutes:", totalMinutes-pomodoroMinute)
		}
		pomodoroMinute++

		if pomodoroMinute == 26 {
			ticker.Stop()
		}

	}
}

func finishPomodoro() {

	endTime := time.Now()
	fmt.Println("Pomodoro Finished:", endTime)

	f, err := os.Open("../assets/tone.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	displayEnd()

}

func displayEnd() {
	a := app.New()
	w := a.NewWindow("Gomodoro")

	text := widget.NewLabel("Gomodoro Finished")
	w.SetContent(widget.NewVBox(
		text,
	))

	w.ShowAndRun()
}
