#!/bin/bash

set -eu

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
NODE_HOME=${NODE_HOME:-"$PWD/.oraid"}
ARGS="--from $USER --chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block --home $NODE_HOME"
docker_command="docker-compose -f $PWD/docker-compose-e2e-upgrade.yml exec"
validator1_command="$docker_command validator1 bash -c"

# prepare a new contract for gasless
fee_params=`$validator1_command "oraid query tokenfactory params --output json | jq '.params.denom_creation_fee[0].denom'"`
if ! [[ $fee_params =~ "orai" ]] ; then
   echo "The tokenfactory fee params is not orai"; exit 1
fi

# try creating a new denom
denom_name="usdt"
$validator1_command "oraid tx tokenfactory create-denom $denom_name $ARGS"

# try querying list denoms afterwards
user_address_result=`$validator1_command "oraid keys show $USER --home $NODE_HOME --keyring-backend test --output json"`
user_address=$(echo $user_address_result | jq '.address' | tr -d '"')
first_denom_before_trim=`$validator1_command "oraid query tokenfactory denoms-from-creator $user_address --output json"`
first_denom=$(echo $first_denom_before_trim | jq '.denoms[0]' | tr -d '"')
echo "first denom: $first_denom"

if ! [[ $first_denom =~ "factory/$user_address/$denom_name" ]] ; then
   echo "The tokenfactory denom does not match the created denom"; exit 1
fi

admin=`$validator1_command "oraid query tokenfactory denom-authority-metadata $first_denom --output json | jq '.authority_metadata.admin'"`
echo "admin: $admin"

if ! [[ $admin =~ $user_address ]] ; then
   echo "The tokenfactory admin does not match the creator"; exit 1
fi

# try mint
$validator1_command "oraid tx tokenfactory mint 10$first_denom $ARGS"

# query balance after mint
tokenfactory_balance_result=`$validator1_command "oraid query bank balances $user_address --denom=$first_denom --output json"`
tokenfactory_balance=$(echo $tokenfactory_balance_result | jq '.amount | tonumber')
echo "token factory balance: $tokenfactory_balance"
if [[ $tokenfactory_balance -ne 10 ]] ; then
   echo "The tokenfactory balance does not increase after mint"; exit 1
fi

# try burn
$validator1_command "oraid tx tokenfactory burn 10$first_denom $ARGS"
tokenfactory_balance_result=`$validator1_command "oraid query bank balances $user_address --denom=$first_denom --output json"`
tokenfactory_balance=$(echo $tokenfactory_balance_result | jq '.amount | tonumber')
if [[ $tokenfactory_balance -ne 0 ]] ; then
   echo "The tokenfactory balance does not decrease after burn"; exit 1
fi

echo "Tokenfactory tests passed!"