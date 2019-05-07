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
	"path/filepath"

	"github.com/gradient-images/teflon"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// `new` creates a new Teflon object from a prototype
var newCmd = &cobra.Command{
	Use:   "new <target..>",
	Short: "Creates a new Teflon object based on a prototype.",
	Long: `Command 'teflon new' looks for a matching prototype upstream the show hierarchy and
creates a new object based on the one it finds matching the requested type.`,
	Run: newRun,
}

var protoFlag string

func init() {
	newCmd.Flags().StringVarP(&protoFlag, "proto", "p", "", "Prototype to create.")
	rootCmd.AddCommand(newCmd)
}

func newRun(cmd *cobra.Command, args []string) {
	for _, target := range args {
		fspath, err := teflon.FSPath(target)
		if err != nil {
			log.Fatalln("ABORT: Malformed target:", err)
		}

		// Create object for parent dir.
		targetDir, targetName := filepath.Split(fspath)
		parent, err := teflon.NewTeflonObject(targetDir)
		if err != nil {
			log.Fatalln("ABORT: Couldn't create object for containing dir:", err)
		}
		if parent.Show == nil {
			log.Fatalln("ABORT: Prototyping not supported outside shows.")
		}

		// Check if explicit prototype is given then find appropriate proto.
		var proto string
		if protoFlag == "" {
			proto, err = parent.FindProtoForTarget(targetName)
			if err != nil {
				log.Fatalln("ABORT: Can't find prototype:", err)
			}
		} else {
			proto, err = parent.FindProto(protoFlag)
			if err != nil {
				log.Fatalln("ABORT: Can't find appropriate prototype:", err)
			}
		}

		// Copy the prototype if target doesn't exist.
		if !teflon.Exist(fspath) {
			err = copy.Copy(proto, fspath)
			if err != nil {
				log.Fatalln("ABORT: Couldn't copy prototype:", err)
			}
			log.Printf("DEBUG: Copied '%s' to '%s'.", proto, fspath)
		} else {
			log.Println("DEBUG: Target exists, skipped copying of proto.")
		}

		// Set Proto: field and clean instances on target.
		o, err := teflon.NewTeflonObject(fspath)
		if err != nil {
			log.Fatalln("ABORT: Couldn't create object:", err)
		}
		err = o.SetProto(proto)
		if err != nil {
			log.Fatalln("ABORT: Couldn't set up proto references:", err)
		}

		log.Printf("SUCCESS: Prototyped '%v' based on '%v'.", target, proto)

	}
}
