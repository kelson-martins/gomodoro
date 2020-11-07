package main

import (
	"fmt"
	"time"
)

var pomodoroMinute int = 1
var totalMinutes int = 25

func main() {

	now := time.Now()
	fmt.Println("Pomodoro Started:", now)

	go pomodoroHeartbeat()
	time.Sleep(26 * time.Second)

}

func pomodoroHeartbeat() {
	for range time.Tick(1 * time.Second) {
		timeMod := pomodoroMinute%5 == 0
		if timeMod == true && pomodoroMinute != 25 {
			println("Remaining minutes: ", totalMinutes-pomodoroMinute)
		}
		pomodoroMinute++

		if pomodoroMinute == 26 {
			now := time.Now()
			fmt.Println("Pomodoro Finished:", now)
		}
	}
}
