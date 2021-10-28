package main

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ethereum/go-ethereum/common"
)

func NewLvlDB(cnfg Config) *lvldb.LVLDB {
	db, err := lvldb.NewLvlDB(cnfg.blockstoreFlagName)
	if err != nil {
		panic(err)
	}
	return db
}

func _newEVMClientWithConfig(cnfg Config, name string) *evmclient.EVMClient {
	client := evmclient.NewEVMClient()
	err := client.Configurate(cnfg.configFlagName, name)
	if err != nil {
		panic(err)
	}
	return client
}

func _registerNewEventHandler(cnfg Config, client *evmclient.EVMClient) *listener.ETHEventHandler {
	cc := client.GetConfig()
	eh := listener.NewETHEventHandler(common.HexToAddress(cc.SharedEVMConfig.Bridge), client)
	eh.RegisterEventHandler(cc.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	return eh
}

func _newEVMListener(cnfg Config, client *evmclient.EVMClient, eventHandler *listener.ETHEventHandler) *listener.EVMListener {
	cc := client.GetConfig()
	return listener.NewEVMListener(client, eventHandler, common.HexToAddress(cc.SharedEVMConfig.Bridge))
}

func _newEVMMessageHandler(cnfg Config, client *evmclient.EVMClient) *voter.EVMMessageHandler {
	cc := client.GetConfig()
	mh := voter.NewEVMMessageHandler(client, common.HexToAddress(cc.SharedEVMConfig.Bridge))
	mh.RegisterMessageHandler(common.HexToAddress(cc.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	return mh
}

func _newEVMVoter(mh *voter.EVMMessageHandler, client *evmclient.EVMClient, fabric voter.TxFabric) *voter.EVMVoter {
	return voter.NewVoter(mh, client, fabric)
}

func _newEVMChain(client *evmclient.EVMClient, listener *listener.EVMListener, voter *voter.EVMVoter, db *lvldb.LVLDB) *evm.EVMChain {
	cc := client.GetConfig()
	return evm.NewEVMChain(listener, voter, db, *cc.SharedEVMConfig.GeneralChainConfig.Id, &cc.SharedEVMConfig)
}
