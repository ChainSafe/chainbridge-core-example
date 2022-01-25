// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

//nolint
package e2e_test

import (
	"testing"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/stretchr/testify/suite"
)

const ETHEndpoint = "ws://localhost:8846"
const CeloEndpoint = "ws://localhost:8545"

var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	ethClient, err := evmclient.NewEVMClientFromParams(ETHEndpoint, local.EveKp.PrivateKey())
	if err != nil {
		panic(err)
	}

	celoClient, err := evmclient.NewEVMClientFromParams(CeloEndpoint, local.EveKp.PrivateKey())
	if err != nil {
		panic(err)
	}

	suite.Run(
		t,
		evm.SetupEVM2EVMTestSuite(
			evmtransaction.NewTransaction,
			transaction.NewCeloTransaction,
			ethClient,
			celoClient,
			local.DefaultRelayerAddresses,
			local.DefaultRelayerAddresses,
		),
	)
}
