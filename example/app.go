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
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/flags"
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

	configuration, err := config.GetConfig(viper.GetString(flags.ConfigFlagName))

	db, err := lvldb.NewLvlDB(viper.GetString(flags.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}

	chains := []relayer.RelayedChain{}
	for _, chainConfig := range configuration.ChainConfigs {
		switch chainConfig["type"] {
		case "evm":
			{
				chain, err := evm.SetupDefaultEVMChain(chainConfig, db)
				if err != nil {
					panic(err)
				}

				chains = append(chains, chain)
			}
		case "celo":
			{
				config, err := chain.NewEVMConfig(chainConfig)
				if err != nil {
					panic(err)
				}

				client := evmclient.NewEVMClient()
				err = client.Configurate(config)
				if err != nil {
					panic(err)
				}

				eventHandler := listener.NewETHEventHandler(common.HexToAddress(config.Bridge), client)
				eventHandler.RegisterEventHandler(config.Erc20Handler, listener.Erc20EventHandler)
				eventHandler.RegisterEventHandler(config.Erc721Handler, listener.Erc721EventHandler)
				eventHandler.RegisterEventHandler(config.GenericHandler, listener.GenericEventHandler)
				evmListener := listener.NewEVMListener(client, eventHandler, common.HexToAddress(config.Bridge))

				mh := voter.NewEVMMessageHandler(client, common.HexToAddress(config.Bridge))
				mh.RegisterMessageHandler(config.Erc20Handler, voter.ERC20MessageHandler)
				mh.RegisterMessageHandler(config.Erc721Handler, voter.ERC721MessageHandler)
				mh.RegisterMessageHandler(config.GenericHandler, voter.GenericMessageHandler)
				evmVoter := voter.NewVoter(mh, client, transaction.NewCeloTransaction, evmgaspricer.NewLondonGasPriceClient(client, nil))

				chains = append(chains, evm.NewEVMChain(evmListener, evmVoter, db, config))
			}
		}
	}

	r := relayer.NewRelayer(chains, &opentelemetry.ConsoleTelemetry{})
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
