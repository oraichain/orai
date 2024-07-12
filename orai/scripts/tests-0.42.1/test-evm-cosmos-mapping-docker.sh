#!/bin/bash

set -eux

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
NODE_HOME=${NODE_HOME:-"$PWD/.oraid"}
ARGS="--from $USER --chain-id $CHAIN_ID -y --keyring-backend test --gas auto --gas-adjustment 1.5 -b block --home $NODE_HOME"
docker_command="docker-compose -f $PWD/docker-compose-e2e-upgrade.yml exec"
validator1_command="$docker_command validator1 bash -c"
HIDE_LOGS="/dev/null"

user_address_result=`$validator1_command "oraid keys show $USER --home $NODE_HOME --keyring-backend test --output json"`
user_address=$(echo $user_address_result | jq '.address' | tr -d '"')

user_pubkey_result=`$validator1_command "oraid keys show $USER --home $NODE_HOME --keyring-backend test -p"`
user_pubkey=$(echo $user_pubkey_result | jq '.key' | tr -d '"')
$validator1_command "oraid tx evm set-mapping-evm $user_pubkey $ARGS > $HIDE_LOGS"

expected_evm_address=`$validator1_command "oraid debug pubkey-simple $user_pubkey"`
actual_evm_address_result=`$validator1_command "oraid query evm mappedevm $user_address --output json"`
actual_evm_address=$(echo $actual_evm_address_result | jq '.evm_address' | tr -d '"')
if ! [[ $actual_evm_address =~ $expected_evm_address ]] ; then
   echo "The evm addresses dont match"; exit 1
fi

echo "EVM cosmos mapping docker tests passed!"