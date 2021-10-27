// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package example

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"

	"github.com/ChainSafe/chainbridge-core/chains/evm"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"

	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/relayer"
	optimismListener "github.com/ChainSafe/chainbridge-optimism-module/listener"
	"github.com/ChainSafe/chainbridge-optimism-module/optimismclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Run() error {
	errChn := make(chan error)
	stopChn := make(chan struct{})

	db, err := lvldb.NewLvlDB(viper.GetString(config.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}

	// Optimism setup
	ethClientOptimism := optimismclient.NewEVMClient()
	err = ethClientOptimism.Configurate(viper.GetString(config.ConfigFlagName), "config_optimism")
	if err != nil {
		panic(err)
	}
	optimismCfg := ethClientOptimism.GetConfig()

	eventHandlerOptimism := optimismListener.NewETHEventHandler(common.HexToAddress(optimismCfg.SharedEVMConfig.Bridge), ethClientOptimism)
	eventHandlerOptimism.RegisterEventHandler(optimismCfg.SharedEVMConfig.Erc20Handler, optimismListener.Erc20EventHandler)
	evmListenerOptimism := optimismListener.NewEVMListener(ethClientOptimism, eventHandlerOptimism, common.HexToAddress(optimismCfg.SharedEVMConfig.Bridge))
	mhOptimism := voter.NewEVMMessageHandler(ethClientOptimism, common.HexToAddress(optimismCfg.SharedEVMConfig.Bridge))
	mhOptimism.RegisterMessageHandler(common.HexToAddress(optimismCfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	evmVoterOptimism := voter.NewVoter(mhOptimism, ethClientOptimism, evmtransaction.NewTransaction)
	optimismChain := evm.NewEVMChain(evmListenerOptimism, evmVoterOptimism, db, *optimismCfg.SharedEVMConfig.GeneralChainConfig.Id, &optimismCfg.SharedEVMConfig)

	////EVM setup
	evmClient := evmclient.NewEVMClient()
	err = evmClient.Configurate(viper.GetString(config.ConfigFlagName), "config_local")
	if err != nil {
		panic(err)
	}
	evmConfig := evmClient.GetConfig()
	eventHandlerEVM := listener.NewETHEventHandler(common.HexToAddress(evmConfig.SharedEVMConfig.Bridge), evmClient)
	eventHandlerEVM.RegisterEventHandler(evmConfig.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	evmListener := listener.NewEVMListener(evmClient, eventHandlerEVM, common.HexToAddress(evmConfig.SharedEVMConfig.Bridge))
	mhEVM := voter.NewEVMMessageHandler(evmClient, common.HexToAddress(evmConfig.SharedEVMConfig.Bridge))
	mhEVM.RegisterMessageHandler(common.HexToAddress(evmConfig.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	evmVoter := voter.NewVoter(mhEVM, evmClient, evmtransaction.NewTransaction)
	evmChain := evm.NewEVMChain(evmListener, evmVoter, db, *evmConfig.SharedEVMConfig.GeneralChainConfig.Id, &evmConfig.SharedEVMConfig)

	r := relayer.NewRelayer([]relayer.RelayedChain{optimismChain, evmChain})

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
		log.Info().Msgf("terminating got ` [%v] signal", sig)
		return nil
	}
}
