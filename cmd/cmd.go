package cmd

import (
	"github.com/ChainSafe/chainbridge-core-example/example"
	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCMD = &cobra.Command{
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
	//root.AddCommand(evmClient.CLI()) // Example of how CLI should be registered

	// TODO: possibly see if these can be bound in core
	rootCMD.Flags().String(config.KeystoreFlagName, "./keys", "Path to keystore directory")
	viper.BindPFlag(config.KeystoreFlagName, rootCMD.Flags().Lookup(config.KeystoreFlagName))

	rootCMD.Flags().String(config.BlockstoreFlagName, "./lvldbdata", "Specify path for blockstore")
	viper.BindPFlag(config.BlockstoreFlagName, rootCMD.Flags().Lookup(config.BlockstoreFlagName))

	rootCMD.Flags().Bool(config.FreshStartFlagName, false, "Disables loading from blockstore at start. Opts will still be used if specified.")
	viper.BindPFlag(config.FreshStartFlagName, rootCMD.Flags().Lookup(config.FreshStartFlagName))

	rootCMD.Flags().Bool(config.LatestBlockFlagName, false, "Overrides blockstore and start block, starts from latest block")
	viper.BindPFlag(config.LatestBlockFlagName, rootCMD.Flags().Lookup(config.LatestBlockFlagName))

	rootCMD.Flags().String(config.TestKeyFlagName, "", "Applies a predetermined test keystore to the chains.")
	viper.BindPFlag(config.TestKeyFlagName, rootCMD.Flags().Lookup(config.TestKeyFlagName))

}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root cmd")
	}
}
