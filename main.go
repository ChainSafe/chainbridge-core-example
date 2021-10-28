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

	// Decide depending on flags
	e := CreateEVMCelloChains(
		cnfg,
		transaction.NewCeloTransaction,
	)

	fmt.Println(e)
}

func CreateEVMCelloChains(
	cnfg Config,
	txFabric voter.TxFabric,
) EVMChain {
	wire.Build(NewLvlDB, CeloSet, EvmSet)
	return &evm.EVMChain{}
}

// func CreateEVMSubstrateChains(
// 	cnfg Config,
// 	txFabric voter.TxFabric,
// ) {
// 	wire.Build(NewLvlDB, SubstrateSet, EvmSet)
// }
