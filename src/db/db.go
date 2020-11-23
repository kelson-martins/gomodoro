package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var gomodoroDB, _ = sql.Open("sqlite3", os.Getenv("HOME")+"/gomodoro/gomodoro.db")

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
	var monthTotals int

	records := [][]string{}

	startDatetime := time.Now()
	startDate := startDatetime.Format("02-01-2006")
	currentMonth := startDate[3:]

	rows, err := gomodoroDB.Query(`SELECT strftime(date) as date, substr(date,4) as monthYear,substr(date,7) as year, startTimestamp, category, subcategory FROM gomodoros`)
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
		if dbDate == startDate {
			todayTotals++
		}

		if dbMonth == currentMonth {
			monthTotals++
		}
	}

	fmt.Println("[Go]modoros today: ", todayTotals)
	fmt.Println("[Go]modoros month: ", monthTotals)
	fmt.Println("[Go]modoros all-time: ", allTimeTotals)
}

func ListGomodoros(days int) {

	queryStatement := "SELECT id, strftime(date) as date, startTimestamp, category, subCategory from gomodoros WHERE date > date('now', '-? day') ORDER BY id DESC"

	rows, err := gomodoroDB.Query(queryStatement, days)

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
