module github.com/ChainSafe/chainbridge-core-example

go 1.15

replace github.com/ChainSafe/chainbridge-optimism-module v0.0.0-20211027031914-8cbe669d45e4 => ../chainbridge-optimism-module

require (
	github.com/ChainSafe/chainbridge-celo-module v0.0.0-20210812101441-b6d7ad422a53
	github.com/ChainSafe/chainbridge-core v0.0.0-20211004151830-0d1addc4c2b1
	github.com/ChainSafe/chainbridge-optimism-module v0.0.0-20211027031914-8cbe669d45e4
	github.com/ethereum/go-ethereum v1.10.11
	github.com/rs/zerolog v1.25.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
)
