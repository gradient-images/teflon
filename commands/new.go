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

var newShowCmd = &cobra.Command{
	Use: "show <proto> <target..>",
	Short: "Creates a new show.",
	Long: `Command 'teflon new show' creates a new show at the tartget location based on a
prototype found in the $TEFLON/show_proto directory.`,
  Run: newShowRun,
}

var showProto string

func init() {
	// metaSetCmd.Flags().StringSliceVarP(&DataList, "data", "d", []string{},
	// 	"Data entry in the form of 'key:value' pairs")
	newShowCmd.Flags().StringVarP(&showProto, "proto", "p", "Default", "Show proto to create.")
	newCmd.AddCommand(newShowCmd)
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func newRun(cmd *cobra.Command, args []string) {
	log.Print("'new' command called")
	for _, target := range args {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			log.Fatalln(err)
		}

		targetDir, targetName := filepath.Split(absTarget)

		showDir, err := teflon.FindShowRoot(targetDir)
		if err != nil {
			log.Println(err)
		}

		split := strings.SplitN(targetName, "_", 2)
		protoName, identity := split[0], ""
		if len(split) == 2 {
			protoName, identity = split[0], split[1]
		}

		protoDir := filepath.Join(showDir, teflon.ProtoDirName)
		proto := filepath.Join(protoDir, protoName)

		fmt.Println("protoName:", protoName, "Identitiy:", identity, "Proto:", proto)
		err = copy.Copy(proto, absTarget)
		if err != nil {
			log.Fatalln(err)
		}
		// for k, v := range md.UserSection.UserData {
		// 	fmt.Println(k, ":", v)
		// }
	}
}

func newShowRun(cmd *cobra.Command, args []string) {
	log.Print("'new show' command called")
	for _, target := range args {
		absTarget, err := filepath.Abs(target)
		if err != nil {
			log.Fatalln(err)
		}

		targetDir, targetName := filepath.Split(absTarget)

		teflonDir, err := teflon.FindTeflonRoot(targetDir)
		if err != nil {
			log.Println(err)
		}

		protoDir := filepath.Join(teflonDir, teflon.TeflonDirName, teflon.ShowProtoDirName)
		proto := filepath.Join(protoDir, showProto)

		fmt.Println("protoName:", showProto, "Identitiy:", targetName, "Proto:", proto)

		err = copy.Copy(proto, absTarget)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
