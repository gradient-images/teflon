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
	// "fmt"
	"os"
	"os/exec"
	"log"

	// "github.com/gradient-images/teflon"
	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:   "conf [<targets>]",
	Short: "Manages the teflon process' configuration",
	Long: `'teflon conf' manages various aspects of configuration of the teflon process.`,
	Run: Conf,
}

var (
	confInitFlag bool
	confForceFlag bool
	confDirFlag string
	confRepoFlag string
)

func init() {
	confCmd.Flags().BoolVarP(&confInitFlag, "init", "I", false, "Initializes a config dir.")
	confCmd.Flags().StringVarP(&confDirFlag, "conf-dir", "c", ".teflonconf", "Configuration directory to work on.")
	confCmd.Flags().StringVarP(&confRepoFlag, "repository", "r",
			"https://github.com/gradient-images/teflon-reference-config.git",
			"Configuration directory to work on.",
		)
	rootCmd.AddCommand(confCmd)
}

// List (`teflon list`) lists various information about objects.
func Conf(cmd *cobra.Command, args []string) {
	if confInitFlag {
		InitConf(cmd, args)
		os.Exit(0)
	}
}

func InitConf(cmd *cobra.Command, args []string) {
	log.Println("DEBUG: Initializing configuration directory.")
	task := exec.Command("git", "clone", confRepoFlag, confDirFlag)
	err := task.Run()
	if err != nil {
		log.Fatalln("Git returned with error code.", err)
	}

}
