// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

//nolint
package e2e_test

import (
	"testing"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"

	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/stretchr/testify/suite"
)

const ETHEndpoint = "http://localhost:8845"
const CeloEndpoint = "http://localhost:8546"

func TestRunE2ETests(t *testing.T) {
	suite.Run(t, evm.PreSetupTestSuite(evmtransaction.NewTransaction, transaction.NewCeloTransaction, ETHEndpoint, CeloEndpoint))
}
