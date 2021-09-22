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

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	//c, err := evmclient.NewEVMClientFromParams(CeloEndpoint, evm.AliceKp.PrivateKey(), big.NewInt(20000000000))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//_, err = c.Simulate(big.NewInt(200), common.HexToHash("0x7025762ad07d80ff6eb560d47e6be146331e5ae66563dd772e2423c16ec09698"), evm.AliceKp.CommonAddress())
	//if err != nil {
	//	t.Fatal(err)
	//}
	suite.Run(t, evm.PreSetupTestSuite(evmtransaction.NewTransaction, transaction.NewCeloTransaction, ETHEndpoint, CeloEndpoint, evm.EveKp))
}
