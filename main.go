// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/google/wire"
)

type Config struct {
	blockstoreFlagName string
	configFlagName     string
}

func main() {
	// This would be initialized from config/flags
	cnfg := Config{
		blockstoreFlagName: "",
		configFlagName:     "",
	}

	e := InitializeEvent(
		cnfg,
		transaction.NewCeloTransaction,
	)
	fmt.Println(e)

}

var CeloSet = wire.NewSet(
	NewEVMCeloClientWithConfig,
	RegisterNewCeloEventHandler,
	NewCeloEVMListener,
	NewCeloEVMMessageHandler,
	NewCeloEVMVoter,
	NewCeloEVMChain,
)

var EvmSet = wire.NewSet(
	NewEVMClientWithConfig,
	RegisterNewEVMEventHandler,
	NewEVMListener,
	NewEVMMessageHandler,
	NewEVMVoter,
	NewEVMChain,
)

func InitializeEvent(
	cnfg Config,
	txFabric voter.TxFabric,
) EVMChain {
	wire.Build(NewLvlDB, CeloSet, EvmSet)
	return &evm.EVMChain{}
}
