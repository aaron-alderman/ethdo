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

package chainqueues

import (
	"context"
	"fmt"

	standardchaintime "github.com/aaron-alderman/ethdo/services/chaintime/standard"
	"github.com/aaron-alderman/ethdo/util"
	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
)

func (c *command) process(ctx context.Context) error {
	// Obtain information we need to process.
	if err := c.setup(ctx); err != nil {
		return err
	}

	epoch, err := util.ParseEpoch(ctx, c.chainTime, c.epoch)
	if err != nil {
		return err
	}

	validators, err := c.validatorsProvider.Validators(ctx, fmt.Sprintf("%d", c.chainTime.FirstSlotOfEpoch(epoch)), nil)
	if err != nil {
		return errors.Wrap(err, "failed to obtain validators")
	}

	for _, validator := range validators {
		if validator.Validator == nil {
			continue
		}
		if validator.Validator.ActivationEligibilityEpoch <= epoch && validator.Validator.ActivationEpoch > epoch {
			c.activationQueue++
		}
		if validator.Validator.ExitEpoch != 0xffffffffffffffff && validator.Validator.ExitEpoch > epoch {
			c.exitQueue++
		}
	}

	return nil
}

func (c *command) setup(ctx context.Context) error {
	var err error

	// Connect to the client.
	c.eth2Client, err = util.ConnectToBeaconNode(ctx, c.connection, c.timeout, c.allowInsecureConnections)
	if err != nil {
		return errors.Wrap(err, "failed to connect to beacon node")
	}

	c.chainTime, err = standardchaintime.New(ctx,
		standardchaintime.WithSpecProvider(c.eth2Client.(eth2client.SpecProvider)),
		standardchaintime.WithForkScheduleProvider(c.eth2Client.(eth2client.ForkScheduleProvider)),
		standardchaintime.WithGenesisTimeProvider(c.eth2Client.(eth2client.GenesisTimeProvider)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to set up chaintime service")
	}

	var isProvider bool
	c.validatorsProvider, isProvider = c.eth2Client.(eth2client.ValidatorsProvider)
	if !isProvider {
		return errors.New("connection does not provide validator information")
	}

	return nil
}
