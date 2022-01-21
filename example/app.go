// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package example

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/chainbridge-celo-module/transaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/ChainSafe/chainbridge-core/flags"
	"github.com/ChainSafe/chainbridge-core/lvldb"
	"github.com/ChainSafe/chainbridge-core/opentelemetry"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ChainSafe/chainbridge-core/store"
	optimism "github.com/ChainSafe/chainbridge-optimism-module"
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
	blockstore := store.NewBlockStore(db)

	chains := []relayer.RelayedChain{}
	for _, chainConfig := range configuration.ChainConfigs {
		switch chainConfig["type"] {
		case "evm":
			{
				chain, err := evm.SetupDefaultEVMChain(chainConfig, evmtransaction.NewTransaction, blockstore)
				if err != nil {
					panic(err)
				}

				chains = append(chains, chain)
			}
		case "celo":
			{
				chain, err := evm.SetupDefaultEVMChain(chainConfig, transaction.NewCeloTransaction, blockstore)
				if err != nil {
					panic(err)
				}

				chains = append(chains, chain)
			}
		case "optimism":
			{
				chain, err := optimism.SetupDefaultOptimismChain(chainConfig, evmtransaction.NewTransaction, blockstore)
				if err != nil {
					panic(err)
				}

				chains = append(chains, chain)
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
