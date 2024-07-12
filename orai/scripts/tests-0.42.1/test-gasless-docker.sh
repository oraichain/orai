#!/bin/bash

set -eu

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
WASM_PATH=${WASM_PATH:-"$PWD/scripts/wasm_file/swapmap.wasm"}
EXECUTE_MSG=${EXECUTE_MSG:-'{"ping":{}}'}
NODE_HOME=${NODE_HOME:-"$PWD/.oraid"}
ARGS="--from $USER --chain-id $CHAIN_ID -y --keyring-backend test --gas 20000000 --gas-adjustment 1.5 -b block --home $NODE_HOME"
docker_command="docker-compose -f $PWD/docker-compose-e2e-upgrade.yml exec"
validator1_command="$docker_command validator1 bash -c"
HIDE_LOGS="/dev/null"

# prepare a new contract for gasless
store_ret=`$validator1_command "oraid tx wasm store $WASM_PATH --from validator1 --home $NODE_HOME $ARGS --output json"`
echo "store ret: $store_ret"
code_id=$(echo $store_ret | jq -r '.logs[0].events[1].attributes[] | select(.key | contains("code_id")).value')
$validator1_command "oraid tx wasm instantiate $code_id \"{}\" --label 'testing' --from validator1 --home $NODE_HOME -b block --no-admin $ARGS > $HIDE_LOGS"
contract_address_res=`$validator1_command "oraid query wasm list-contract-by-code $code_id --output json | jq -r '.contracts[0]'"`
contract_address=$(echo "$contract_address_res" | tr -d -c '[:alnum:]') # remove all special characters because of the command's result

# set gasless proposal
$validator1_command "oraid tx gov submit-proposal set-gasless $contract_address --title 'gasless' --description 'gasless' --deposit 10000000orai $ARGS > $HIDE_LOGS"
proposal_id_result=`$validator1_command "oraid query gov proposals --reverse --output json"`
proposal_id=$(echo $proposal_id_result | jq '.proposals[0].proposal_id | tonumber')

$validator1_command "oraid tx gov vote $proposal_id yes $ARGS > $HIDE_LOGS"

# wait til proposal passes
sleep 6
proposal_status_result=`$validator1_command "oraid query gov proposal $proposal_id --output json"`
proposal_status=$(echo $proposal_status_result | jq '.status')
if ! [[ $proposal_status =~ "PROPOSAL_STATUS_PASSED" ]] ; then
   echo "The proposal has not passed yet"; exit 1
fi