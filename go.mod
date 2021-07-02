module github.com/ChainSafe/chainbridge-core-example

go 1.16

replace github.com/ChainSafe/chainbridge-core => /private/var/www/ChainSafe/rnd/chainbridge-core

replace github.com/ChainSafe/chainbridge-celo-module => /private/var/www/ChainSafe/rnd/chainbridge-celo-module

require (
	github.com/ChainSafe/chainbridge-celo-module v0.0.0-00010101000000-000000000000
	github.com/ChainSafe/chainbridge-core v0.0.0-20210520113638-fb0ff8dc9606
	github.com/ChainSafe/chainbridge-substrate-module v0.0.0-20210521092722-d2ee3d9d63cc
	github.com/centrifuge/go-substrate-rpc-client v2.0.0+incompatible
	github.com/ethereum/go-ethereum v1.10.4
	github.com/rs/zerolog v1.23.0
	github.com/spf13/cobra v1.1.3
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
)
