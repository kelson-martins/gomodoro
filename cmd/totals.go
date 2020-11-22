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

// totalsCmd represents the totals command
var totalsCmd = &cobra.Command{
	Use:   "totals",
	Short: "Query information about [Go]modoro usage totals",
	Long: `Use the 'total' functionlity to query information about [Go]modoro usage overtime. For example:

gomodoro totals --days 30.`,
	Run: func(cmd *cobra.Command, args []string) {
		db.GetTotalsRecord()
		db.Close()
	},
}

func init() {
	rootCmd.AddCommand(totalsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// totalsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
