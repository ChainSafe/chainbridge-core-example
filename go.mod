module github.com/ChainSafe/chainbridge-core-example

go 1.16

replace (
	github.com/ChainSafe/chainbridge-core v0.0.0-20210531163541-ba144909335d => ../chainbridge-core
	github.com/ChainSafe/chainbridge-eth-module v0.0.0-20210521100422-24b8768656d3 => ../chainbridge-eth-module
	github.com/ChainSafe/chainbridge-substrate-module v0.0.0-20210521092722-d2ee3d9d63cc => ../chainbridge-substrate-module
)

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20210531163541-ba144909335d
	github.com/ChainSafe/chainbridge-eth-module v0.0.0-20210521100422-24b8768656d3
	github.com/ChainSafe/chainbridge-substrate-module v0.0.0-20210521092722-d2ee3d9d63cc
	github.com/centrifuge/go-substrate-rpc-client v2.0.0+incompatible
	github.com/ethereum/go-ethereum v1.10.3
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/rs/zerolog v1.22.0
	github.com/spf13/cobra v1.1.1
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)
