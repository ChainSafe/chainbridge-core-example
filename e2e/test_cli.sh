#!/usr/bin/env bash
# Copyright 2021 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

CMD=./chainbridge-core-example

# Default values
GAS_LIMIT=6721975
GAS_PRICE=20000000000

################
# USER CONFIGURATIONS
################

# Node WS endpoint
NODE_URL=""

# Address to receive transfers
RECIPIENT_ADDRESS=""

# Addresses of the relayers
RELAYER_ADDRESS1=""
RELAYER_ADDRESS2=""
RELAYER_ADDRESS3=""

################
# OUTPUTS
################

# Deployed bridge contract
BRIDGE_ADDRESS=""

# Deployed ERC20 contracts
ERC20_ADDRESS=""
ERC20_HANDLER=""

# ERC20 Resource ID
ERC20_RESOURCE_ID=""

# Deployed ERC721 contracts
ERC721_HANDLER=""
ERC721_CONTRACT=""

# ERC721 Resource ID
ERC721_RESOURCE_ID=""

# Deployed Generic contracts
GENERIC_HANDLER=""
GENERIC_RESOURCE_ID=""

set -eux

## deploy
# deploy bridge
{
    $CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --bridge --url $NODE_URL --relayerThreshold 1 --relayers="$RELAYER_ADDRESS1,$RELAYER_ADDRESS2,$RELAYER_ADDRESS3"
} || {
	exit
}

# deploy erc20
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --erc20 --url $NODE_URL --erc20Symbol "TKN" --erc20Name "token"

# deploy erc20 handler
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --erc20Handler --url $NODE_URL --bridgeAddress $BRIDGE_ADDRESS

## bridge
# register resource
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge register-resource --url $NODE_URL --resourceId $ERC20_RESOURCE_ID --target $ERC20_HANDLER

# register contract as mintable/burnable
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge set-burn --url $NODE_URL --bridge $BRIDGE_ADDRESS --handler $ERC20_HANDLER --tokenContract $ERC20_ADDRESS

## erc20
# register handler as mintable
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 add-minter --url $NODE_URL --minter $ERC20_HANDLER --erc20Address $ERC20_ADDRESS

# approve erc20
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 approve --url $NODE_URL --erc20Address $ERC20_ADDRESS --recipient $RECIPIENT_ADDRESS --amount "1.11"  --decimals 2

# mint erc20
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 mint --url $NODE_URL --amount "100" --erc20Address $ERC20_ADDRESS