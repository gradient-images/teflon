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
	Use:   "get -m [<expr>..] [<target>..]",
	Short: "Reads Teflon metadata",
	Long: `'teflon get' prints the metadata belonging to the targets.  If no <target>
is specified it will run for '.'.`,
	Run: Get,
}

func init() {
	getCmd.Flags().StringSliceVarP(&metaListFlag, "meta", "m", []string{},
		"Comma separated list of metadata entry names.")
	rootCmd.AddCommand(getCmd)
}

func Get(cmd *cobra.Command, args []string) {
	// Set default argument if no arguments given.
	if len(args) == 0 {
		args = append(args, ".")
	}

	// Cycle through all args as Targets.
	for _, target := range args {
		o, err := teflon.NewTeflonObject(target)
		if err != nil {
			log.Fatalln("FATAL: Couldn't init object:", err)
		}

		// Print the whole meta if there was no selected meta
		if len(metaListFlag) == 0 {
			dj, err := json.MarshalIndent(&o, "", "  ")
			if err != nil {
				log.Fatalln("FATAL: Couldn't marshal entire object to JSON:", err)
			}

			fmt.Println(string(dj))
			return
		}

		// Get the context for ident resolution.
		c, err := o.GetContext()
		if err != nil {
			log.Fatalln("FATAL: Couldn't create context:", err)
		}

		// Cycle through metaList, 'es' stands for expression string
		for _, es := range metaListFlag {
			log.Println("DEBUG: es:", es)
			// Create and parse expression
			e, err := teflon.NewExpr(es)
			if err != nil {
				log.Fatalf("FATAL: Couldn't create expression: %s, %s", es, err)
			}
			log.Printf("DEBUG: Expr: %s\n", e)

			// Evaluate expression.
			v, err := e.Eval(c)
			if err != nil {
				log.Fatalln("FATAL: Couldn't evaluate expression:", err)
			}

			// Create display JSON.
			dj, err := json.MarshalIndent(v, "", "  ")

			fmt.Printf("%s: %s\n", e, dj)
			if err != nil {
				log.Fatalln("FATAL: Couldnt marshal result JSON:", err)

				// // Init context
				// v := m
				//
				// // Value to return and display name
				// var val interface{}
				// var dn string
				//
				// // Cycle through names in Identifier
				// for i, n := range e {
				// 	// Create lower map for case insensitive matching
				// 	lm := map[string]string{}
				// 	for k := range v {
				// 		lm[strings.ToLower(k)] = k
				// 	}
				//
				// 	val, ok := v[lm[strings.ToLower(n)]]
				// 	if !ok {
				// 		log.Fatalf("Couldn't find key in meta 3: %s", n)
				// 	}
				// 	dn = dn + lm[strings.ToLower(n)]
				//
				// 	// If there is more name to come
				// 	if i < len(e)-1 {
				// 		switch val.(type) {
				// 		case map[string]interface{}:
				// 			// Convert interface to map of interfaces
				// 			v = val.(map[string]interface{})
				// 			dn = dn + "."
				// 		default:
				// 			log.Fatalf("Couldn't find key in meta 2: %s", n)
				// 		}
				// 	}
				// }
				// vs, err := json.MarshalIndent(&val, "", "  ")
			}

		}
	}
}
