# Chainbridge example
<a href="https://discord.gg/ykXsJKfhgq">
  <img alt="discord" src="https://img.shields.io/discord/593655374469660673?label=Discord&logo=discord&style=flat" />
</a>

ChainBridge example is the repository that show an example of running Bridge with chainbridge-core [framework](https://github.com/ChainSafe/chainbridge-core).

1. [Installation](#installation)
2. [Global Flags](#global-flags)
3. [Usage](#usage)

## Installation
Refer to [installation](https://github.com/ChainSafe/chainbridge-docs/blob/develop/docs/installation.md) guide for assistance in installing.

## Global Flags

Though the EVM-CLI and Celo-CLI may differ in their implementation of subcommands, both modules share global flag values.

These global flags are so-called due to their presence in every sucommand, irrespective of their relevance to the subcommand. 

In situations where the global flag values are not needed or are otherwise irrelevant, ie passing a `--gasLimit` flag when the `accounts generate` subcommand is invoked, the flag will be ignored.

```bash
Global Flags:
      --gasLimit uint               gasLimit used in transactions (default 6721975)
      --gasPrice uint               gasPrice used for transactions (default 20000000000)
      --jsonWallet string           Encrypted JSON wallet
      --jsonWalletPassword string   Password for encrypted JSON wallet
      --networkid uint              networkid
      --privateKey string           Private key to use
      --url string                  node url (default "ws://localhost:8545")
```

## Usage
This section will detail how to properly configure `example/app.go` and get started with a new, more modular chainbridge.

The main section of the app we will look at is within [example/app.go](example/app.go). It is within here that we will initialize and register all of the necessary components for adding a relayed chain for our bridge relayer to become aware of.

_The code comments will best describe what each step is doing._

```go
// example/app.go

// new chain setup: EVM (flavor)

// initialize a new client to communicate with the network, in this case Ethereum
evmClient := evmclient.NewEVMClient()

// configure the client with values from our `config.json`. includes:
// chain ID (unique)
// network RPC endpoint
// bridge contract address (deployed)
// ERC20 handler contract address (deployed)
// relayer wallet address
// gas limit
// maximum gas price
err = evmClient.Configurate(viper.GetString(config.ConfigFlagName), "config.json")
if err != nil {
      // typically panics are to be avoided, but here it is necessary to ensure that the client be properly configured before relaying
      panic(err)
}

// store the client configuration in a local variable for easier readability and usage
evmConfig := evmClient.GetConfig()

// instantiate a new event handler by passing in the deployed bridge address and newly initialized client
eventHandlerEVM := listener.NewETHEventHandler(common.HexToAddress(evmConfig.SharedEVMConfig.Bridge), evmClient)

// register the event handler by providing the ERC20 handler contract and event handler
eventHandlerEVM.RegisterEventHandler(evmConfig.SharedEVMConfig.Erc20Handler, listener.Erc20EventHandler)

// instantiate a new relayer listener by passing in the client, newly created and registered event handler as well as the deployed bridge address
evmListener := listener.NewEVMListener(evmClient, eventHandlerEVM, common.HexToAddress(evmConfig.SharedEVMConfig.Bridge))

// instantiate a new message handler by passing in the client and deployed bridge address
mhEVM := voter.NewEVMMessageHandler(evmClient, common.HexToAddress(evmConfig.SharedEVMConfig.Bridge))

// register the message handler by providing the ERC20 handler contract as well as the message handler
mhEVM.RegisterMessageHandler(common.HexToAddress(evmConfig.SharedEVMConfig.Erc20Handler), voter.ERC20MessageHandler)

// instantiate a new new voter by passing in the message handler, client
// and a generic transaction type that can serve as either a contact invocation or a contract deployment transaction
evmVoter := voter.NewVoter(mhEVM, evmClient, evmtransaction.NewTransaction)

// instantiate a new chain (in our case EVM) by passing in the listener, voter, database (leveldb), the chain ID (unique) as well as the client configuration
evmChain := evm.NewEVMChain(evmListener, evmVoter, db, *evmConfig.SharedEVMConfig.GeneralChainConfig.Id, &evmConfig.SharedEVMConfig)

// instantiate a new chainbridge relayer by passing it the chains wishing to be bridged
// Note: there can be more than 2 chains
r := relayer.NewRelayer([]relayer.RelayedChain{evmChain, otherChain1, otherChain2})
```

### Run
Run the chainbridge relayer.

```bash
Run example app

Usage:
   run [flags]

Flags:
      --blockstore string   Specify path for blockstore (default "./lvldbdata")
      --config string       Path to JSON configuration files directory (default ".")
      --fresh               Disables loading from blockstore at start. Opts will still be used if specified. (default: false)
  -h, --help                help for run
      --keystore string     Path to keystore directory (default "./keys")
      --latest              Overrides blockstore and start block, starts from latest block (default: false)
      --testkey string      Applies a predetermined test keystore to the chains.
```

Running the relayer with the following flags:
1. Path to our chain configurations: JSON config file(s) stored within a directory called configs.
2. Path to our relayer's keystore: an ethereum keypair used for signing transactions.
3. Path to our blockstore: used to record the last block the relayer processed, allowing the relayer to pick up where it left off.

_example:_
```bash
./chainbridge-core-example run --config configs --keystore keys --blockstore blockstore
```

# ChainSafe Security Policy

## Reporting a Security Bug

We take all security issues seriously, if you believe you have found a security issue within a ChainSafe
project please notify us immediately. If an issue is confirmed, we will take all necessary precautions
to ensure a statement and patch release is made in a timely manner.

Please email us a description of the flaw and any related information (e.g. reproduction steps, version) to
[security at chainsafe dot io](mailto:security@chainsafe.io).

## License

_GNU Lesser General Public License v3.0_