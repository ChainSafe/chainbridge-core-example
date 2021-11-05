// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"fmt"

	celotx "github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
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
	celloChain := CreateEVMCelloChain(
		cnfg,
		celotx.NewCeloTransaction,
	)

	evmChain := CreateEVMChain(
		cnfg,
		evmtransaction.NewTransaction,
	)

	fmt.Println(celloChain)
	fmt.Println(evmChain)
}

func CreateEVMCelloChain(
	cnfg Config,
	txFabric voter.TxFabric,
) CeloEVMChain {
	wire.Build(NewLvlDB, CeloSet)
	return &evm.EVMChain{}
}

func CreateEVMChain(
	cnfg Config,
	txFabric voter.TxFabric,
) EVMChain {
	wire.Build(NewLvlDB, EvmSet)
	return &evm.EVMChain{}
}
