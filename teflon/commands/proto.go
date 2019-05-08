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
	Use:   "proto [<targets>]",
	Short: "Manages prototype information",
	Long: `'teflon proto' prints the protoype information belonging to its targets. If no
<target> is specified it will run for '.' .`,
	Run: Proto,
}

var protoListCmd = &cobra.Command{
	Use:   "list [<target>]",
	Args:  cobra.MaximumNArgs(1),
	Short: "List available prototypes in the given target.",
	Long: `'teflon proto list' prints all the available prototypes at the target's
location. If no <target> is specified it will run for '.' .`,
	Run: ProtoList,
}

func init() {
	protoCmd.AddCommand(protoListCmd)
	rootCmd.AddCommand(protoCmd)
}

// Proto (`teflon proto`) prints the prototype the target belongs to.
func Proto(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		args = append(args, ".")
	}
	for _, target := range args {
		o, err := teflon.NewTeflonObject(target)
		if err != nil {
			log.Fatalln("Couldn't create object:", err)
		}
		if o.Proto != "" {
			fmt.Println(o.Proto)
		}
	}
}

// ProtoList (`teflon proto list`) prints the available prototypes at the target's
// location.
func ProtoList(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		args = append(args, ".")
	}
	o, err := teflon.NewTeflonObject(args[0])
	if err != nil {
		log.Fatalln("Couldn't create object:", err)
	}
	protoMap, err := o.ListProtos()
	if err != nil {
		log.Fatalln("Couldn't assemble proto list.", err)
	}
	for k, v := range protoMap {
		fmt.Printf("%s: %s\n", k, v)
	}
}
