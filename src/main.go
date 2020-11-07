package main

import (
	"fmt"
	"time"
)

func main() {

	minutesPomodoro := 25
	now := time.Now()
	minuteTimer := time.NewTimer(60 * time.Second)

	fmt.Println("Pomodoro Started:", now)

	for currentMinute := 1; currentMinute <= minutesPomodoro; currentMinute++ {

		<-minuteTimer.C

		// printing remaining time every 5 minutes
		timeMod := currentMinute%5 == 0
		if timeMod {
			remainingMinutes := minutesPomodoro - currentMinute
			println("Remaining minutes: ", remainingMinutes)
		}

	}

}
