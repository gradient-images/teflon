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

package cmd

import (
	"fmt"
	"os"
	"log"
	"path/filepath"

	"github.com/gradient-images/teflon/internal/metadata"

	"github.com/spf13/cobra"
)

const (
	teflonDirName = ".teflon"
	teflonDirMetaName = "_"
	teflonMetaExt = "._"
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

func init() {
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
	log.Print("meta called")
	for _, baseName := range args {
		baseInfo, err := os.Stat(baseName)
		if err != nil {
			log.Fatal(err)
		}
		var metaName string
		if baseInfo.IsDir() {
			metaName = filepath.Join(baseName, teflonDirName, teflonDirMetaName)
		} else {
			d, n := filepath.Split(baseName)
			metaName = filepath.Join(d, teflonDirName, n + teflonMetaExt)
		}

		if _, err := os.Stat(metaName); os.IsNotExist(err) {
			log.Print("Meta file doesn't exists.")
		} else {
			log.Print("Meta file exists.")
		}
		fmt.Println(baseName, metaName)
		us := metadata.UserSection{}
		us.UserData = make(map[string]string)
		us.UserData["valami"] = "azta"
		fmt.Println(us)
	}
}
