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

	"github.com/gradient-images/teflon/internal/teflon"

	"github.com/spf13/cobra"
	"github.com/otiai10/copy"
)

// moveCmd represents the move command
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

		if _, err := os.Stat(absTarget); !os.IsNotExist(err) {
			log.Printf("ABORT: '%v' already exists.", target)
			os.Exit(1)
		}

		targetDir, targetName := filepath.Split(absTarget)

		showDir, err := teflon.FindShowRoot(targetDir)
		if err != nil {
			log.Println(err)
		}

		if proto == "" {
			// Infering prototype from file name.
			split := strings.SplitN(targetName, "_", 2)
			proto = split[0]
		}

		protoDir := filepath.Join(showDir, teflon.ProtoDirName)
		proto = filepath.Join(protoDir, proto)

		err = copy.Copy(proto, absTarget)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("DONE: Created '%v' based on '%v'.", target, proto)
	}
}
