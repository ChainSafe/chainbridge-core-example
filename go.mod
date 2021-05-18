module github.com/ChainSafe/chainbridge-example

go 1.16

replace github.com/ChainSafe/chainbridgev2 => /var/www/ChainSafe/rnd/chainbridgev2

replace github.com/ChainSafe/chainbridge-eth-module => /var/www/ChainSafe/rnd/chainbridge-eth-module

replace github.com/ChainSafe/chainbridge-substrate-module => /var/www/ChainSafe/rnd/chainbridge-substrate-module

require (
	github.com/ChainSafe/chainbridge-eth-module v0.0.0-00010101000000-000000000000
	github.com/ChainSafe/chainbridge-substrate-module v0.0.0-00010101000000-000000000000
	github.com/ChainSafe/chainbridgev2 v0.0.0-00010101000000-000000000000
	github.com/centrifuge/go-substrate-rpc-client v2.0.0+incompatible
	github.com/ethereum/go-ethereum v1.10.3
	github.com/rs/zerolog v1.21.0
	github.com/spf13/cobra v1.1.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)
