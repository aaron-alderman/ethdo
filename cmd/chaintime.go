// Copyright © 2021 Weald Technology Trading
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

	chaintime "github.com/aaron-alderman/ethdo/cmd/chain/time"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chainTimeCmd = &cobra.Command{
	Use:   "time",
	Short: "Obtain info about the chain at a given time",
	Long: `Obtain info about the chain at a given time.  For example:

    ethdo chain time --slot=12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := chaintime.Run(cmd)
		if err != nil {
			return err
		}
		if res != "" {
			fmt.Print(res)
		}
		return nil
	},
}

func init() {
	chainCmd.AddCommand(chainTimeCmd)
	chainFlags(chainTimeCmd)
	chainTimeCmd.Flags().String("slot", "", "The slot for which to obtain information")
	chainTimeCmd.Flags().String("epoch", "", "The epoch for which to obtain information")
	chainTimeCmd.Flags().String("timestamp", "", "The timestamp for which to obtain information (format YYYY-MM-DDTHH:MM:SS+ZZZZ)")
}

func chainTimeBindings() {
	if err := viper.BindPFlag("slot", chainTimeCmd.Flags().Lookup("slot")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("epoch", chainTimeCmd.Flags().Lookup("epoch")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("timestamp", chainTimeCmd.Flags().Lookup("timestamp")); err != nil {
		panic(err)
	}
}
