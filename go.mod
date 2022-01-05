module github.com/ChainSafe/chainbridge-core-example

go 1.15

//replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core
//
// replace github.com/ChainSafe/chainbridge-celo-module => ../chainbridge-celo-module

require (
	github.com/ChainSafe/chainbridge-celo-module v0.0.0-20211125083019-55c33ab406b4
	github.com/ChainSafe/chainbridge-core v0.0.0-20211126125725-1525c043927c
	github.com/rs/zerolog v1.26.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
)
