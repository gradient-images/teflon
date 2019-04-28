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

var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Manages prototype information",
	Long: `'teflon proto' prints the protoype information belonging to its targets. If no
<target> is specified it will run for '.'. As a side effect the metadata access the
command creates the meta files if they are not already exist or refreshes
them if they are not up-to-date.`,
	Run: protoRun,
}

func init() {
	rootCmd.AddCommand(protoCmd)
}

func protoRun(cmd *cobra.Command, args []string) {
	log.Println("'proto' command called")
	if len(args) == 0 {
		args = append(args, ".")
		log.Println("No targets given, running for '.' .")
	}
	for _, target := range args {
		o, err := teflon.NewInitObject(target)
		if err != nil {
			log.Fatalln("Couldn't create object:", err)
		}
		fmt.Printf("%s: %s\n", target, o.Proto)
	}
}
