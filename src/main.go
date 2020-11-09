package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	time.Sleep(25 * time.Second)
	finishPomodoro()
}

func pomodoroHeartbeat() {
	for range time.Tick(1 * time.Second) {
		timeMod := pomodoroMinute%5 == 0
		if timeMod == true && pomodoroMinute != 25 {
			println("Remaining minutes:", totalMinutes-pomodoroMinute)
		}
		pomodoroMinute++
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

}
