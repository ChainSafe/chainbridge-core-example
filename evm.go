package main

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/google/wire"
)

var EvmSet = wire.NewSet(
	NewEVMClientWithConfig,
	RegisterNewEVMEventHandler,
	NewEVMListener,
	NewEVMMessageHandler,
	NewEVMVoter,
	NewEVMChain,
)

type EVMClient *evmclient.EVMClient

func NewEVMClientWithConfig(cnfg Config) EVMClient {
	return _newEVMClientWithConfig(cnfg, "config_evm.json")
}

type EVMEventListener *listener.ETHEventHandler

func RegisterNewEVMEventHandler(cnfg Config, client EVMClient) EVMEventListener {
	return _registerNewEventHandler(cnfg, client)
}

type EVMListener *listener.EVMListener

func NewEVMListener(cnfg Config, client EVMClient, eventHandler EVMEventListener) EVMListener {
	return _newEVMListener(cnfg, client, eventHandler)
}

type EVMMessageHandler *voter.EVMMessageHandler

func NewEVMMessageHandler(cnfg Config, client EVMClient) EVMMessageHandler {
	return _newEVMMessageHandler(cnfg, client)
}

type EVMVoter *voter.EVMVoter

func NewEVMVoter(mh EVMMessageHandler, client EVMClient, fabric voter.TxFabric) EVMVoter {
	return _newEVMVoter(mh, client, fabric)
}

type EVMChain *evm.EVMChain

func NewEVMChain(client EVMClient, listener EVMListener, voter EVMVoter, db *lvldb.LVLDB) EVMChain {
	return _newEVMChain(client, listener, voter, db)
}
