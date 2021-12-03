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
This section will detail how to properly configure `chainbridge-core` and get started with a new, more modular chainbridge.

### Configuring chains
Chains are configured as an array under `chains` tag inside provided config. Inside
[example](./example/app.go) we then go through all chain configs and create an
instance of `RelayedChain` via `SetupDefaultEVMChain` that registers handlers
for provided handler addresses in config (currently supported: erc20, erc721, generic).


### Run
Run the chainbridge relayer.

```bash
Run example app

Usage:
   run [flags]

Flags:
      --blockstore string   Specify path for blockstore (default "./lvldbdata")
      --config string       Path to JSON configuration file
      --fresh               Disables loading from blockstore at start. Opts will still be used if specified. (default: false)
  -h, --help                help for run
      --keystore string     Path to keystore directory (default "./keys")
      --latest              Overrides blockstore and start block, starts from latest block (default: false)
      --testkey string      Applies a predetermined test keystore to the chains.
```

Running the relayer with the following flags:
1. Path to config: JSON config file path
2. Path to our relayer's keystore: an ethereum keypair used for signing transactions.
3. Path to our blockstore: used to record the last block the relayer processed, allowing the relayer to pick up where it left off.

_example:_
```bash
./chainbridge-core-example run --config config.json --keystore keys --blockstore blockstore
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
