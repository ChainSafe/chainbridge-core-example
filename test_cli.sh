#!/usr/bin/env bash
# Copyright 2021 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

CMD=./chainbridge-core-example

KALEIDO_NODE_URL="wss://u0oo16l2hl:WQgUMKVenv9mefKXW57fV-TzSWhFz06YILiXlHGykd4@u0ize6fsgo-u0np93mhu4-wss.us0-aws.kaleido.io"

BRIDGE_ADDRESS="0xD83cA9C4EEa4B33F110d9F9e5E462D4111335090"
ERC20_ADDRESS="0xa0a041d451c3ae2Cb32C0a8c9fD054D6bf281D02"
ERC20_HANDLER="0x42e529d0249410fBf3b0690Ca2f100a83fB78F5F"
RESOURCE_ID="000000000000000000000021605f71845f372A9ed84253d2D024B7B10999f400"

GAS_LIMIT=6721975
GAS_PRICE=20000000000

BRIDGE2_ADDRESS=""
ERC721_HANDLER=""
ERC721_RESOURCE_ID="0000000000000000000000d7E33e1bbf65dC001A0Eb1552613106CD7e40C3100"
ERC721_CONTRACT=""

GENERIC_HANDLER=""
GENERIC_RESOURCE_ID="0000000000000000000000106C24dc2D480b5559C9E0e97bAaDf0750d9F0B800"

set -eux

# deploy bridge
$CMD evm-cli  --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --bridge --url $KALEIDO_NODE_URL

# deploy erc20 handler
$CMD evm-cli  --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --erc20Handler --url $KALEIDO_NODE_URL --bridgeAddress $BRIDGE_ADDRESS

# register resource
$CMD evm-cli  --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge register-resource --url $KALEIDO_NODE_URL --resourceId $RESOURCE_ID --target $ERC20_HANDLER

# register contract as mintable/burnable
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge set-burn --url $KALEIDO_NODE_URL --bridge $BRIDGE_ADDRESS --handler $ERC20_HANDLER --tokenContract $ERC20_ADDRESS

# register handler as mintable
$CMD evm-cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 add-minter --url $KALEIDO_NODE_URL --minter $ERC20_HANDLER
