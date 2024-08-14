#!/bin/bash

set -eu

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
NODE_HOME=${NODE_HOME:-"$PWD/.oraid"}
ARGS="--from $USER --chain-id $CHAIN_ID -y --keyring-backend test --gas 20000000 -b block --home $NODE_HOME"
HIDE_LOGS="/dev/null"

# TODO: fix the content of evmutil.json to match the kava erc20 address and the tokenfactory created for regression testing
oraid tx gov submit-proposal param-change ./evmutil.json --from $USER --home $NODE_HOME $ARGS

# FIXME: get proposal value from result instead of hard-coding
oraid tx gov vote 2 yes --from $USER --keyring-backend test --chain-id testing --home $NODE_HOME -y -b block
oraid tx gov vote 2 yes --from validator2 --keyring-backend test --chain-id testing --home ~/.oraid/validator2 -y -b block

# the proposal takes about 8s
sleep 8s

params=$(oraid query evmutil params --home $NODE_HOME)