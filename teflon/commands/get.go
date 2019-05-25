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
	"strings"

	"github.com/gradient-images/teflon"
	"github.com/gradient-images/teflon/expr"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
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
			log.Fatalln("Couldn't init object:", err)
		}

		// Create textual representation for display and unmarshal
		dj, err := json.MarshalIndent(&o, "", "  ")
		if err != nil {
			log.Fatalln("Couldnt marshal meta JSON:", err)
		}

		// Print the whole meta if there was no selected meta
		if len(metaListFlag) == 0 {
			fmt.Println(string(dj))
			return
		}

		// Otherwise find selected meta
		var m map[string]interface{}
		err = json.Unmarshal(dj, &m)
		if err != nil {
			log.Fatalln("Couldn't marshal JSON into map:", err)
		}

		// Cycle through metaList, 'es' stands for expression string
		for _, es := range metaListFlag {
			ei, err := expr.Parse("", []byte(es))
			if err != nil {
				log.Fatalf("Couldn't parse expression: %s, %s", es, err)
			}

			// Now we know that it's a slice of strings
			e := ei.([]string)
			log.Printf("Expr: %s\n", e)

			// // Init vi for the first address, vi stands for value interface
			// vi, ok := m[e[0]]
			// if !ok {
			// 	log.Fatalf("Couldn't find key in meta: %s", e[0])
			// }

			// Create value to descend into and interface object to hold the result
			v := m

			// Value to return
			var val interface{}

			// Display name
			dn := ""

			// Cycle through names in expression
			for i, n := range e {
				// vi, ok := v[n]
				// if !ok {
				// 	log.Fatalf("Couldn't find key in meta 1: %s", n)
				// }
				lm := map[string]string{}
				for k := range v {
					lm[strings.ToLower(k)] = k
				}

				var ok bool
				val, ok = v[lm[strings.ToLower(n)]]
				if !ok {
					log.Fatalf("Couldn't find key in meta 3: %s", n)
				}
				dn = dn + lm[strings.ToLower(n)]

				// If there is more name to come
				if i < len(e)-1 {
					switch val.(type) {
					case map[string]interface{}:
						// Convert interface to map of interfaces
						v = val.(map[string]interface{})
						dn = dn + "."
					default:
						log.Fatalf("Couldn't find key in meta 2: %s", n)
					}
				}
			}
			vs, err := json.MarshalIndent(&val, "", "  ")
			fmt.Printf("%s: %v\n", dn, string(vs))
			if err != nil {
				log.Fatalln("Couldnt marshal result JSON:", err)
			}

		}
	}
}
