# Chainbridge example
<a href="https://discord.gg/ykXsJKfhgq">
  <img alt="discord" src="https://img.shields.io/discord/593655374469660673?label=Discord&logo=discord&style=flat" />
</a>

Chainbridge example is the repository that show an example of running Bridge with chainbridge-core [framework](https://github.com/ChainSafe/chainbridge-core).

1. [Installation](#installation)
2. [Modules](#modules)
3. [Global Flags](#global-flags)
4. [EVM-CLI](#evm-cli)
5. [Celo-CLI](#celo-cli)

### Installation
Refer to [installation](https://github.com/ChainSafe/chainbridge-docs/blob/develop/docs/installation.md) guide for assistance in installing.

### Modules

The chainbridge-core-example currently supports two modules:
1. [EVM-CLI](#evm-cli)
2. [Celo-CLI](celo-cli)

### Global Flags

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

### EVM-CLI
This module provides instruction for communicating with EVM-compatible chains.

```bash
Usage:
   evm-cli [command]

Available Commands:
  accounts    Account instructions
  admin       Admin-related instructions
  bridge      Bridge-related instructions
  deploy      Deploy smart contracts
  erc20       ERC20-related instructions
  erc721      ERC721-related instructions
  utils       Utils-related instructions

Flags:
  -h, --help   help for evm-cli
```

#### Accounts
Account instructions, allowing us to generate keypairs or import existing keypairs for use.

```bash
Usage:
   evm-cli accounts [command]

Available Commands:
  generate    Generate bridge keystore (Secp256k1)
  import      Import bridge keystore

Flags:
  -h, --help   help for accounts
```

*generate*
The generate subcommand is used to generate the bridge keystore. If no options are specified, a Secp256k1 key will be made.

```bash
Usage:
   evm-cli accounts generate [flags]

Flags:
  -h, --help   help for generate
```

*import*
The import subcommand is used to import a keystore for the bridge.

```bash
Usage:
   evm-cli accounts import [flags]

Flags:
  -h, --help              help for import
      --password string   password to encrypt with
```

#### Admin
Admin-related instructions.

```bash
Usage:
   evm-cli admin [command]

Available Commands:
  add-admin      Add a new admin
  add-relayer    Add a new relayer
  is-relayer     Check if an address is registered as a relayer
  pause          Pause deposits and proposals
  remove-admin   Remove an existing admin
  remove-relayer Remove a relayer
  set-fee        Set a new fee for deposits
  set-threshold  Set a new relayer vote threshold
  unpause        Unpause deposits and proposals
  withdraw       Withdraw tokens from the handler contract

Flags:
  -h, --help   help for admin
```

*add-admin*
Add a new admin.

```bash
Usage:
   evm-cli admin add-admin [flags]

Flags:
      --admin string    address to add
      --bridge string   bridge contract address
  -h, --help            help for add-admin
```

*add-relayer*
Add a new relayer.

```bash
Usage:
   evm-cli admin add-relayer [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for add-relayer
      --relayer string   address to add
```

*is-relayer*
Check if an address is registered as a relayer.

```bash
Usage:
   evm-cli admin is-relayer [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for is-relayer
      --relayer string   address to check
```

*pause*
Pause deposits and proposals,

```bash
Usage:
   evm-cli admin pause [flags]

Flags:
      --bridge string   bridge contract address
  -h, --help            help for pause

```

*remove-admin*
Remove an existing admin.

```bash
Usage:
   evm-cli admin remove-admin [flags]

Flags:
      --admin string    address to remove
      --bridge string   bridge contract address
  -h, --help            help for remove-admin
```

*remove-relayer*
Remove a relayer.

```bash
Usage:
   evm-cli admin remove-relayer [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for remove-relayer
      --relayer string   address to remove
```

*set-fee*
Set a new fee for deposits.

```bash
Usage:
   evm-cli admin set-fee [flags]

Flags:
      --bridge string   bridge contract address
      --fee string      New fee (in ether)
  -h, --help            help for set-fee
```

*set-threshold*
et a new relayer vote threshold.

```bash
Usage:
   evm-cli admin set-threshold [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for set-threshold
      --threshold uint   new relayer threshold
```

*unpause*
Unpause deposits and proposals.

```bash
Usage:
   evm-cli admin unpause [flags]

Flags:
      --bridge string   bridge contract address
  -h, --help            help for unpause
```

*withdraw*
Withdraw tokens from the handler contract.

```bash
Usage:
   evm-cli admin withdraw [flags]

Flags:
      --amount string      token amount to withdraw. Should be set or ID or amount if both set error will occur
      --bridge string      bridge contract address
      --decimals uint      ERC20 token decimals
      --handler string     handler contract address
  -h, --help               help for withdraw
      --id string          token ID to withdraw. Should be set or ID or amount if both set error will occur
      --recipient string   address to withdraw to
      --token string       ERC20 or ERC721 token contract address
```

#### Bridge
Bridge-related instructions.

```bash
Usage:
   evm-cli bridge [command]

Available Commands:
  cancel-proposal           Cancel an expired proposal
  query-proposal            Query an inbound proposal
  query-resource            Query the contract address
  register-generic-resource Register a generic resource ID
  register-resource         Register a resource ID
  set-burn                  Set a token contract as mintable/burnable

Flags:
  -h, --help   help for bridge
```

*cancel-proposal*
Cancel an expired proposal.

```bash
Usage:
   evm-cli bridge cancel-proposal [flags]

Flags:
      --bridge string       bridge contract address
      --chainId uint        chain ID of proposal to cancel
      --dataHash string     hash of proposal metadata
      --depositNonce uint   deposit nonce of proposal to cancel
  -h, --help                help for cancel-proposal
```

*query-proposal*
Query an inbound proposal.

```bash
Usage:
   evm-cli bridge query-proposal [flags]

Flags:
      --bridge string       bridge contract address
      --chainId uint        source chain ID of proposal
      --dataHash string     hash of proposal metadata
      --depositNonce uint   deposit nonce of proposal
  -h, --help                help for query-proposal
```

*query-resource*
Query the contract address with the provided resource ID for a specific handler contract.

```bash
Usage:
   evm-cli bridge query-resource [flags]

Flags:
      --handler string      handler contract address
  -h, --help                help for query-resource
      --resourceId string   resource ID to query
```

*register-generic-resource*
Register a resource ID with a contract address for a generic handler.

```bash
Usage:
   evm-cli bridge register-generic-resource [flags]

Flags:
      --bridge string       bridge contract address
      --deposit string      deposit function signature (default "0x00000000")
      --execute string      execute proposal function signature (default "0x00000000")
      --handler string      handler contract address
      --hash                treat signature inputs as function prototype strings, hash and take the first 4 bytes
  -h, --help                help for register-generic-resource
      --resourceId string   resource ID to query
      --target string       contract address to be registered
```
*register-resource*
Register a resource ID

```bash
Usage:
   evm-cli bridge register-resource [flags]

Flags:
      --bridge string       bridge contract address
      --handler string      handler contract address
  -h, --help                help for register-resource
      --resourceId string   resource ID to be registered
      --target string       contract address to be registered
```

*set-burn*
Set a token contract as mintable/burnable

```bash
Usage:
   evm-cli bridge set-burn [flags]

Flags:
      --bridge string          bridge contract address
      --handler string         ERC20 handler contract address
  -h, --help                   help for set-burn
      --tokenContract string   token contract to be registered
```

#### Deploy
Deploy smart contracts.

Used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags

```bash
Usage:
   evm-cli deploy [flags]

Flags:
      --all                     deploy all
      --bridge                  deploy bridge
      --bridgeAddress string    bridge contract address. Should be provided if handlers are deployed separately
      --chainId string          chain ID for the instance (default "1")
      --erc20                   deploy ERC20
      --erc20Handler            deploy ERC20 handler
      --erc20Name string        ERC20 contract name
      --erc20Symbol string      ERC20 contract symbol
      --erc721                  deploy ERC721
      --fee string              fee to be taken when making a deposit (in ETH, decimas are allowed) (default "0")
  -h, --help                    help for deploy
      --relayerThreshold uint   number of votes required for a proposal to pass (default 1)
      --relayers strings        list of initial relayers
```

#### ERC20
ERC20-related instructions.

```bash
Usage:
   evm-cli erc20 [command]

Available Commands:
  add-minter  Add a minter to an Erc20 mintable contract
  allowance   Get the allowance of a spender for an address
  approve     Approve tokens in an ERC20 contract for transfer
  balance     Query balance of an account in an ERC20 contract
  deposit     Initiate a transfer of ERC20 tokens
  mint        Mint tokens on an ERC20 mintable contract

Flags:
  -h, --help   help for erc20
```

*add-minter*
Add a minter to an Erc20 mintable contract.

```bash
Usage:
   evm-cli erc20 add-minter [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for add-minter
      --minter string         address of minter

```

*allowance*
Get the allowance of a spender for an address.

```bash
Usage:
   evm-cli erc20 allowance [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for allowance
      --owner string          address of token owner
      --spender string        address of spender
```

*approve*
Approve tokens in an ERC20 contract for transfer.

```bash
Usage:
   evm-cli erc20 approve [flags]

Flags:
      --amount string         amount to grant allowance
      --decimals uint         ERC20 token decimals (default 18)
      --erc20address string   ERC20 contract address
  -h, --help                  help for approve
      --recipient string      address of recipient
```

*balance*
Query balance of an account in an ERC20 contract.

```bash
Usage:
   evm-cli erc20 balance [flags]

Flags:
      --accountAddress string   address to receive balance of
      --erc20Address string     ERC20 contract address
  -h, --help                    help for balance
```

*deposit*
Initiate a transfer of ERC20 tokens.

```bash
Usage:
   evm-cli erc20 deposit [flags]

Flags:
      --amount string       amount to deposit
      --bridge string       address of bridge contract
      --decimals uint       ERC20 token decimals
      --destId string       destination chain ID
  -h, --help                help for deposit
      --recipient string    address of recipient
      --resourceId string   resource ID for transfer
```

*mint*
Mint tokens on an ERC20 mintable contract.

```bash
Usage:
   evm-cli erc20 mint [flags]

Flags:
      --amount string         amount to mint fee (in ETH)
      --decimal uint          ERC20 token decimals (default 18)
      --dstAddress string     Where tokens should be minted. Defaults to TX sender
      --erc20Address string   ERC20 contract address
  -h, --help                  help for mint
```

#### ERC721
ERC721-related instructions.

*add-minter*
Add a minter to an ERC721 mintable contract.

```bash
Usage:
   evm-cli erc721 add-minter [flags]

Flags:
      --erc721Address string   ERC721 contract address
  -h, --help                   help for add-minter
      --minter string          address of minter
```

#### Utils
Utils-related instructions.
*Useful for debugging*

```bash
Usage:
   evm-cli utils [command]

Available Commands:
  hashList    List tx hashes
  simulate    Simulate transaction invocation

Flags:
  -h, --help   help for utils
```

*hashlist*
List tx hashes.

```bash
Usage:
   evm-cli utils hashList [flags]

Flags:
      --blockNumber string   block number
  -h, --help                 help for hashList
```

*simulate*
Replay a failed transaction by simulating invocation; not state-altering

```bash
Usage:
   evm-cli utils simulate [flags]

Flags:
      --blockNumber string   block number
      --fromAddress string   address of sender
  -h, --help                 help for simulate
      --txHash string        transaction hash
```

### Celo-CLI
Though Celo is an EVM-compatible chain, it deviates in its implementation of the original Ethereum specifications, and therefore is deserving of its own separate module.

```bash
Root command for starting Celo CLI

Usage:
   celo-cli [command]

Available Commands:
  bridge      Bridge-related instructions
  deploy      Deploy smart contracts
  erc20       erc20-related instructions

Flags:
  -h, --help   help for celo-cli
```

#### Bridge
Bridge-related instructions.

```bash
Usage:
   celo-cli bridge [command]

Available Commands:
  register-resource Register a resource ID
  set-burn          Set a token contract as mintable/burnable

Flags:
  -h, --help   help for bridge
```

*register-resource*
Register a resource ID with a contract address for a handler

```bash
Usage:
   celo-cli bridge register-resource [flags]

Flags:
      --bridge string       bridge contract address
      --handler string      handler contract address
  -h, --help                help for register-resource
      --resourceId string   resource ID to be registered
      --target string       contract address to be registered
```

*set-burn*
Set a token contract as mintable/burnable in a handler

```bash
Usage:
   celo-cli bridge set-burn [flags]

Flags:
      --bridge string          bridge contract address
      --handler string         ERC20 handler contract address
  -h, --help                   help for set-burn
      --tokenContract string   token contract to be registered
```

#### Deploy
Deploy smart contracts.

This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags.

```bash
Usage:
   celo-cli deploy [flags]

Flags:
      --all                     deploy all
      --bridge                  deploy bridge
      --bridgeAddress string    bridge contract address. Should be provided if handlers are deployed separately
      --chainId string          chain ID for the instance (default "1")
      --erc20                   deploy ERC20
      --erc20Handler            deploy ERC20 handler
      --erc20Name string        ERC20 contract name
      --erc20Symbol string      ERC20 contract symbol
      --erc721                  deploy ERC721
      --fee string              fee to be taken when making a deposit (in ETH, decimas are allowed) (default "0")
  -h, --help                    help for deploy
      --relayerThreshold uint   number of votes required for a proposal to pass (default 1)
      --relayers strings        list of initial relayers
```

#### ERC20
erc20-related instructions

```bash
Usage:
   celo-cli erc20 [command]

Available Commands:
  add-minter  Add a minter to an Erc20 mintable contract
  allowance   Set a token contract as mintable/burnable
  approve     Approve tokens in an ERC20 contract for transfer
  balance     Query balance of an account in an ERC20 contract
  deposit     Initiate a transfer of ERC20 tokens
  mint        Mint tokens on an ERC20 mintable contract

Flags:
  -h, --help   help for erc20
```

*add-minter*
Add a minter to an Erc20 mintable contract.

```bash
Usage:
   celo-cli erc20 add-minter [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for add-minter
      --minter string         address of minter
```

*allowance*
Set a token contract as mintable/burnable in a handler.

```bash
Usage:
   celo-cli erc20 allowance [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for allowance
      --owner string          address of token owner
      --spender string        address of spender
```

*approve*
Approve tokens in an ERC20 contract for transfer.

```bash
Usage:
   celo-cli erc20 approve [flags]

Flags:
      --amount string         amount to grant allowance
      --decimals uint         ERC20 token decimals (default 18)
      --erc20address string   ERC20 contract address
  -h, --help                  help for approve
      --recipient string      address of recipient
```

*balance*
Query balance of an account in an ERC20 contract.

```bash
Usage:
   celo-cli erc20 balance [flags]

Flags:
      --accountAddress string   address to receive balance of
      --erc20Address string     ERC20 contract address
  -h, --help                    help for balance
```

*deposit*
Initiate a transfer of ERC20 tokens.

```bash
Usage:
   celo-cli erc20 deposit [flags]

Flags:
      --amount string       amount to deposit
      --bridge string       address of bridge contract
      --decimals uint       ERC20 token decimals
      --destId string       destination chain ID
  -h, --help                help for deposit
      --recipient string    address of recipient
      --resourceId string   resource ID for transfer
```

*mint*
Mint tokens on an ERC20 mintable contract.

```bash
Usage:
   celo-cli erc20 mint [flags]

Flags:
      --amount string         amount to mint fee (in ETH)
      --decimal uint          ERC20 token decimals (default 18)
      --dstAddress string     Where tokens should be minted. Defaults to TX sender
      --erc20Address string   ERC20 contract address
  -h, --help                  help for mint
```

#### Run
Run the actual chainbridge relayer.

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

*run*
Running the relayer with the following flags:
1. Path to our chain configurations: JSON config file(s) stored within a directory called configs.
2. Path to our relayer's keystore: an ethereum keypair used for signing transactions.
3. Path to our blockstore: used to record the last block the relayer processed, allowing the relayer to pick up where it left off.

```bash
run --config configs --keystore keys --blockstore blockstore
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