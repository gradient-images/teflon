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
	"os"
	"path/filepath"
	"strings"

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

var proto string

func init() {
	newCmd.Flags().StringVarP(&proto, "proto", "p", "", "Prototype to create.")
	rootCmd.AddCommand(newCmd)
}

func newRun(cmd *cobra.Command, args []string) {
	for _, target := range args {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			log.Fatalln(err)
		}

		// Check if target exists
		if _, err := os.Stat(absTarget); !os.IsNotExist(err) {
			log.Printf("ABORT: '%v' already exists.", target)
			os.Exit(1)
		}

		targetDir, targetName := filepath.Split(absTarget)

		pdl := teflon.FindProtoDirs(targetDir)
		log.Println("Proto dir list:", pdl)

		if proto == "" {
			// Infering prototype from file name.
			split := strings.SplitN(targetName, "_", 2)
			proto = split[0]
		}

		var p string
		for _, pd := range pdl {
			p = filepath.Join(pd, proto)
			if _, err := os.Stat(p); os.IsNotExist(err) {
				p = ""
				continue
			}
			break
		}

		if p == "" {
			log.Fatalln("No prototype found:", proto)
		}

		err = copy.Copy(p, absTarget)
		if err != nil {
			log.Fatalln(err)
		}

		// Set Proto: field on target
		o, err := teflon.NewInitObject(absTarget)
		if err != nil {
			log.Fatalln("Couldn't create object:", err)
		}

		o.Proto = p
		o.Instances = []string{}

		if o.SyncMeta() != nil {
			log.Fatalln("Couldn't write meta of newly created show:", err)
		}

		// Set Instances: field on proto
		o, err = teflon.NewInitObject(p)
		if err != nil {
			log.Fatalln("Couldn't create object:", err)
		}

		o.Instances = append(o.Instances, absTarget)
		log.Println("Instances: ", o.Instances)

		if o.SyncMeta() != nil {
			log.Fatalln("Couldn't write meta of newly created show:", err)
		}
		log.Printf("DONE: Created '%v' based on '%v'.", target, p)

	}
}
