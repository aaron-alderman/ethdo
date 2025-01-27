// Copyright © 2019 - 2022 Weald Technology Trading
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

	validatordepositdata "github.com/aaron-alderman/ethdo/cmd/validator/depositdata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validatorDepositDataCmd = &cobra.Command{
	Use:   "depositdata",
	Short: "Generate deposit data for one or more validators",
	Long: `Generate data for deposits to the Ethereum 1 validator contract.  For example:

    ethdo validator depositdata --validatoraccount=primary/validator --withdrawalaccount=primary/current --value="32 Ether"

If validatoraccount is provided with an account path it will generate deposit data for all matching accounts.

The information generated can be passed to ethereal to create a deposit from the Ethereum 1 chain.

In quiet mode this will return 0 if the the data can be generated correctly, otherwise 1.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := validatordepositdata.Run(cmd)
		if err != nil {
			return err
		}
		if viper.GetBool("quiet") {
			return nil
		}
		fmt.Println(res)
		return nil
	},
}

func init() {
	validatorCmd.AddCommand(validatorDepositDataCmd)
	validatorFlags(validatorDepositDataCmd)
	validatorDepositDataCmd.Flags().String("validatoraccount", "", "Account carrying out the validation")
	validatorDepositDataCmd.Flags().String("withdrawalaccount", "", "Account to which the validator funds will be withdrawn")
	validatorDepositDataCmd.Flags().String("withdrawalpubkey", "", "Public key of the account to which the validator funds will be withdrawn")
	validatorDepositDataCmd.Flags().String("withdrawaladdress", "", "Ethereum 1 address of the account to which the validator funds will be withdrawn")
	validatorDepositDataCmd.Flags().String("depositvalue", "", "Value of the amount to be deposited")
	validatorDepositDataCmd.Flags().Bool("raw", false, "Print raw deposit data transaction data")
	validatorDepositDataCmd.Flags().String("forkversion", "", "Use a hard-coded fork version (default is to use mainnet value)")
	validatorDepositDataCmd.Flags().Bool("launchpad", false, "Print launchpad-compatible JSON")
}

func validatorDepositdataBindings() {
	if err := viper.BindPFlag("validatoraccount", validatorDepositDataCmd.Flags().Lookup("validatoraccount")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("withdrawalaccount", validatorDepositDataCmd.Flags().Lookup("withdrawalaccount")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("withdrawalpubkey", validatorDepositDataCmd.Flags().Lookup("withdrawalpubkey")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("withdrawaladdress", validatorDepositDataCmd.Flags().Lookup("withdrawaladdress")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("depositvalue", validatorDepositDataCmd.Flags().Lookup("depositvalue")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("raw", validatorDepositDataCmd.Flags().Lookup("raw")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("forkversion", validatorDepositDataCmd.Flags().Lookup("forkversion")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("launchpad", validatorDepositDataCmd.Flags().Lookup("launchpad")); err != nil {
		panic(err)
	}
}
