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
	"path/filepath"

	"github.com/gradient-images/teflon"
	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:   "conf [<conf_dir>]",
	Short: "Manages the configuration of a teflon process.",
	Long: `'teflon conf' manages various aspects of configuration of the teflon process.
If no arguments are given it will default to '.teflonconf'`,
	Run: Conf,
}

var (
	confInitFlag bool
	confForceFlag bool
	confRepoFlag string
	confEmptyFlag bool
)

func init() {
	confCmd.Flags().BoolVarP(&confInitFlag, "init", "I", false, "Initializes a config dir.")
	confCmd.Flags().BoolVarP(&confEmptyFlag, "empty", "E", false, "Makes '-I' to create an empty config dir.")
	confCmd.Flags().StringVarP(&confRepoFlag, "repository", "r",
			"https://github.com/gradient-images/teflon-reference-config.git",
			"Configuration directory to work on.",
		)
	rootCmd.AddCommand(confCmd)
}

// List (`teflon conf`) manipulates the config dir.
func Conf(cmd *cobra.Command, args []string) {
	// Set default argument if no arguments given.
	if len(args) == 0 {
		args = append(args, ".teflonconf")
	}

	// Only one config dir can be manipulated right now.
	if len(args) > 1 {
		log.Fatalln("ABORT: Only one conf dir path can be given. Ask for more if seems useful! :)")
	}

	// Initializing a new config dir.
	if confInitFlag {
		if confEmptyFlag {
			err := os.MkdirAll(filepath.Join(args[0], teflon.ShowProtoDirName, "Default"), os.FileMode(0755))
			if err != nil {
				log.Fatalln("ABORT: Couldn't create empty config dir.")
			}
			log.Println("SUCCESS: Created empty conf dir:", args[0])
			os.Exit(0)
		}
		InitConf(cmd, args)
		os.Exit(0)
	}
}

func InitConf(cmd *cobra.Command, args []string) {
	log.Println("DEBUG: Initializing configuration directory.")
	task := exec.Command("git", "clone", confRepoFlag, args[0])
	err := task.Run()
	if err != nil {
		log.Fatalln("Git returned with error code.", err)
	}

}
