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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aaron-alderman/ethdo/util"
	eth2client "github.com/attestantio/go-eth2-client"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var chainInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Obtain information about a chain",
	Long: `Obtain information about a chain.  For example:

    ethdo chain info

In quiet mode this will return 0 if the chain information can be obtained, otherwise 1.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		eth2Client, err := util.ConnectToBeaconNode(ctx, viper.GetString("connection"), viper.GetDuration("timeout"), viper.GetBool("allow-insecure-connections"))
		errCheck(err, "Failed to connect to Ethereum 2 beacon node")

		config, err := eth2Client.(eth2client.SpecProvider).Spec(ctx)
		errCheck(err, "Failed to obtain beacon chain specification")

		genesis, err := eth2Client.(eth2client.GenesisProvider).Genesis(ctx)
		errCheck(err, "Failed to obtain beacon chain genesis")

		fork, err := eth2Client.(eth2client.ForkProvider).Fork(ctx, "head")
		errCheck(err, "Failed to obtain current fork")

		if quiet {
			os.Exit(_exitSuccess)
		}

		if genesis.GenesisTime.Unix() == 0 {
			fmt.Println("Genesis time: undefined")
		} else {
			fmt.Printf("Genesis time: %s\n", genesis.GenesisTime.Format(time.UnixDate))
			outputIf(verbose, fmt.Sprintf("Genesis timestamp: %v", genesis.GenesisTime.Unix()))
		}
		fmt.Printf("Genesis validators root: %#x\n", genesis.GenesisValidatorsRoot)
		fmt.Printf("Genesis fork version: %#x\n", config["GENESIS_FORK_VERSION"].(spec.Version))
		fmt.Printf("Current fork version: %#x\n", fork.CurrentVersion)
		if verbose {
			forkData := &spec.ForkData{
				CurrentVersion:        fork.CurrentVersion,
				GenesisValidatorsRoot: genesis.GenesisValidatorsRoot,
			}
			forkDataRoot, err := forkData.HashTreeRoot()
			if err == nil {
				var forkDigest spec.ForkDigest
				copy(forkDigest[:], forkDataRoot[:])
				fmt.Printf("Fork digest: %#x\n", forkDigest)
			}
		}
		fmt.Printf("Seconds per slot: %d\n", int(config["SECONDS_PER_SLOT"].(time.Duration).Seconds()))
		fmt.Printf("Slots per epoch: %d\n", config["SLOTS_PER_EPOCH"].(uint64))

		os.Exit(_exitSuccess)
	},
}

func init() {
	chainCmd.AddCommand(chainInfoCmd)
	chainFlags(chainInfoCmd)
}
