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

func init() {
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

	records := [][]string{}

	startDatetime := time.Now()
	startDate := startDatetime.Format("2006-01-02")
	currentMonth := startDate[0:7]
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

		if dbMonth == currentMonth {
			monthTotals++
		}
	}

	fmt.Println("[Go]modoros today: ", todayTotals)
	fmt.Println("[Go]modoros yesterday: ", yesterdayTotals)
	fmt.Println("[Go]modoros this month: ", monthTotals)
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
