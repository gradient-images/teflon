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
	"encoding/json"
	"fmt"
	"log"

	"github.com/gradient-images/teflon"

	"github.com/spf13/cobra"
)

var contractPatternFlag string

var contractCmd = &cobra.Command{
	Use:   "contract [-p <pattern>] [<expr>]",
	Short: "Manipulates contracts",
	Long:  `'teflon contract' sets and gets contract related system metadata on contract objects.`,
	Args:  cobra.ExactArgs(1),
	Run:   Contract,
}

func init() {
	contractCmd.Flags().StringVarP(
		&contractPatternFlag,
		"pattern",
		"p",
		"",
		"Set the contracts Pattern value to the given object selector expression.")
	rootCmd.AddCommand(contractCmd)
}

// CLIGet is the command line wrapper around the Teflon API's Get function.
func Contract(cmd *cobra.Command, args []string) {
	var res interface{}
	var err error

	// Create object for current working directory
	pwd, err := teflon.NewTeflonObject(".")
	if err != nil {
		log.Fatalln("Couldn't create object for '.' :", err)
	}

	if contractPatternFlag != "" {
		log.Println("DEBUG: Setting pattern for contracts.")
		// Run Get.
		res, err = pwd.SetContractPattern(args[0], contractPatternFlag)
		if err != nil {
			log.Fatalln("ABORT: Couldn't set contract:", err)
		}
	}

	// Create display string of result (dres).
	dres, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatalln("ABORT: Couldnt marshal result JSON:", err)
	}

	// Print display result to terminal.
	fmt.Printf("%s: %s\n", args[0], dres)
}
