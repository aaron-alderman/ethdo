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
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aaron-alderman/ethdo/util"
	eth2client "github.com/attestantio/go-eth2-client"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

var exitVerifyPubKey string

var exitVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify exit data is valid",
	Long: `Verify that exit data generated by "ethdo validator exit" is correct for a given account.  For example:

    ethdo exit verify --data=exitdata.json --account=primary/current

In quiet mode this will return 0 if the the exit is verified correctly, otherwise 1.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		assert(viper.GetString("account") != "" || exitVerifyPubKey != "", "account or public key is required")
		account, err := exitVerifyAccount(ctx)
		errCheck(err, "Failed to obtain account")

		assert(viper.GetString("exit") != "", "exit is required")
		data, err := obtainExitData(viper.GetString("exit"))
		errCheck(err, "Failed to obtain exit data")

		// Confirm signature is good.
		eth2Client, err := util.ConnectToBeaconNode(ctx, viper.GetString("connection"), viper.GetDuration("timeout"), viper.GetBool("allow-insecure-connections"))
		errCheck(err, "Failed to connect to Ethereum 2 beacon node")

		genesis, err := eth2Client.(eth2client.GenesisProvider).Genesis(ctx)
		errCheck(err, "Failed to obtain beacon chain genesis")

		domain := e2types.Domain(e2types.DomainVoluntaryExit, data.ForkVersion[:], genesis.GenesisValidatorsRoot[:])
		var exitDomain spec.Domain
		copy(exitDomain[:], domain)
		exit := &spec.VoluntaryExit{
			Epoch:          data.Exit.Message.Epoch,
			ValidatorIndex: data.Exit.Message.ValidatorIndex,
		}
		exitRoot, err := exit.HashTreeRoot()
		errCheck(err, "Failed to obtain exit hash tree root")
		signatureBytes := make([]byte, 96)
		copy(signatureBytes, data.Exit.Signature[:])
		sig, err := e2types.BLSSignatureFromBytes(signatureBytes)
		errCheck(err, "Invalid signature")
		verified, err := util.VerifyRoot(account, exitRoot, exitDomain, sig)
		errCheck(err, "Failed to verify voluntary exit")
		assert(verified, "Voluntary exit failed to verify")

		fork, err := eth2Client.(eth2client.ForkProvider).Fork(ctx, "head")
		errCheck(err, "Failed to obtain current fork")
		assert(bytes.Equal(data.ForkVersion[:], fork.CurrentVersion[:]) || bytes.Equal(data.ForkVersion[:], fork.PreviousVersion[:]), "Exit is for an old fork version and is no longer valid")

		outputIf(verbose, "Verified")
		os.Exit(_exitSuccess)
	},
}

// obtainExitData obtains exit data from an input, could be JSON itself or a path to JSON.
func obtainExitData(input string) (*util.ValidatorExitData, error) {
	var err error
	var data []byte
	// Input could be JSON or a path to JSON
	if strings.HasPrefix(input, "{") {
		// Looks like JSON
		data = []byte(input)
	} else {
		// Assume it's a path to JSON
		data, err = ioutil.ReadFile(input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find deposit data file")
		}
	}
	exitData := &util.ValidatorExitData{}
	err = json.Unmarshal(data, exitData)
	if err != nil {
		return nil, errors.Wrap(err, "data is not valid JSON")
	}

	return exitData, nil
}

// exitVerifyAccount obtains the account for the exitVerify command.
func exitVerifyAccount(ctx context.Context) (e2wtypes.Account, error) {
	var account e2wtypes.Account
	var err error
	if viper.GetString("account") != "" {
		_, account, err = walletAndAccountFromPath(ctx, viper.GetString("account"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain account")
		}
	} else {
		pubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(exitVerifyPubKey, "0x"))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode public key %s", exitVerifyPubKey))
		}
		account, err = util.NewScratchAccount(nil, pubKeyBytes)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("invalid public key %s", exitVerifyPubKey))
		}
	}
	return account, nil
}

func init() {
	exitCmd.AddCommand(exitVerifyCmd)
	exitFlags(exitVerifyCmd)
	exitVerifyCmd.Flags().String("exit", "", "JSON data, or path to JSON data")
	exitVerifyCmd.Flags().StringVar(&exitVerifyPubKey, "pubkey", "", "Public key for which to verify exit")
}

func exitVerifyBindings() {
	if err := viper.BindPFlag("exit", exitVerifyCmd.Flags().Lookup("exit")); err != nil {
		panic(err)
	}
}
