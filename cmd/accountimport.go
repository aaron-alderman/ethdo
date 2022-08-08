// Copyright © 2019 Weald Technology Trading
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

	accountimport "github.com/aaron-alderman/ethdo/cmd/account/import"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accountImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import an account",
	Long: `Import an account from its private key.  For example:

    ethdo account import --account="primary/testing" --key="0x..." --passphrase="my secret"

In quiet mode this will return 0 if the account is imported successfully, otherwise 1.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := accountimport.Run(cmd)
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
	accountCmd.AddCommand(accountImportCmd)
	accountFlags(accountImportCmd)
	accountImportCmd.Flags().String("key", "", "Private key of the account to import (0x...)")
	accountImportCmd.Flags().String("keystore", "", "Keystore, or path to keystore ")
	accountImportCmd.Flags().String("keystore-passphrase", "", "Passphrase of keystore")
}

func accountImportBindings() {
	if err := viper.BindPFlag("key", accountImportCmd.Flags().Lookup("key")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("keystore", accountImportCmd.Flags().Lookup("keystore")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("keystore-passphrase", accountImportCmd.Flags().Lookup("keystore-passphrase")); err != nil {
		panic(err)
	}
}
