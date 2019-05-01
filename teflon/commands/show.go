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
	"os"
	"path/filepath"

	"github.com/gradient-images/teflon"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var showCmd = &cobra.Command{
	Use:   "show [<target>]",
	Args:  cobra.MaximumNArgs(1),
	Run:   Show,
	Short: "Manipulates show objects.",
	Long: `Without any subcommand 'teflon show' prints the absolute path of the show the target
belongs to. If target is omitted, it defaults to the current directory.`,
}

var showNewCmd = &cobra.Command{
	Use:   "new <target> [<targets>...]",
	Short: "Creates a new show from a prototype.",
	Long: `Command 'teflon show new' creates a new show at the tartget location based on a
prototype found in the $TEFLON/show_proto directory.`,
	Run: ShowNew,
}

var showProto string

func init() {
	showNewCmd.Flags().StringVarP(&showProto, "proto", "p", "Default", "Show proto to create.")
	showCmd.AddCommand(showNewCmd)
	rootCmd.AddCommand(showCmd)
}

// Show (`teflon show`) prints the absolute path to the show the target belongs to.
func Show(cmd *cobra.Command, args []string) {
	log.Print("DEBUG: 'show' command called")
	var target string

	// Check if a target is given
	if len(args) == 0 {
		target = "."
		log.Println("DEBUG: No targets given, running for '.' .")
	} else {
		target = args[0]
	}

	// Create new object for target
	o, err := teflon.NewTeflonObject(target)
	if err != nil {
		log.Fatalln("ABORT: Unable to init object:", target, err)
	}
	log.Printf("DEBUG: Path: %v", o.Path)
	log.Printf("DEBUG: Objects size: %v", len(teflon.Objects))
	if o.Show != nil {
		fmt.Println(o.Show.Path)
	} else {
		fmt.Println("ORPHAN")
	}
}

// ShowNew (`teflon show new`) creates a new show based on a template in
// `teflon.TeflonConf`. The arguments are targets.
func ShowNew(cmd *cobra.Command, targets []string) {
	for _, target := range targets {
		fspath, err := teflon.FSPath(target)
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := os.Stat(fspath); !os.IsNotExist(err) {
			log.Fatalf("ABORT: Target already exists: '%s'", fspath)
		}

		proto := filepath.Join(teflon.TeflonConf, teflon.ShowProtoDirName, showProto)

		err = copy.Copy(proto, fspath)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("SUCCESS: Created new show: %s (%s)", fspath, showProto)

		o, err := teflon.NewTeflonObject(fspath)
		if err != nil {
			log.Fatalln("ABORT: Couldn't create object:", err)
		}

		o.ShowRoot = true
		o.Proto = proto

		if o.SyncMeta() != nil {
			log.Fatalln("ABORT: Couldn't write meta of newly created show:", err)
		}
	}
}
