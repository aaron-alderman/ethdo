// Copyright © 2019, 2020 Weald Technology Trading
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

	accountderive "github.com/aaron-alderman/ethdo/cmd/account/derive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountDeriveCmd = &cobra.Command{
	Use:   "derive",
	Short: "Derive an account",
	Long: `Derive an account from a mnemonic and path.  For example:

    ethdo account derive --mnemonic="..." --path="m/12381/3600/0/0"

In quiet mode this will return 0 if the inputs can derive an account account, otherwise 1.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := accountderive.Run(cmd)
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
	accountCmd.AddCommand(accountDeriveCmd)
	accountFlags(accountDeriveCmd)
	accountDeriveCmd.Flags().String("mnemonic", "", "mnemonic from which to derive the HD seed")
	accountDeriveCmd.Flags().String("path", "", "path from which to derive the account")
	accountDeriveCmd.Flags().Bool("show-private-key", false, "show private key for derived account")
	accountDeriveCmd.Flags().Bool("show-withdrawal-credentials", false, "show withdrawal credentials for derived account")
}

func accountDeriveBindings() {
	if err := viper.BindPFlag("mnemonic", accountDeriveCmd.Flags().Lookup("mnemonic")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("path", accountDeriveCmd.Flags().Lookup("path")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("show-private-key", accountDeriveCmd.Flags().Lookup("show-private-key")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("show-withdrawal-credentials", accountDeriveCmd.Flags().Lookup("show-withdrawal-credentials")); err != nil {
		panic(err)
	}
}
