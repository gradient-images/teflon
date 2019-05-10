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

var showCmd = &cobra.Command{
	Use:   "show [<targets>]",
	Short: "Shows information about teflon objects",
	Long: `'teflon show' prints various kind of information about Teflon objects. If no
<target> is specified it will run for '.' .`,
	Run: Show,
}

func init() {
	showCmd.Flags().BoolVarP(&protoFlag, "proto", "P", false, "Lists available protos in the context of target.")
	rootCmd.AddCommand(showCmd)
}

// List (`teflon list`) lists various information about objects.
func Show(cmd *cobra.Command, args []string) {
	// Set default target if none is given.
	if len(args) == 0 {
		args = append(args, ".")
	}

	// If '-p' flag is set call listProtos.
	if protoFlag {
		showProtos(cmd, args)
		return
	}
	log.Println("DEBUG: Only proto listing is implemented, doing nothing.")
}

// listProtos (`teflon list -p`) prints the available prototypes in the target's
// context.
func showProtos(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		log.Fatalln("ABORT: Only one target allowed for proto listing.")
	}
	o, err := teflon.NewTeflonObject(args[0])
	if err != nil {
		log.Fatalln("Couldn't create object:", err)
	}
	protoMap, err := o.ListProtos()
	if err != nil {
		log.Fatalln("Couldn't assemble proto list.", err)
	}
	for k, v := range protoMap {
		fmt.Printf("%s: %s\n", k, v)
	}
}
