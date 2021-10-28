package main

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/google/wire"
)

var CeloSet = wire.NewSet(
	NewEVMCeloClientWithConfig,
	RegisterNewCeloEventHandler,
	NewCeloEVMListener,
	NewCeloEVMMessageHandler,
	NewCeloEVMVoter,
	NewCeloEVMChain,
)

type CeloEVMClient *evmclient.EVMClient

func NewEVMCeloClientWithConfig(cnfg Config) CeloEVMClient {
	return _newEVMClientWithConfig(cnfg, "celo_config.json")
}

type CeloEventListener *listener.ETHEventHandler

func RegisterNewCeloEventHandler(cnfg Config, client CeloEVMClient) CeloEventListener {
	return _registerNewEventHandler(cnfg, client)
}

type CeloEVMListener *listener.EVMListener

func NewCeloEVMListener(cnfg Config, client CeloEVMClient, eventHandler CeloEventListener) CeloEVMListener {
	return _newEVMListener(cnfg, client, eventHandler)
}

type CeloEVMMessageHandler *voter.EVMMessageHandler

func NewCeloEVMMessageHandler(cnfg Config, client CeloEVMClient) CeloEVMMessageHandler {
	return _newEVMMessageHandler(cnfg, client)
}

type CeloEVMVoter *voter.EVMVoter

func NewCeloEVMVoter(mh CeloEVMMessageHandler, client CeloEVMClient, fabric voter.TxFabric) CeloEVMVoter {
	return _newEVMVoter(mh, client, fabric)
}

type CeloEVMChain *evm.EVMChain

func NewCeloEVMChain(client CeloEVMClient, listener CeloEVMListener, voter CeloEVMVoter, db *lvldb.LVLDB) CeloEVMChain {
	return _newEVMChain(client, listener, voter, db)
}
