module github.com/ChainSafe/chainbridge-core-example

go 1.16

replace (
	github.com/ChainSafe/chainbridge-core => ../chainbridge-core
	github.com/ChainSafe/chainbridge-substrate-module => ../chainbridge-substrate-module
)

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20210520113638-fb0ff8dc9606
	github.com/ChainSafe/chainbridge-substrate-module v0.0.0-20210521092722-d2ee3d9d63cc
	github.com/centrifuge/go-substrate-rpc-client v2.0.0+incompatible
	github.com/ethereum/go-ethereum v1.10.4
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
)
