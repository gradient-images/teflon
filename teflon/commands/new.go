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

var newFileFlag bool
var newShowProtoFlag string

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
			log.Fatalln("ABORT: Couldnt create shows:", err)
		}
		for _, shw := range nshws {
			fmt.Println(shw.Path)
		}
		return
	}

	nobjs, err := pwd.CreateObject(args[0], newFileFlag)
	if err != nil {
		log.Fatalln("ABORT: Couldn't create objects:", err)
	}

	for _, obj := range nobjs {
		fmt.Println(obj.Path)
	}
}
