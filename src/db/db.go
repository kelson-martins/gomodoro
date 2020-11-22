package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	// var allTimeTotals int
	// var todayTotals int
	// var year int

	records := [][]string{}

	// startDatetime := time.Now()
	// startDate := startDatetime.Format("02-01-2006")

	rows, err := gomodoroDB.Query(`SELECT strftime(date) as date, substr(date,7) as year, startTimestamp, category, subcategory FROM gomodoros`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for rows.Next() {

		var date, year, startTimestamp, category, subcategory string
		rows.Scan(&date, &year, &startTimestamp, &category, &subcategory)

		data := []string{date, year, startTimestamp, category, subcategory}
		records := append(records, data)
		fmt.Println(records)
		// todayTotals = today
	}

	for row, _ := range records {
		fmt.Println(row)
	}
	// rows2, err := gomodoroDB.Query(`SELECT COUNT(id) as total FROM gomodoros`)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	// for rows2.Next() {
	// 	var total int
	// 	rows2.Scan(&total)

	// 	allTimeTotals = total
	// }

	// fmt.Println("[Go]modoros today: ", todayTotals)
	// fmt.Println("[Go]modoros all-time: ", allTimeTotals)
}
