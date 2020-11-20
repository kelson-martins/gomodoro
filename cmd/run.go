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

	"database/sql"
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

var applicationPath = "/etc/gomodoro/"

var gomodoroDB, _ = sql.Open("sqlite3", os.Getenv("HOME") + "/gomodoro/gomodoro.db")
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

		defer gomodoroDB.Close()
		initDB(gomodoroDB)
	
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

		if pomodoroMinute == gomodoroMinutes + 1 {
			ticker.Stop()
		}

	}
}

func finishPomodoro(startTime string) {

	endTime := startDatetime.Add(time.Minute * time.Duration(gomodoroMinutes)).Format("15:04:05")

	fmt.Println("Pomodoro Finished:", startDate, endTime)

	f, err := os.Open("/etc/gomodoro/tone.mp3")
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

	insertRecord(gomodoroDB, startDate, startTime, endTime, gomodoroMinutes, category, subcategory)

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

func initDB(db *sql.DB) {
	createTable(db)
}

func createTable(db *sql.DB) {
	createTableGomodoros := `CREATE TABLE IF NOT EXISTS gomodoros (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"date" DATE NOT NULL,		
		"startTimestamp" TIME,
		"endTimestamp" TIME,
		"minutes" INT NOT NULL,
		"category" TEXT NOT NULL,
		"subCategory" TEXT
	  );`

	statement, err := db.Prepare(createTableGomodoros)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func insertRecord(db *sql.DB, date string, startTimestamp string, endTimestamp string, minutes int, category string, subcategory string) {
	insertGomodoroSQL := `INSERT INTO gomodoros(date, startTimestamp, endTimestamp, minutes, category, subCategory) VALUES (?,?,?,?,?,?)`
	statement, err := db.Prepare(insertGomodoroSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(date, startTimestamp, endTimestamp, minutes, category, subcategory)
	if err != nil {
		log.Fatalln(err.Error())
	}
}