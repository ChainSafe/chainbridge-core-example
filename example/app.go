// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package example

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-core-example/example/keystore"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"

	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
var EveKp = keystore.TestKeyRing.EthereumKeys[keystore.EveKey]

var (
	DefaultRelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}
)

//Bridge:             0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B
//Erc20 Handler:      0x3167776db165D8eA0f51790CA2bbf44Db5105ADF
func Run() error {
	errChn := make(chan error)
	stopChn := make(chan struct{})

	db, err := lvldb.NewLvlDB(viper.GetString(config.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}

	// rinkeby
	ethClient := evmclient.NewEVMClient()
	err = ethClient.Configurate(viper.GetString(config.ConfigFlagName), "config_rinkeby")
	if err != nil {
		panic(err)
	}
	ethCfg := ethClient.GetConfig()

	eventHandler := listener.NewETHEventHandler(common.HexToAddress(ethCfg.SharedEVMConfig.Bridge), ethClient)
	eventHandler.RegisterEventHandler(ethCfg.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	evmListener := listener.NewEVMListener(ethClient, eventHandler, common.HexToAddress(ethCfg.SharedEVMConfig.Bridge))
	mh := voter.NewEVMMessageHandler(ethClient, common.HexToAddress(ethCfg.SharedEVMConfig.Bridge))
	mh.RegisterMessageHandler(common.HexToAddress(ethCfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	evmVoter := voter.NewVoter(mh, ethClient)
	evmChain := evm.NewEVMChain(evmListener, evmVoter, db, *ethCfg.SharedEVMConfig.GeneralChainConfig.Id, &ethCfg.SharedEVMConfig)

	// goerli setup
	goerliClient := evmclient.NewEVMClient()
	if err != nil {
		panic(err)
	}
	goerliCfg := ethClient.GetConfig()
	err = goerliClient.Configurate(viper.GetString(config.ConfigFlagName), "config_goerli")
	if err != nil {
		panic(err)
	}
	eventHandlerGoerli := listener.NewETHEventHandler(common.HexToAddress(goerliCfg.SharedEVMConfig.Bridge), goerliClient)
	eventHandlerGoerli.RegisterEventHandler(goerliCfg.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	goerliListener := listener.NewEVMListener(goerliClient, eventHandlerGoerli, common.HexToAddress(goerliCfg.SharedEVMConfig.Bridge))
	mhGoerli := voter.NewEVMMessageHandler(goerliClient, common.HexToAddress(goerliCfg.SharedEVMConfig.Bridge))
	mhGoerli.RegisterMessageHandler(common.HexToAddress(goerliCfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	goerliVoter := voter.NewVoter(mhGoerli, goerliClient)
	goerliChain := evm.NewEVMChain(goerliListener, goerliVoter, db, 2, &goerliCfg.SharedEVMConfig)

	r := relayer.NewRelayer([]relayer.RelayedChain{evmChain, goerliChain})

	go r.Start(stopChn, errChn)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	select {
	case err := <-errChn:
		log.Error().Err(err).Msg("failed to listen and serve")
		close(stopChn)
		return err
	case sig := <-sysErr:
		log.Info().Msgf("terminating got [%v] signal", sig)
		return nil
	}
}
