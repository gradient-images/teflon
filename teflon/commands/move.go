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

	// "github.com/gradient-images/teflon"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move <source> [<source..>] <destination>",
	Short: "Moves Teflon objects maintaining relations",
	Long: `'teflon move' moves the object specified by the arguments to the location
specified by the last argument. In high contrast to unix mv it also maintains all
Teflon relations like metadata, prototypes, links, etc.`,
	Run: Move,
}

func init() {
	// moveCmd.Flags().StringSliceVarP(&metaListFlag, "meta", "m", []string{},
	// 	"Comma separated list of metadata entry names.")
	rootCmd.AddCommand(moveCmd)
}

func Move(cmd *cobra.Command, args []string) {
	sources := args[:len(args)-1]
	dest := args[len(args)-1]
	log.Printf("DEBUG: sources: %v, dest: %v", sources, dest)
}
