// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"log"

	"github.com/gradient-images/teflon/internal/metadata"

	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Gets or sets Teflon metadata",
	Long: `Command 'meta' prints the metadata belonging to the files given as
arguments. As a side effect the metadata access process creates the meta
file if it's not already there or not up to date.`,
	Run: meta,
}

var metaSetCmd = &cobra.Command{
	Use: "set <-e key:value> <target>",
	Short: "Sets a user metadata entry on the given target",
	Long: `Command 'meta set' sets a metadata entry on the the target. If the meta
file doesn't exist 'meta set' will create a new one.`,
  Run: metaSet,
}

var DataList []string

func init() {
	metaSetCmd.Flags().StringSliceVarP(&DataList, "data", "d", []string{}, "Data entry in the form of key:value")
	metaCmd.AddCommand(metaSetCmd)
	rootCmd.AddCommand(metaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func meta(cmd *cobra.Command, args []string) {
	log.Print("'meta' command called")
	for _, baseName := range args {
		us := metadata.Get(baseName)
		us.UserData["valami"] = "azta"
		fmt.Println(us)
	}
}

func metaSet(cmd *cobra.Command, args []string) {
	log.Print("'meta set' command called")
	for _, baseName := range args {
		metadata.Get(baseName)
	}
	fmt.Println(DataList)
}
