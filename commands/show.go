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
	"log"
	"os"
	"path/filepath"

	// "strings"

	"github.com/gradient-images/teflon/internal/teflon"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var showCmd = &cobra.Command{
	Use:   "show [target]",
	Short: "Manipulates show objects.",
	Long: `Without any subcommand 'teflon show' prints the metadata of the show the target
belongs to. If target is omitted, it defaults to the current directory.`,
	Run: showRun,
}

var showNewCmd = &cobra.Command{
	Use:   "new <target> [<targets>...]",
	Short: "Creates a new show from a prototype.",
	Long: `Command 'teflon show new' creates a new show at the tartget location based on a
prototype found in the $TEFLON/show_proto directory.`,
	Run: showNewRun,
}

var showProto string

func init() {
	showNewCmd.Flags().StringVarP(&showProto, "proto", "p", "Default", "Show proto to create.")
	showCmd.AddCommand(showNewCmd)
	rootCmd.AddCommand(showCmd)
}

func showRun(cmd *cobra.Command, args []string) {
	log.Print("DEBUG: 'show' command called")
}

func showNewRun(cmd *cobra.Command, args []string) {
	for _, target := range args {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := os.Stat(absTarget); !os.IsNotExist(err) {
			log.Printf("ABORT: '%v' already exists.", target)
			os.Exit(1)
		}

		proto := filepath.Join(teflon.TeflonDir, teflon.ShowProtoDirName, showProto)

		err = copy.Copy(proto, absTarget)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("DONE: Created new show '%v' based on '%v'.", absTarget, showProto)

		o, err := teflon.InitObject(target)
		if err != nil {
			log.Fatalln("Couldn't create object:", err)
		}

		o.Proto = teflon.ShowPrefix + showProto

		if o.SyncMeta() != nil {
			log.Fatalln("Couldn't write meta of newly created show:", err)
		}
	}
}
