// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package example

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"

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

	// CELO1 setup
	celo1Client := evmclient.NewEVMClient()
	err = celo1Client.Configurate(viper.GetString(config.ConfigFlagName), "config_celo1.json")
	if err != nil {
		panic(err)
	}
	celo1Cfg := celo1Client.GetConfig()

	eventHandler := listener.NewETHEventHandler(common.HexToAddress(celo1Cfg.SharedEVMConfig.Bridge), celo1Client)
	eventHandler.RegisterEventHandler(celo1Cfg.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	celoListener1 := listener.NewEVMListener(celo1Client, eventHandler, common.HexToAddress(celo1Cfg.SharedEVMConfig.Bridge))
	mh := voter.NewEVMMessageHandler(celo1Client, common.HexToAddress(celo1Cfg.SharedEVMConfig.Bridge))
	mh.RegisterMessageHandler(common.HexToAddress(celo1Cfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	celoVoter1 := voter.NewVoter(mh, celo1Client, transaction.NewCeloTransaction)
	celoChain1 := evm.NewEVMChain(celoListener1, celoVoter1, db, *celo1Cfg.SharedEVMConfig.GeneralChainConfig.Id, &celo1Cfg.SharedEVMConfig)

	//// CELO2 setup
	celoClient2 := evmclient.NewEVMClient()
	err = celoClient2.Configurate(viper.GetString(config.ConfigFlagName), "config_celo2.json")
	if err != nil {
		panic(err)
	}
	celoConfig2 := celoClient2.GetConfig()

	eventHandlerCelo2 := listener.NewETHEventHandler(common.HexToAddress(celoConfig2.SharedEVMConfig.Bridge), celoClient2)
	eventHandlerCelo2.RegisterEventHandler(celoConfig2.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	goerliListener := listener.NewEVMListener(celoClient2, eventHandlerCelo2, common.HexToAddress(celoConfig2.SharedEVMConfig.Bridge))
	mhCelo2 := voter.NewEVMMessageHandler(celoClient2, common.HexToAddress(celoConfig2.SharedEVMConfig.Bridge))
	mhCelo2.RegisterMessageHandler(common.HexToAddress(celoConfig2.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	celoVoter2 := voter.NewVoter(mhCelo2, celoClient2, transaction.NewCeloTransaction)
	celoChain2 := evm.NewEVMChain(goerliListener, celoVoter2, db, *celoConfig2.SharedEVMConfig.GeneralChainConfig.Id, &celoConfig2.SharedEVMConfig)

	r := relayer.NewRelayer([]relayer.RelayedChain{celoChain1, celoChain2})

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
