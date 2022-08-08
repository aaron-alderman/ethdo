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

	blockinfo "github.com/aaron-alderman/ethdo/cmd/block/info"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var blockInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Obtain information about a block",
	Long: `Obtain information about a block.  For example:

    ethdo block info --blockid=12345

In quiet mode this will return 0 if the block information is present and not skipped, otherwise 1.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := blockinfo.Run(cmd)
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
	blockCmd.AddCommand(blockInfoCmd)
	blockFlags(blockInfoCmd)
	blockInfoCmd.Flags().String("blockid", "head", "the ID of the block to fetch")
	blockInfoCmd.Flags().Bool("stream", false, "continually stream blocks as they arrive")
	blockInfoCmd.Flags().Bool("json", false, "output data in JSON format")
	blockInfoCmd.Flags().Bool("ssz", false, "output data in SSZ format")
}

func blockInfoBindings() {
	if err := viper.BindPFlag("blockid", blockInfoCmd.Flags().Lookup("blockid")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("stream", blockInfoCmd.Flags().Lookup("stream")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("json", blockInfoCmd.Flags().Lookup("json")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("ssz", blockInfoCmd.Flags().Lookup("ssz")); err != nil {
		panic(err)
	}
}
