module github.com/ChainSafe/chainbridge-core-example

go 1.15

replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core

replace github.com/ChainSafe/chainbridge-celo-module => ../chainbridge-celo-module

require (
	github.com/ChainSafe/chainbridge-celo-module v0.0.0-20210702092144-957c6185d362
	github.com/ChainSafe/chainbridge-core v0.0.0-20210712095225-dd3876a066e4
	github.com/ethereum/go-ethereum v1.10.4
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
)
