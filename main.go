// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"fmt"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ethereum/go-ethereum/common"
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
	RegisterNewETHEventHandler,
	NewEVMListener,
	NewEVMMessageHandler,
	NewEVMVoter,
	NewEVMChain,
)

// var evmSet = wire.NewSet(
// 	NewEVMevmClientWithConfig,
// 	RegisterNewETHEventHandler,
// 	NewEVMListener,
// 	NewEVMMessageHandler,
// 	NewCeloVoter,
// 	NewCeloChain,
// )

func InitializeEvent(
	cnfg Config,
	txFabric voter.TxFabric,
) *evm.EVMChain {
	wire.Build(NewLvlDB, CeloSet)
	return &evm.EVMChain{}
}

func NewLvlDB(cnfg Config) *lvldb.LVLDB {
	db, err := lvldb.NewLvlDB(cnfg.blockstoreFlagName)
	if err != nil {
		panic(err)
	}
	return db
}

func NewEVMCeloClientWithConfig(cnfg Config) *evmclient.EVMClient {
	celoClient := evmclient.NewEVMClient()
	err := celoClient.Configurate(cnfg.configFlagName, "config_celo.json")
	if err != nil {
		panic(err)
	}
	return celoClient
}

func NewEVMevmClientWithConfig(cnfg Config) *evmclient.EVMClient {
	evmClient := evmclient.NewEVMClient()
	err := evmClient.Configurate(cnfg.configFlagName, "config_evm.json")
	if err != nil {
		panic(err)
	}
	return evmClient
}

func RegisterNewETHEventHandler(cnfg Config, client *evmclient.EVMClient) *listener.ETHEventHandler {
	cc := client.GetConfig()
	eh := listener.NewETHEventHandler(common.HexToAddress(cc.SharedEVMConfig.Bridge), client)
	eh.RegisterEventHandler(cc.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	return eh
}

func NewEVMListener(cnfg Config, client *evmclient.EVMClient, eventHandler *listener.ETHEventHandler) *listener.EVMListener {
	cc := client.GetConfig()
	return listener.NewEVMListener(client, eventHandler, common.HexToAddress(cc.SharedEVMConfig.Bridge))
}

func NewEVMMessageHandler(cnfg Config, client *evmclient.EVMClient) *voter.EVMMessageHandler {
	cc := client.GetConfig()
	mh := voter.NewEVMMessageHandler(client, common.HexToAddress(cc.SharedEVMConfig.Bridge))
	mh.RegisterMessageHandler(common.HexToAddress(cc.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	return mh
}

func NewEVMVoter(mh *voter.EVMMessageHandler, client *evmclient.EVMClient, fabric voter.TxFabric) *voter.EVMVoter {
	return voter.NewVoter(mh, client, fabric)
}

func NewEVMChain(client *evmclient.EVMClient, listener *listener.EVMListener, voter *voter.EVMVoter, db *lvldb.LVLDB) *evm.EVMChain {
	cc := client.GetConfig()
	return evm.NewEVMChain(listener, voter, db, *cc.SharedEVMConfig.GeneralChainConfig.Id, &cc.SharedEVMConfig)
}
