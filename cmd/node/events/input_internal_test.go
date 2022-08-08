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

package nodeevents

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aaron-alderman/ethdo/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	e2wallet "github.com/wealdtech/go-eth2-wallet"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	nd "github.com/wealdtech/go-eth2-wallet-nd/v2"
	scratch "github.com/wealdtech/go-eth2-wallet-store-scratch"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestInput(t *testing.T) {
	if os.Getenv("ETHDO_TEST_CONNECTION") == "" {
		t.Skip("ETHDO_TEST_CONNECTION not configured; cannot run tests")
	}

	require.NoError(t, e2types.InitBLS())

	store := scratch.New()
	require.NoError(t, e2wallet.UseStore(store))
	testWallet, err := nd.CreateWallet(context.Background(), "Test wallet", store, keystorev4.New())
	require.NoError(t, err)
	require.NoError(t, testWallet.(e2wtypes.WalletLocker).Unlock(context.Background(), nil))
	viper.Set("passphrase", "pass")
	_, err = testWallet.(e2wtypes.WalletAccountImporter).ImportAccount(context.Background(),
		"Interop 0",
		testutil.HexToBytes("0x25295f0d1d592a90b333e26e85149708208e9f8e8bc18f6c77bd62f8ad7a6866"),
		[]byte("pass"),
	)
	require.NoError(t, err)

	tests := []struct {
		name string
		vars map[string]interface{}
		res  *dataIn
		err  string
	}{
		{
			name: "TimeoutMissing",
			vars: map[string]interface{}{},
			err:  "timeout is required",
		},
		{
			name: "ConnectionMissing",
			vars: map[string]interface{}{
				"timeout": "5s",
			},
			err: "failed to connect to any beacon node",
		},
		{
			name: "ConnectionBad",
			vars: map[string]interface{}{
				"timeout":    "5s",
				"connection": "localhost:1",
				"topics":     []string{"one", "two"},
			},
			err: "failed to connect to beacon node: failed to confirm node connection: failed to fetch genesis: failed to request genesis: failed to call GET endpoint: Get \"http://localhost:1/eth/v1/beacon/genesis\": dial tcp 127.0.0.1:1: connect: connection refused",
		},
		{
			name: "TopicsNil",
			vars: map[string]interface{}{
				"timeout":    "5s",
				"connection": os.Getenv("ETHDO_TEST_CONNECTION"),
			},
			res: &dataIn{
				timeout: 5 * time.Second,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			viper.Reset()

			for k, v := range test.vars {
				viper.Set(k, v)
			}
			res, err := input(context.Background())
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.res.timeout, res.timeout)
				require.Equal(t, test.res.topics, res.topics)
			}
		})
	}
}
