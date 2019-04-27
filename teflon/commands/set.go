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
	"strings"

	"github.com/gradient-images/teflon"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <-d key:value..> <target..>",
	Short: "Sets a user metadata entry on the given target",
	Long: `Command 'teflon meta set' sets a metadata entry on the the target. If no <target>
is specified it will run for '.'. If the meta file doesn't exist 'meta set'
will create a new one. If only a key is given to the -d flag, the entry for
the key will be deleted.`,
	Run: setRun,
}

func init() {
	setCmd.Flags().StringSliceVarP(&DataList, "data", "d", []string{},
		"Data entry in the form of 'key:value' pairs")
	rootCmd.AddCommand(setCmd)
}

func setRun(cmd *cobra.Command, args []string) {
	log.Print("'set' command called")
	if len(args) == 0 {
		args = append(args, ".")
		log.Println("No targets given, running for '.' .")
	}
	for _, target := range args {
		o, err := teflon.InitObject(target)
		if err != nil {
			log.Fatalln("Couldn't create object:", err)
		}
		for _, data := range DataList {
			s := strings.SplitN(data, ":", 2)
			if len(s) < 2 {
				log.Fatalln("FALTAL: Malformed metadata:", data)
			}
			if s[1] == "" {
				log.Println("Deleting metadata entry:", s[0])
				o.DelMeta(s[0])
			} else {
				log.Println("Setting metadata entry:", s[0]+":", s[1])
				o.SetMeta(s[0], s[1])
			}
		}
		o.SyncMeta()
	}
}
