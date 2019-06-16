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

	"github.com/gradient-images/teflon"
	"github.com/spf13/cobra"
)

// `new` creates a new Teflon object from a prototype
var newCmd = &cobra.Command{
	Use:   "new <target..>",
	Short: "Creates a new Teflon object based on a prototype",
	Long: `Command 'teflon new' looks for a matching prototype upstream the show hierarchy and
creates a new object based on the one it finds matching the requested type. If the
'-S' flag is present, the command will create a new show from a show prototype.`,
	Args: cobra.ExactArgs(1),
	Run:  New,
}

var newFileFlag bool
var newShowProtoFlag string

func init() {
	newCmd.Flags().BoolVarP(
		&newFileFlag,
		"file",
		"f",
		false,
		"Create an empty file instead of an empty directory.",
	)
	newCmd.Flags().BoolVarP(
		&showFlag,
		"show",
		"S",
		false,
		"Create new show.",
	)
	newCmd.Flags().StringVarP(
		&newShowProtoFlag,
		"show-proto",
		"p",
		"Default",
		"Prototype to use during show creation.",
	)
	rootCmd.AddCommand(newCmd)
}

// Creates new teflon object and triggers new event.
func New(cmd *cobra.Command, args []string) {
	// Create object for current working directory
	pwd, err := teflon.NewTeflonObject(".")
	if err != nil {
		log.Fatalln("Couldn't create object for '.' :", err)
	}

	// If showFlag is set `new` will create a shows instead of a regular object.
	if showFlag {
		nshws, err := pwd.CreateShow(args[0], newShowProtoFlag)
		if err != nil {
			log.Fatalln("ABORT: Couldnt create show:", err)
		}
		for _, shw := range nshws {
			fmt.Println(shw.Path)
		}
	}

	// // Create regular objects.
	// for _, target := range args {
	// 	fspath, err := teflon.Path(target)
	// 	if err != nil {
	// 		log.Fatalln("ABORT: Malformed target:", err)
	// 	}
	//
	// 	// Create object for parent dir.
	// 	targetDir, targetName := filepath.Split(fspath)
	// 	parent, err := teflon.NewTeflonObject(targetDir)
	// 	if err != nil {
	// 		log.Fatalln("ABORT: Couldn't create object for containing dir:", err)
	// 	}
	// 	if parent.Show == nil {
	// 		log.Fatalln("ABORT: Prototyping not supported outside shows.")
	// 	}
	//
	// 	// Check if explicit prototype is given then find appropriate proto.
	// 	var proto string
	// 	if newForceProtoFlag == "" {
	// 		proto, err = parent.FindProtoForTarget(targetName)
	// 		if err != nil {
	// 			log.Fatalln("ABORT: Can't find prototype:", err)
	// 		}
	// 	} else {
	// 		proto, err = parent.FindProto(newForceProtoFlag)
	// 		if err != nil {
	// 			log.Fatalln("ABORT: Can't find appropriate prototype:", err)
	// 		}
	// 	}
	//
	// 	// Copy the prototype if target doesn't exist.
	// 	if !teflon.Exist(fspath) {
	// 		err = copy.Copy(proto, fspath)
	// 		if err != nil {
	// 			log.Fatalln("ABORT: Couldn't copy prototype:", err)
	// 		}
	// 		log.Printf("DEBUG: Copied '%s' to '%s'.", teflon.ProtoPath(proto), teflon.ShowAbs(fspath))
	// 	} else {
	// 		log.Println("DEBUG: Target exists, skipped copying of proto.")
	// 	}
	//
	// 	// Set Proto: field and clean instances on target.
	// 	o, err := teflon.NewTeflonObject(fspath)
	// 	if err != nil {
	// 		log.Fatalln("ABORT: Couldn't create object:", err)
	// 	}
	// 	err = o.SetProto(proto)
	// 	if err != nil {
	// 		log.Fatalln("ABORT: Couldn't set up proto references:", err)
	// 	}
	//
	// 	log.Printf("SUCCESS: Prototyped '%v' based on '%v'.", target, teflon.ProtoPath(proto))
	//
	// }
}

// newShow (`teflon new -s`) creates a new show based on a template in
// `teflon.TeflonConf`. The arguments are targets.
// func newShow(cmd *cobra.Command, targets []string) {
// 	for _, target := range targets {
// 		fspath, err := teflon.Path(target)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
//
// 		if _, err := os.Stat(fspath); !os.IsNotExist(err) {
// 			log.Fatalf("ABORT: Target already exists: '%s'", fspath)
// 		}
//
// 		// Set default proto if not set.
// 		if newForceProtoFlag == "" {
// 			newForceProtoFlag = "Default"
// 		}
//
// 		proto := filepath.Join(teflon.TeflonConf, teflon.ShowProtoDirName, newForceProtoFlag)
//
// 		err = copy.Copy(proto, fspath)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		log.Printf("SUCCESS: Created new show: %s (%s)", fspath, newForceProtoFlag)
//
// 		o, err := teflon.NewTeflonObject(fspath)
// 		if err != nil {
// 			log.Fatalln("ABORT: Couldn't create object:", err)
// 		}
//
// 		o.ShowRoot = true
// 		o.Proto = proto
//
// 		if o.SyncMeta() != nil {
// 			log.Fatalln("ABORT: Couldn't write meta of newly created show:", err)
// 		}
// 	}
// }
