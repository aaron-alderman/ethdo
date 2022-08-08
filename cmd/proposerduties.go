// Copyright © 2022 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	proposerduties "github.com/aaron-alderman/ethdo/cmd/proposer/duties"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var proposerDutiesCmd = &cobra.Command{
	Use:   "duties",
	Short: "Obtain information about duties of an proposer",
	Long: `Obtain information about dutes of an proposer.  For example:

    ethdo proposer duties --epoch=12345

In quiet mode this will return 0 if duties can be obtained, otherwise 1.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := proposerduties.Run(cmd)
		if err != nil {
			return err
		}
		if viper.GetBool("quiet") {
			return nil
		}
		if res != "" {
			fmt.Println(res)
		}
		return nil
	},
}

func init() {
	proposerCmd.AddCommand(proposerDutiesCmd)
	proposerFlags(proposerDutiesCmd)
	proposerDutiesCmd.Flags().String("epoch", "", "the epoch for which to fetch duties")
	proposerDutiesCmd.Flags().Bool("json", false, "output data in JSON format")
}

func proposerDutiesBindings() {
	if err := viper.BindPFlag("epoch", proposerDutiesCmd.Flags().Lookup("epoch")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("json", proposerDutiesCmd.Flags().Lookup("json")); err != nil {
		panic(err)
	}
}
