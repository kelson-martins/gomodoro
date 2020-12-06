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
	"db"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var id int

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a [Go]modoro from the database",
	Long: `The 'delete' command is used for deleting a particular [Go]modoro instance. For example:

gomodoro delete --id 2
gomodoro delete --id latest`,
	Run: func(cmd *cobra.Command, args []string) {
		db.Init()
		defer db.Close()

		if id == 0 {
			db.DeleteLatest()
		} else {
			db.ExternalDeleteID(id)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deleteCmd.Flags().IntVarP(&id, "id", "i", 0, "ID of the [Go]modoro instance to delete. 0 for latest.")
	deleteCmd.MarkFlagRequired("id")
}
