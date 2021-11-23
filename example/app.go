// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package example

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/opentelemetry"
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

	// goerli setup
	goerliClient := evmclient.NewEVMClient()
	if err != nil {
		panic(err)
	}

	err = goerliClient.Configurate(viper.GetString(config.ChainConfigFlagName), "config_goerli.json")
	if err != nil {
		panic(err)
	}
	goerliCfg := goerliClient.GetConfig()
	eventHandlerGoerli := listener.NewETHEventHandler(common.HexToAddress(goerliCfg.SharedEVMConfig.Bridge), goerliClient)
	eventHandlerGoerli.RegisterEventHandler(goerliCfg.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	goerliListener := listener.NewEVMListener(goerliClient, eventHandlerGoerli, common.HexToAddress(goerliCfg.SharedEVMConfig.Bridge))
	mhGoerli := voter.NewEVMMessageHandler(goerliClient, common.HexToAddress(goerliCfg.SharedEVMConfig.Bridge))
	mhGoerli.RegisterMessageHandler(common.HexToAddress(goerliCfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	goerliVoter := voter.NewVoter(mhGoerli, goerliClient, evmtransaction.NewTransaction, evmgaspricer.NewLondonGasPriceClient(goerliClient, nil))
	goerliChain := evm.NewEVMChain(goerliListener, goerliVoter, db, *goerliCfg.SharedEVMConfig.GeneralChainConfig.Id, &goerliCfg.SharedEVMConfig)

	// rinkeby setup
	rinkebyClient := evmclient.NewEVMClient()
	err = rinkebyClient.Configurate(viper.GetString(config.ChainConfigFlagName), "config_rinkeby.json")
	if err != nil {
		panic(err)
	}
	ethCfg := rinkebyClient.GetConfig()
	eventHandler := listener.NewETHEventHandler(common.HexToAddress(ethCfg.SharedEVMConfig.Bridge), rinkebyClient)
	eventHandler.RegisterEventHandler(ethCfg.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	evmListener := listener.NewEVMListener(rinkebyClient, eventHandler, common.HexToAddress(ethCfg.SharedEVMConfig.Bridge))
	mhRinkeby := voter.NewEVMMessageHandler(rinkebyClient, common.HexToAddress(ethCfg.SharedEVMConfig.Bridge))
	mhRinkeby.RegisterMessageHandler(common.HexToAddress(ethCfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	rinkebyVoter := voter.NewVoter(mhRinkeby, rinkebyClient, evmtransaction.NewTransaction, evmgaspricer.NewLondonGasPriceClient(rinkebyClient, nil))
	rinkebyChain := evm.NewEVMChain(evmListener, rinkebyVoter, db, *ethCfg.SharedEVMConfig.GeneralChainConfig.Id, &ethCfg.SharedEVMConfig)

	// celo setup
	celoClient := evmclient.NewEVMClient()
	err = celoClient.Configurate(viper.GetString(config.ChainConfigFlagName), "config_celo.json")
	if err != nil {
		panic(err)
	}
	celoCfg := celoClient.GetConfig()
	eventHandler = listener.NewETHEventHandler(common.HexToAddress(celoCfg.SharedEVMConfig.Bridge), celoClient)
	eventHandler.RegisterEventHandler(celoCfg.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)
	celoListener := listener.NewEVMListener(celoClient, eventHandler, common.HexToAddress(celoCfg.SharedEVMConfig.Bridge))
	mhCelo := voter.NewEVMMessageHandler(celoClient, common.HexToAddress(celoCfg.SharedEVMConfig.Bridge))
	mhCelo.RegisterMessageHandler(common.HexToAddress(celoCfg.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)
	celoVoter := voter.NewVoter(mhCelo, celoClient, evmtransaction.NewTransaction, evmgaspricer.NewLondonGasPriceClient(celoClient, nil))
	celoChain := evm.NewEVMChain(celoListener, celoVoter, db, *celoCfg.SharedEVMConfig.GeneralChainConfig.Id, &celoCfg.SharedEVMConfig)

	r := relayer.NewRelayer([]relayer.RelayedChain{rinkebyChain, goerliChain, celoChain}, &opentelemetry.ConsoleTelemetry{})

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
