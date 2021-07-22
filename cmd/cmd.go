package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "chainbridge-core-example",
		Short: "Chainbridge Core Example app",
		Long:  "Chainbridge Core Example app",
	}
)

func init() {
	// rootCMD.AddCommand(evmClient.CLI()) // Example of how CLI should be registered
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root cmd")
	}
}
