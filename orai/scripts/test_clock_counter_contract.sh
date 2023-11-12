#!/bin/bash
# Before running this script, you must setup local network:
# sh $PWD/scripts/multinode-local-testnet.sh

WASM_PATH=${WASM_PATH:-"$PWD/scripts/wasm_file/cw-clock-example.wasm"}
ARGS="--chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block"
VALIDATOR_HOME=${VALIDATOR_HOME:-"$HOME/.oraid/validator1"}
QUERY_MSG=${QUERY_MSG:-'{"get_config":{}}'}

CONTRACT_GAS_LIMIT=${CONTRACT_GAS_LIMIT:-"123000000"}
TITLE=${TITLE:-"add contract to clock module"}
INITIAL_DEPOSIT=${INITIAL_DEPOSIT:-"10000000orai"}
DESCRIPTION=${DESCRIPTION:-"add cw-clock contract to clock module"}

store_ret=$(oraid tx wasm store $WASM_PATH --from validator1 --home $VALIDATOR_HOME $ARGS --output json)
code_id=$(echo $store_ret | jq -r '.logs[0].events[1].attributes[] | select(.key | contains("code_id")).value')
oraid tx wasm instantiate $code_id '{}' --label 'cw clock contract' --from validator1 --home $VALIDATOR_HOME -b block --admin $(oraid keys show validator1 --keyring-backend test --home $VALIDATOR_HOME -a) $ARGS
CONTRACT_ADDRESS=$(oraid query wasm list-contract-by-code $code_id --output json | jq -r '.contracts[0]')
echo "cw-clock contract address: $CONTRACT_ADDRESS"

oraid tx clock add-contract $CONTRACT_ADDRESS $CONTRACT_GAS_LIMIT $TITLE $INITIAL_DEPOSIT $DESCRIPTION --from validator1 --home $VALIDATOR_HOME $ARGS

oraid tx gov vote 2 yes --from validator3 --home "$HOME/.oraid/validator3" $ARGS && oraid tx gov vote 1 yes --from validator2 --home "$HOME/.oraid/validator2" $ARGS && oraid tx gov vote 1 yes --from validator1 --home $VALIDATOR_HOME $ARGS

# sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."
sleep 96

# Query the counter
counter_before=$(oraid query wasm contract-state smart $contract_address $QUERY_MSG --node "tcp://localhost:26657" --output json | jq -r '.data.val | tonumber')
sleep 7
echo "cw-clock counter_before: $counter_before"

counter_after=$(oraid query wasm contract-state smart $contract_address $QUERY_MSG --node "tcp://localhost:26657" --output json | jq -r '.data.val | tonumber')
sleep 7
echo "cw-clock counter_after: $counter_after"

if [ $counter_after -gt $counter_before ]
then
echo "Test Passed"
else
echo "Test Failed"
fi
