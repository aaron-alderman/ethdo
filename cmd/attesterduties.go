// Copyright © 2020 Weald Technology Trading
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

	attesterduties "github.com/aaron-alderman/ethdo/cmd/attester/duties"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var attesterDutiesCmd = &cobra.Command{
	Use:   "duties",
	Short: "Obtain information about duties of an attester",
	Long: `Obtain information about dutes of an attester.  For example:

    ethdo attester duties --account=Validators/00001 --epoch=12345

In quiet mode this will return 0 if a duty from the attester is found, otherwise 1.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := attesterduties.Run(cmd)
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
	attesterCmd.AddCommand(attesterDutiesCmd)
	attesterFlags(attesterDutiesCmd)
	attesterDutiesCmd.Flags().Int64("epoch", -1, "the last complete epoch")
	attesterDutiesCmd.Flags().String("pubkey", "", "the public key of the attester")
	attesterDutiesCmd.Flags().Bool("json", false, "Generate JSON data for an exit; do not broadcast to network")
}

func attesterDutiesBindings() {
	if err := viper.BindPFlag("epoch", attesterDutiesCmd.Flags().Lookup("epoch")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("pubkey", attesterDutiesCmd.Flags().Lookup("pubkey")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("json", attesterDutiesCmd.Flags().Lookup("json")); err != nil {
		panic(err)
	}
}
