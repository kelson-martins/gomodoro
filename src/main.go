package main

import (
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

var pomodoroMinute int = 1
var gomodoroMinutes int = 25
var gomodoroDB, _ = sql.Open("sqlite3", "./gomodoro.db")
var startDatetime = time.Now()
var startDate = startDatetime.Format("02-01-2006")

func main() {

	defer gomodoroDB.Close()
	initDB(gomodoroDB)

	startTime := startDatetime.Format("15:04:05")
	fmt.Println("Pomodoro Started:", startDate, startTime)

	go pomodoroHeartbeat()
	time.Sleep(25 * time.Second)

	finishPomodoro(startTime)

}

func pomodoroHeartbeat() {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		timeMod := pomodoroMinute%5 == 0
		if timeMod == true && pomodoroMinute != 25 {
			println("Remaining minutes:", gomodoroMinutes-pomodoroMinute)
		}
		pomodoroMinute++

		if pomodoroMinute == 26 {
			ticker.Stop()
		}

	}
}

func finishPomodoro(startTime string) {

	endTime := startDatetime.Add(time.Minute * 25).Format("15:04:05")

	fmt.Println("Pomodoro Finished:", startDate, endTime)

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

	insertRecord(gomodoroDB, startDate, startTime, endTime, gomodoroMinutes, "", "")

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
