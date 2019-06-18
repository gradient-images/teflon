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

var getCmd = &cobra.Command{
	Use:   "get [<expr>]",
	Short: "Reads Teflon metadata",
	Long: `'teflon get' prints the metadata belonging to the targets.  If no <expr>
is specified it will return all the metadata for '.'.`,
	Args: cobra.ExactArgs(1),
	Run:  Get,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

// CLIGet is the command line wrapper around the Teflon API's Get function.
func Get(cmd *cobra.Command, args []string) {
	// Create object for current working directory
	pwd, err := teflon.NewTeflonObject(".")
	if err != nil {
		log.Fatalln("Couldn't create object for '.' :", err)
	}

	// Run Get.
	res, err := pwd.Get(args[0])
	if err != nil {
		log.Fatalln("ABORT: Couldn't get results:", err)
	}

	close(teflon.Events)

	// Create display string of result (dres).
	dres, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatalln("ABORT: Couldnt marshal result JSON:", err)
	}

	// Print display result to terminal.
	fmt.Printf("%s: %s\n", args[0], dres)

}
