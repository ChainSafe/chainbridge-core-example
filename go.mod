module github.com/ChainSafe/chainbridge-core-example

go 1.15

replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core

require (
	github.com/ChainSafe/chainbridge-celo-module v0.0.0-20210805091055-68700289e998
	github.com/ChainSafe/chainbridge-core v0.0.0-20210804101411-6e069815c674
	github.com/ethereum/go-ethereum v1.10.6
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
)
