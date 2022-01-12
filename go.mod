module github.com/ChainSafe/chainbridge-core-example

go 1.15

//replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core
//
// replace github.com/ChainSafe/chainbridge-celo-module => ../chainbridge-celo-module

require (
	github.com/ChainSafe/chainbridge-celo-module v0.0.0-20220107113949-c7fd9c5f44ac
	github.com/ChainSafe/chainbridge-core v0.0.0-20220110124723-abb0bf918502
	github.com/rs/zerolog v1.26.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
)
