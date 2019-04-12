// Copyright © 2019 Máté Birkás <gadfly16@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"

	"github.com/gradient-images/teflon/internal/metadata"

	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Manipulates Teflon metadata",
	Long: `Without any subcommand 'teflon meta' prints the metadata belonging to the targets.
As a side effect the metadata access the command creates the meta
files if they are not already exist or refreshes them if they are not up-to-date.`,
	Run: metaRun,
}

var metaSetCmd = &cobra.Command{
	Use: "set <-d key:value..> <target..>",
	Short: "Sets a user metadata entry on the given target.",
	Long: `Command 'teflon meta set' sets a metadata entry on the the target. If the meta
file doesn't exist 'meta set' will create a new one.`,
  Run: metaSetRun,
}

var DataList []string

func init() {
	metaSetCmd.Flags().StringSliceVarP(&DataList, "data", "d", []string{},
		"Data entry in the form of 'key:value' pairs")
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

func metaRun(cmd *cobra.Command, args []string) {
	log.Print("'meta' command called")
	for _, baseName := range args {
		md := metadata.Get(baseName)
		for k, v := range md.UserSection.UserData {
			fmt.Println(k, ":", v)
		}
	}
}

func metaSetRun(cmd *cobra.Command, args []string) {
	log.Print("'meta set' command called")
	for _, baseName := range args {
		md := metadata.Get(baseName)
		for _, data := range DataList {
			s := strings.SplitN(data, ":", 2)
			key, value := s[0], s[1]
			md.UserSection.UserData[key] = value
		}
		fmt.Println(DataList, md.UserSection.UserData)
		md.Sync()
	}
}
