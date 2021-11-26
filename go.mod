module github.com/ChainSafe/chainbridge-core-example

go 1.15

//replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core
//
// replace github.com/ChainSafe/chainbridge-celo-module => ../chainbridge-celo-module

require (
	github.com/ChainSafe/chainbridge-celo-module v0.0.0-20211123153704-07133071e1ce
	github.com/ChainSafe/chainbridge-core v0.0.0-20211126100456-46d385089a1f
	github.com/rs/zerolog v1.26.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
)
