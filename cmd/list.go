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
	"github.com/kelson-martins/gomodoro/src/db"
	"github.com/spf13/cobra"
)

var days int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List [Go]modoro sessions",
	Long: `The 'list' action is used to list gomodoro sessions. Example:

gomodoro list
gomodoro list --days 10
`,
	Run: func(cmd *cobra.Command, args []string) {
		db.Init()
		defer db.Close()

		db.ListGomodoros(days)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().IntVarP(&days, "days", "d", 30, "Days to list [Go]modoro sessions")

}
