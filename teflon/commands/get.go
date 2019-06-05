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
	"os"

	tfl "github.com/gradient-images/teflon"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [<expr>]",
	Short: "Reads Teflon metadata",
	Long: `'teflon get' prints the metadata belonging to the targets.  If no <expr>
is specified it will return all the metadata for '.'.`,
	Args: cobra.MaximumNArgs(1),
	Run:  CliGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

// CLIGet is the command line wrapper around the Teflon API's Get function.
func CliGet(cmd *cobra.Command, args []string) {
	// Init expression string.
	ex := ""
	if len(args) == 1 {
		ex = args[0]
	}

	// Get current working directory.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("FATAL: Couldn't get current working dir:", err)
	}

	// Run Get.
	res, err := tfl.Get(dir, ex)
	if err != nil {
		log.Fatalln("FATAL: Couldn't evaluate expression:", err)
	}

	// Create display string of result (dres).
	dres, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatalln("FATAL: Couldnt marshal result JSON:", err)
	}

	// Print display result to terminal.
	fmt.Printf("%s: %s\n", ex, dres)

	// // Set default argument if no arguments given.
	// if len(args) == 0 {
	// 	args = append(args, ".")
	// }
	//
	// // Cycle through all args as Targets.
	// for _, target := range args {
	// 	o, err := teflon.NewTeflonObject(target)
	// 	if err != nil {
	// 		log.Fatalln("FATAL: Couldn't init object:", err)
	// 	}
	//
	// 	// Print the whole meta if there was no selected meta
	// 	if len(metaListFlag) == 0 {
	// 		dj, err := json.MarshalIndent(&o, "", "  ")
	// 		if err != nil {
	// 			log.Fatalln("FATAL: Couldn't marshal entire object to JSON:", err)
	// 		}
	//
	// 		fmt.Println(string(dj))
	// 		return
	// 	}
	//
	// 	// Get the context for ident resolution.
	// 	c, err := o.GetContext()
	// 	if err != nil {
	// 		log.Fatalln("FATAL: Couldn't create context:", err)
	// 	}
	//
	// 	// Cycle through metaList, 'es' stands for expression string
	// 	for _, es := range metaListFlag {
	// 		log.Println("DEBUG: es:", es)
	// 		// Create and parse expression
	// 		e, err := teflon.NewExpr(es)
	// 		if err != nil {
	// 			log.Fatalf("FATAL: Couldn't create expression: %s, %s", es, err)
	// 		}
	// 		log.Printf("DEBUG: Expr: %s\n", e)
	//
	// 		// Evaluate expression.
	// 		v, err := e.Eval(c)
	// 		if err != nil {
	// 			log.Fatalln("FATAL: Couldn't evaluate expression:", err)
	// 		}
	//
	// 		// Create display JSON.
	// 		dj, err := json.MarshalIndent(v, "", "  ")
	//
	// 		fmt.Printf("%s: %s\n", e, dj)
	// 		if err != nil {
	// 			log.Fatalln("FATAL: Couldnt marshal result JSON:", err)
	// 		}
	// 	}
	// }
}
