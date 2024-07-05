#!/bin/bash

set -eux

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
NODE_HOME=${NODE_HOME:-"$PWD/.oraid"}
ARGS="--from $USER --chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block --home $NODE_HOME"

# prepare a new contract for gasless
fee_params=$(oraid query tokenfactory params --output json | jq '.params.denom_creation_fee[0].denom')
if ! [[ $fee_params =~ "orai" ]] ; then
   echo "The tokenfactory fee params is not orai"; exit 1
fi

# try creating a new denom
denom_name="usdt"
oraid tx tokenfactory create-denom $denom_name $ARGS

# try querying list denoms afterwards
user_address=$(oraid keys show $USER --home $NODE_HOME --keyring-backend test -a)
first_denom=$(oraid query tokenfactory denoms-from-creator $user_address --output json | jq '.denoms[0]' | tr -d '"')
echo "first denom: $first_denom"

if ! [[ $first_denom =~ "factory/$user_address/$denom_name" ]] ; then
   echo "The tokenfactory denom does not match the created denom"; exit 1
fi

admin=$(oraid query tokenfactory denom-authority-metadata $first_denom --output json | jq '.authority_metadata.admin')
echo "admin: $admin"

if ! [[ $admin =~ $user_address ]] ; then
   echo "The tokenfactory admin does not match the creator"; exit 1
fi

# try mint
oraid tx tokenfactory mint 10$first_denom $ARGS

# query balance after mint
tokenfactory_balance=$(oraid query bank balances $user_address --denom=$first_denom --output json | jq '.amount | tonumber')
if [[ $tokenfactory_balance -ne 10 ]] ; then
   echo "The tokenfactory balance does not increase after mint"; exit 1
fi

# try burn
oraid tx tokenfactory burn 10$first_denom $ARGS
tokenfactory_balance=$(oraid query bank balances $user_address --denom=$first_denom --output json | jq '.amount | tonumber')
if [[ $tokenfactory_balance -ne 0 ]] ; then
   echo "The tokenfactory balance does not decrease after burn"; exit 1
fi

echo "Tokenfactory tests passed!"