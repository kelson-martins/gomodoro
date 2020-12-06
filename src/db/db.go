package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// var gomodoroDB, _ = sql.Open("sqlite3", os.Getenv("HOME")+"/gomodoro/G.db")

var gomodoroDB *sql.DB

type Gomodoro struct {
	Year           string
	Date           string
	StartTimestamp string
	EndTimestamp   string
	Minutes        int
	Category       string
	SubCategory    string
}

func InsertRecord(date string, startTimestamp string, endTimestamp string, minutes int, category string, subcategory string) {
	insertGomodoroSQL := `INSERT INTO gomodoros(date, startTimestamp, endTimestamp, minutes, category, subCategory) VALUES (?,?,?,?,?,?)`
	statement, err := gomodoroDB.Prepare(insertGomodoroSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(date, startTimestamp, endTimestamp, minutes, category, subcategory)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func createTable() {
	createTableGomodoros := `CREATE TABLE IF NOT EXISTS gomodoros (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"date" DATE NOT NULL,		
		"startTimestamp" TIME,
		"endTimestamp" TIME,
		"minutes" INT NOT NULL,
		"category" TEXT NOT NULL,
		"subCategory" TEXT
	  );`

	statement, err := gomodoroDB.Prepare(createTableGomodoros)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func Init() {
	LoadConfig()
	dbPath := viper.GetString("gomodoro.databasePath")
	sqliteDir := strings.Replace(dbPath, "$HOME", os.Getenv("HOME"), -1)
	gomodoroDB, _ = sql.Open("sqlite3", sqliteDir)
	createTable()
}

func Close() {
	gomodoroDB.Close()
}

func DeleteLatest() {

	var toDelete int

	latestStatement := `SELECT id FROM gomodoros ORDER BY id DESC LIMIT 1;`

	rows, err := gomodoroDB.Query(latestStatement)

	if err != nil {
		log.Fatalln(err.Error())
	}
	for rows.Next() {
		rows.Scan(&toDelete)
	}

	if toDelete > 0 {
		deleteID(toDelete)
	}

}

func ExternalDeleteID(id int) {
	deleteID(id)
}

func deleteID(id int) {

	deleteStatement := `DELETE FROM gomodoros WHERE id = ?`
	statement, err := gomodoroDB.Prepare(deleteStatement)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("[Go]modoro ID %v was deleted", id)
}

func GetTotalsRecord() {

	var allTimeTotals int
	var todayTotals int
	var yesterdayTotals int
	var monthTotals int
	var lastMonthTotals int

	records := [][]string{}

	startDatetime := time.Now()
	startDate := startDatetime.Format("2006-01-02")
	currentMonth := startDate[0:7]
	lastMonth := startDatetime.AddDate(0, -1, 0).Format("2006-01-02")[0:7]
	yesterdayDate := startDatetime.AddDate(0, 0, -1).Format("2006-01-02")

	rows, err := gomodoroDB.Query(`SELECT strftime(date) as date, substr(date,1,7) as monthYear,substr(date,1,4) as year, startTimestamp, category, subcategory FROM gomodoros`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for rows.Next() {
		var date, monthYear, year, startTimestamp, category, subcategory string
		rows.Scan(&date, &monthYear, &year, &startTimestamp, &category, &subcategory)

		data := []string{date, monthYear, year, startTimestamp, category, subcategory}
		records = append(records, data)

		allTimeTotals++
	}

	for _, row := range records {
		dbDate, dbMonth := row[0], row[1]

		switch dbDate {
		case startDate:
			todayTotals++
		case yesterdayDate:
			yesterdayTotals++
		}

		switch dbMonth {
		case currentMonth:
			monthTotals++
		case lastMonth:
			lastMonthTotals++
		}

	}

	fmt.Println("[Go]modoros today: ", todayTotals)
	fmt.Println("[Go]modoros yesterday: ", yesterdayTotals)
	fmt.Println("[Go]modoros this month: ", monthTotals)
	fmt.Println("[Go]modoros last month: ", lastMonthTotals)
	fmt.Println("[Go]modoros all-time: ", allTimeTotals)
}

func ListGomodoros(days int) {

	queryStatement := fmt.Sprintf("SELECT id, strftime(date) as date, startTimestamp, category, subCategory from gomodoros WHERE date > date('now', '-%v day') ORDER BY id DESC", days)
	rows, err := gomodoroDB.Query(queryStatement)

	if err != nil {
		fmt.Println("error")
		log.Fatalln(err.Error())
	}

	for rows.Next() {
		var id, date, startTimestamp, category, subcategory string
		rows.Scan(&id, &date, &startTimestamp, &category, &subcategory)

		if subcategory != "" {
			fmt.Printf("ID: %v\t %v %v\t Category: %v\t SubCategory: %v\n", id, date, startTimestamp, category, subcategory)
		} else {
			fmt.Printf("ID: %v\t %v %v\t Category: %v\n", id, date, startTimestamp, category)
		}

	}

}

func GetGomodorosCassandra() []Gomodoro {

	queryStatement := "SELECT substr(date,0,5) as year, date, startTimestamp, endTimestamp, minutes, category, subCategory from gomodoros"
	listGomodoros := []Gomodoro{}

	rows, err := gomodoroDB.Query(queryStatement)
	if err != nil {
		log.Fatalf("[ERROR] error fetching Gomodoros from database, %v", err.Error())
	}

	for rows.Next() {
		var year, date, startTimestamp, endTimestamp, category, subcategory string
		var minutes int

		rows.Scan(&year, &date, &startTimestamp, &endTimestamp, &minutes, &category, &subcategory)

		entry := Gomodoro{
			Year:           year,
			Date:           date,
			StartTimestamp: startTimestamp,
			EndTimestamp:   endTimestamp,
			Minutes:        minutes,
			Category:       category,
			SubCategory:    subcategory,
		}

		listGomodoros = append(listGomodoros, entry)
	}

	return listGomodoros
}

func GetValidYears() []int {

	yearsQuery := "SELECT substr(date,0,5) as year from gomodoros  GROUP BY year"
	validYears := []int{}

	rows, err := gomodoroDB.Query(yearsQuery)
	if err != nil {
		log.Fatalf("[ERROR] error fetching Gomodoros data from database, %v", err.Error())
	}

	for rows.Next() {
		var year int

		rows.Scan(&year)

		validYears = append(validYears, year)
	}

	return validYears
}

func GetLocalGomodorosByYear(year int) []string {
	gomodorosQuery := fmt.Sprintf("SELECT date, startTimestamp from gomodoros  WHERE substr(date,0,5) = '%v'", year)

	localGomodoros := []string{}

	rows, err := gomodoroDB.Query(gomodorosQuery)
	if err != nil {
		log.Fatalf("[ERROR] error fetching Gomodoros data from database, %v", err.Error())
	}

	for rows.Next() {
		var date, startTimestamp string
		rows.Scan(&date, &startTimestamp)

		v := fmt.Sprintf("%v %v", date[0:10], startTimestamp)

		localGomodoros = append(localGomodoros, v)
	}

	return localGomodoros
}
