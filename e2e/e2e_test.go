// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

//nolint
package e2e_test

import (
	"testing"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/stretchr/testify/suite"
)

const ETHEndpoint = "http://localhost:8845"
const CeloEndpoint = "http://localhost:8546"

var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	suite.Run(t, evm.SetupEVM2EVMTestSuite(evmtransaction.NewTransaction, transaction.NewCeloTransaction, ETHEndpoint, CeloEndpoint, EveKp))
}
