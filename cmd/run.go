/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"db"

	_ "github.com/mattn/go-sqlite3"
)

var configFilesPath = "/etc/gomodoro/"
var startDatetime = time.Now()
var startDate = startDatetime.Format("02-01-2006")
var category string
var subcategory string
var gomodoroMinutes int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a [Go]modoro",
	Long: `The 'run' command is used to start a new [go]modoro session.

Examples:

go run --category coding
`,
	Run: func(cmd *cobra.Command, args []string) {
		db.Init()
		defer db.Close()
		startTime := startDatetime.Format("15:04:05")
		fmt.Println("Pomodoro Started:", startDate, startTime)
		go pomodoroHeartbeat()
		time.Sleep(time.Duration(gomodoroMinutes) * time.Minute)
		finishPomodoro(startTime)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&category, "category", "c", "", "gomodoro category (required)")
	runCmd.MarkPersistentFlagRequired("category")
	runCmd.PersistentFlags().StringVarP(&subcategory, "subcategory", "s", "", "gomodoro subcategory")
	runCmd.PersistentFlags().IntVarP(&gomodoroMinutes, "minutes", "m", 25, "gomodoro minutes session")
	runCmd.MarkPersistentFlagRequired("category")
}

func pomodoroHeartbeat() {
	pomodoroMinute := 1
	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		timeMod := pomodoroMinute%5 == 0
		if timeMod == true && pomodoroMinute != gomodoroMinutes {
			println("Remaining minutes:", gomodoroMinutes-pomodoroMinute)
		}
		pomodoroMinute++

		if pomodoroMinute == gomodoroMinutes+1 {
			ticker.Stop()
		}

	}
}

func finishPomodoro(startTime string) {

	dbStartdate := startDatetime.Format("2006-01-02")
	endTime := startDatetime.Add(time.Minute * time.Duration(gomodoroMinutes)).Format("15:04:05")

	fmt.Println("Pomodoro Finished:", startDate, endTime)

	f, err := os.Open(configFilesPath + "tone.mp3")
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

	db.InsertRecord(dbStartdate, startTime, endTime, gomodoroMinutes, category, subcategory)

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
