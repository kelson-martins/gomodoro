package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type astra_config struct {
	dbPath        string
	clusterId     string
	clusterRegion string
	keyspace      string
	table         string
	token         string
}

type gomodoro struct {
	Date            string
	End_timestamp   string
	Start_timestamp string
	Sub_category    string
	Minutes         int
	Category        string
}

var timeout = time.Duration(5 * time.Second)
var client = http.Client{Timeout: timeout}
var config astra_config

type gomodoros struct {
	Rows []gomodoro `json:"rows"`
}

func Sync() {
	token := config.token
	createCassandraTable(token)
	pushToRemote(token)
	readRemote(token)
}

func ValidateSync() {

	if isStringEmpty(config.dbPath, config.token, config.keyspace, config.clusterId, config.clusterRegion) == true {
		log.Println("[ERROR] configure cassandra details at config [$HOME/gomodoro/config.yaml] to use the sync functionality")
		os.Exit(0)
	}

}

func LoadConfig() {
	config = astra_config{
		dbPath:        strings.Replace(viper.GetString("gomodoro.databasePath"), "$HOME", os.Getenv("HOME"), -1),
		keyspace:      viper.GetString("gomodoro.cassandra.keyspace"),
		clusterId:     viper.GetString("gomodoro.cassandra.cluster_id"),
		clusterRegion: viper.GetString("gomodoro.cassandra.cluster_region"),
		token:         viper.GetString("gomodoro.cassandra.token"),
		table:         "gomodoros",
	}

}

func GetDbPath() string {
	return config.dbPath
}

func isStringEmpty(values ...string) bool {

	for _, v := range values {
		if v == "" {
			return true
		}
	}

	return false
}

func createCassandraTable(token string) {

	postURL := fmt.Sprintf("https://%v-%v.apps.astra.datastax.com/api/rest/v1/keyspaces/%v/tables", config.clusterId, config.clusterRegion, config.keyspace)

	type gomodoroTable struct {
		Name              string              `json:"name"`
		IfNotExists       bool                `json:"ifNotExists"`
		TableOptions      map[string]int      `json:"tableOptions"`
		ColumnDefinitions []map[string]string `json:"columnDefinitions"`
		PrimaryKey        map[string][]string `json:"primaryKey"`
	}

	columnDefinitions := []map[string]string{
		{"name": "year", "typeDefinition": "int", "static": "false"},
		{"name": "date", "typeDefinition": "text", "static": "false"},
		{"name": "start_timestamp", "typeDefinition": "text", "static": "false"},
		{"name": "end_timestamp", "typeDefinition": "text", "static": "false"},
		{"name": "minutes", "typeDefinition": "int", "static": "false"},
		{"name": "category", "typeDefinition": "text", "static": "false"},
		{"name": "sub_category", "typeDefinition": "text", "static": "false"},
	}

	partitionKey := []string{"year"}
	clusteringKey := []string{"date", "start_timestamp"}

	primaryKey := map[string][]string{
		"partitionKey":  partitionKey,
		"clusteringKey": clusteringKey,
	}

	table := gomodoroTable{
		Name:              "gomodoros",
		IfNotExists:       true,
		TableOptions:      map[string]int{"defaultTimeToLive": 0},
		ColumnDefinitions: columnDefinitions,
		PrimaryKey:        primaryKey,
	}

	var body []byte
	body, err := json.Marshal(table)

	if err != nil {
		fmt.Println(err)
	}

	request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(body))
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("x-cassandra-token", token)

	if err != nil {
		log.Fatalf("[ERROR] error building createTable httpRequest: %v", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("[ERROR] error creation Cassandra Gomodoro table: %v", err)
	}
	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[ERROR] error reading Cassandra table creation response: %v", err)
	}

	log.Printf("[INFO] cassandra table synced successfully. DSE response: %v", string(bodyResp))
}

func pushToRemote(token string) {
	localGomodoros := GetGomodorosCassandra()
	postURL := fmt.Sprintf("https://%v-%v.apps.astra.datastax.com/api/rest/v1/keyspaces/%v/tables/%v/rows", config.clusterId, config.clusterRegion, config.keyspace, config.table)

	type gomodoroPush struct {
		Columns []map[string]string `json:"columns"`
	}

	for _, gomodoro := range localGomodoros {

		columns := []map[string]string{
			{"name": "year", "value": gomodoro.Year},
			{"name": "date", "value": gomodoro.Date[0:10]},
			{"name": "start_timestamp", "value": gomodoro.StartTimestamp},
			{"name": "end_timestamp", "value": gomodoro.EndTimestamp},
			{"name": "minutes", "value": strconv.Itoa(gomodoro.Minutes)},
			{"name": "category", "value": gomodoro.Category},
			{"name": "sub_category", "value": gomodoro.SubCategory},
		}

		pushEntry := gomodoroPush{
			Columns: columns,
		}

		body, err := json.Marshal(pushEntry)
		if err != nil {
			log.Fatalf("[ERROR] error bulding json entry for remote push, %v", err.Error())
		}

		request, err := http.NewRequest("POST", postURL, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		request.Header.Set("x-cassandra-token", token)

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalf("[ERROR] error pushing entry to remote, %v", err.Error())
		}
		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("[ERROR] error reading entry push response: %v", err.Error())
		}

	}

	log.Printf("[INFO] local [Go]modoros pushed to remote successfully")
}

func readRemote(token string) {
	years := GetValidYears()

	gToPull := []gomodoro{}

	for _, v := range years {
		remoteGomodoros := getRemoteGomodorosByYear(v, token)
		localGomodoros := GetLocalGomodorosByYear(v)
		for _, v := range remoteGomodoros.Rows {
			remoteEntry := fmt.Sprintf("%v %v", v.Date, v.Start_timestamp)
			found := false

			// TODO: implement better approach
			for _, localEntry := range localGomodoros {
				if localEntry == remoteEntry {
					found = true
					break
				}
			}

			if !found {
				gToPull = append(gToPull, v)
			}
		}
	}

	if len(gToPull) > 0 {
		for _, v := range gToPull {
			log.Printf("[INFO] pulling remote [Go]modoro: date: %v, startTimestamp: %v, category: %v", v.Date, v.Start_timestamp, v.Category)
			InsertRecord(v.Date, v.Start_timestamp, v.End_timestamp, v.Minutes, v.Category, v.Sub_category)
		}
		log.Println("[INFO] [Go]modoro app synced successfully!")
	} else {
		log.Printf("[INFO] [Go]modoro already synchronized, no remote changes were pulled.")
	}

}

func getRemoteGomodorosByYear(year int, token string) gomodoros {

	getURL := fmt.Sprintf("https://%v-%v.apps.astra.datastax.com/api/rest/v1/keyspaces/%v/tables/%v/rows/%v", config.clusterId, config.clusterRegion, config.keyspace, config.table, year)

	request, err := http.NewRequest("GET", getURL, nil)
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("x-cassandra-token", token)

	resp, err := client.Do(request)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	remoteGomodoros := gomodoros{}
	jsonErr := json.Unmarshal(body, &remoteGomodoros)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return remoteGomodoros
}
