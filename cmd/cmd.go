// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package cmd

import (
	cCLI "github.com/ChainSafe/chainbridge-celo-module/cli"
	"github.com/ChainSafe/chainbridge-core-example/example"
	evmCLI "github.com/ChainSafe/chainbridge-core/chains/evm/cli"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCMD = &cobra.Command{
		Use: "",
	}
	runCMD = &cobra.Command{
		Use:   "run",
		Short: "Run example app",
		Long:  "Run example app",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := example.Run(); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	config.BindFlags(rootCMD)
}

func Execute() {
	rootCMD.AddCommand(runCMD, cCLI.CeloRootCLI, evmCLI.EvmRootCLI)
	if err := rootCMD.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root cmd")
	}
}
