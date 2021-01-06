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
	"github.com/ermanimer/color/v2"
	"github.com/kelson-martins/gomodoro/src/db"

	"github.com/ermanimer/progress_bar"
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
		runGomodoro()
	},
}

func runGomodoro() {
	startTime := startDatetime.Format("15:04:05")
	fmt.Println("Gomodoro Started:", startDate, startTime)
	gomodoroProgress()
	finishPomodoro(startTime)
}

func gomodoroProgress() {
	output := os.Stdout
	orange := (&color.Color{Foreground: 172}).SprintFunction()
	grey := (&color.Color{Foreground: 246}).SprintFunction()
	bar := orange("{bar}")
	percent := grey("{percent}")
	schema := fmt.Sprintf("%s %s", bar, percent)
	filledCharacter := "▆"
	blankCharacter := " "
	var length float64 = 50
	var totalValue float64 = float64(gomodoroMinutes * 2)
	pb := progress_bar.NewProgressBar(output, schema, filledCharacter, blankCharacter, length, totalValue)
	// start
	err := pb.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// update
	for value := 1; value <= int(totalValue); value++ {
		time.Sleep(30 * time.Second)
		err := pb.Update(float64(value))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&category, "category", "c", "", "gomodoro category (required)")
	runCmd.MarkPersistentFlagRequired("category")
	runCmd.PersistentFlags().StringVarP(&subcategory, "subcategory", "s", "", "gomodoro subcategory")
	runCmd.PersistentFlags().IntVarP(&gomodoroMinutes, "minutes", "m", 25, "gomodoro minutes session")
	runCmd.MarkPersistentFlagRequired("category")
}

func finishPomodoro(startTime string) {

	dbStartdate := startDatetime.Format("2006-01-02")
	endTime := startDatetime.Add(time.Minute * time.Duration(gomodoroMinutes)).Format("15:04:05")

	fmt.Println("Gomodoro Finished:", startDate, endTime)

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
