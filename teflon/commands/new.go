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
	"os"

	"github.com/gradient-images/teflon"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// `new` creates a new Teflon object from a prototype
var newCmd = &cobra.Command{
	Use:   "new <target..>",
	Short: "Creates a new Teflon object based on a prototype",
	Long: `Command 'teflon new' looks for a matching prototype upstream the show hierarchy and
creates a new object based on the one it finds matching the requested type. If the
'-S' flag is present, the command will create a new show from a show prototype.`,
	Run: New,
}

var forceProtoFlag string

func init() {
	newCmd.Flags().StringVarP(&forceProtoFlag, "force-proto", "p", "", "Prototype to use. If '-S' is set, it defaults to 'Default'.")
	newCmd.Flags().BoolVarP(&showFlag, "show", "S", false, "Create a show.")
	rootCmd.AddCommand(newCmd)
}

// Creates a new teflon object from a prototype.
func New(cmd *cobra.Command, args []string) {
	// If show flag is set create shows instead of regular objects.
	if showFlag {
		newShow(cmd, args)
		return
	}

	// Create regular objects.
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
		if forceProtoFlag == "" {
			proto, err = parent.FindProtoForTarget(targetName)
			if err != nil {
				log.Fatalln("ABORT: Can't find prototype:", err)
			}
		} else {
			proto, err = parent.FindProto(forceProtoFlag)
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
			log.Printf("DEBUG: Copied '%s' to '%s'.", teflon.ProtoPath(proto), teflon.ShowAbs(fspath))
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

		log.Printf("SUCCESS: Prototyped '%v' based on '%v'.", target, teflon.ProtoPath(proto))

	}
}

// newShow (`teflon new -s`) creates a new show based on a template in
// `teflon.TeflonConf`. The arguments are targets.
func newShow(cmd *cobra.Command, targets []string) {
	for _, target := range targets {
		fspath, err := teflon.FSPath(target)
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := os.Stat(fspath); !os.IsNotExist(err) {
			log.Fatalf("ABORT: Target already exists: '%s'", fspath)
		}

		// Set default proto if not set.
		if forceProtoFlag == "" {
			forceProtoFlag = "Default"
		}

		proto := filepath.Join(teflon.TeflonConf, teflon.ShowProtoDirName, forceProtoFlag)

		err = copy.Copy(proto, fspath)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("SUCCESS: Created new show: %s (%s)", fspath, forceProtoFlag)

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
