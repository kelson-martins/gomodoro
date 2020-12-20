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
	"context"
	"db"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var backup_dir = "gomodoro_backup"

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the [Go]modoro database into Google Drive",
	Long: `The backup feature integrates the application into your Google Drive to backup the [Go]modoro SQLite database. 
The database is saved in the 'gomodoro_backup' folder at the root of your Google Drive.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[INFO] initiating Gomodoro database backup")
		db.LoadConfig()
		dbPath := db.GetDbPath()

		// Step 1. Open the file
		f, err := os.Open(dbPath)

		if err != nil {
			panic(fmt.Sprintf("[ERROR] cannot open gomodoro database file: %v", err))
		}

		defer f.Close()

		// Step 2. Get the Google Drive service
		service, err := getService()

		// Step 3. Create the directory
		dir, err := createDir(service, backup_dir, "root")

		// Step 4. Create the file and upload its content
		createFile(service, "gomodoro.db", "Unknown", f, dir)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	db.LoadConfig()

	tokFile := fmt.Sprint(getUserHome(), "/gomodoro/token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("[INFO] go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("[ERROR] unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func createDir(service *drive.Service, name string, parentID string) (*drive.File, error) {

	fileID := ""

	r, _ := service.Files.List().Do()
	for _, i := range r.Files {
		if i.Name == backup_dir {
			fmt.Printf("[INFO] backup folder %v found, re-using it.\n", i.Name)
			fileID = i.Id
			break
		}
	}

	// create directory if not existing
	if fileID == "" {
		fmt.Println("[INFO] creating backup directory")
		d := &drive.File{
			Name:     name,
			MimeType: "application/vnd.google-apps.folder",
			Parents:  []string{parentID},
		}

		file, err := service.Files.Create(d).Do()

		if err != nil {
			panic(fmt.Sprintf("[ERROR] could not create dir: %v\n", err))
		}

		return file, nil
	}

	file, err := service.Files.Get(fileID).Do()

	if err != nil {
		panic(fmt.Sprintf("[ERROR] error getting backup directory: %v\n", err))
	}

	return file, nil

}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parent *drive.File) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parent.Id},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("[ERROR] could not create file: " + err.Error())
	} else {
		fmt.Printf("[INFO] [Go]modoro database '%s' successfully saved at '%s' directory\n", file.Name, parent.Name)
	}

}

func getService() (*drive.Service, error) {

	driveFile := fmt.Sprint(getUserHome(), "/gomodoro/credentials.json")

	b, err := ioutil.ReadFile(driveFile)
	if err != nil {
		fmt.Printf("[ERROR] unable to read credentials.json file. Err: %v\n", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)

	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, err := drive.New(client)

	if err != nil {
		fmt.Printf("[EROR] cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

func getUserHome() string {

	user, err := user.Current()

	if err != nil {
		log.Fatalln("[ERROR] failure to load user's HOME dir")
		os.Exit(1)
	}

	return user.HomeDir
}
