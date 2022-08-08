// Copyright © 2022 Weald Technology Trading.
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

package proposerduties

import (
	"context"
	"time"

	"github.com/aaron-alderman/ethdo/services/chaintime"
	eth2client "github.com/attestantio/go-eth2-client"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type command struct {
	quiet   bool
	verbose bool
	debug   bool

	// Beacon node connection.
	timeout                  time.Duration
	connection               string
	allowInsecureConnections bool

	// Operation.
	epoch      string
	jsonOutput bool

	// Data access.
	eth2Client             eth2client.Service
	chainTime              chaintime.Service
	proposerDutiesProvider eth2client.ProposerDutiesProvider

	// Results.
	results *results
}

type results struct {
	Epoch  phase0.Epoch          `json:"epoch"`
	Duties []*apiv1.ProposerDuty `json:"duties"`
}

func newCommand(ctx context.Context) (*command, error) {
	c := &command{
		quiet:   viper.GetBool("quiet"),
		verbose: viper.GetBool("verbose"),
		debug:   viper.GetBool("debug"),
		results: &results{},
	}

	// Timeout.
	if viper.GetDuration("timeout") == 0 {
		return nil, errors.New("timeout is required")
	}
	c.timeout = viper.GetDuration("timeout")

	if viper.GetString("connection") == "" {
		return nil, errors.New("connection is required")
	}
	c.connection = viper.GetString("connection")
	c.allowInsecureConnections = viper.GetBool("allow-insecure-connections")

	c.epoch = viper.GetString("epoch")
	c.jsonOutput = viper.GetBool("json")

	return c, nil
}
