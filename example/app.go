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

	// EVM1 setup
	evmClient := evmclient.NewEVMClient()
	err = evmClient.Configurate(viper.GetString(config.ConfigFlagName), "config_evm.json")
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

	// EVM2 setup
	evmClient2 := evmclient.NewEVMClient()
	err = evmClient.Configurate(viper.GetString(config.ConfigFlagName), "config_evm.json")
	if err != nil {
		panic(err)
	}
	evmConfig2 := evmClient2.GetConfig()
	eventHandlerEVM2 := listener.NewETHEventHandler(
		common.HexToAddress(evmConfig2.SharedEVMConfig.Bridge), evmClient2,
	)
	eventHandlerEVM2.RegisterEventHandler(
		evmConfig2.SharedEVMConfig.Erc20Handler,
		listener.Erc20EventHandler,
	)
	evmListener2 := listener.NewEVMListener(
		evmClient2,
		eventHandlerEVM2,
		common.HexToAddress(evmConfig2.SharedEVMConfig.Bridge),
	)
	mhEVM2 := voter.NewEVMMessageHandler(
		evmClient2,
		common.HexToAddress(evmConfig2.SharedEVMConfig.Bridge),
	)
	mhEVM2.RegisterMessageHandler(
		common.HexToAddress(evmConfig2.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler,
	)
	evmVoter2 := voter.NewVoter(
		mhEVM2,
		evmClient2,
		evmtransaction.NewTransaction,
	)
	evmChain2 := evm.NewEVMChain(
		evmListener2,
		evmVoter2,
		db,
		*evmConfig2.SharedEVMConfig.GeneralChainConfig.Id,
		&evmConfig2.SharedEVMConfig,
	)

	r := relayer.NewRelayer([]relayer.RelayedChain{evmChain, evmChain2})

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
