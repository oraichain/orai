WASM_PATH=${WASM_PATH:-"$PWD/scripts/wasm_file/cw-clock-example.wasm"}
ARGS="--chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block"
VALIDATOR_HOME=${VALIDATOR_HOME:-"$HOME/.oraid/validator1"}

# setup local network
sh $PWD/scripts/multinode-local-testnet.sh

store_ret=$(oraid tx wasm store $WASM_PATH --from validator1 --home $VALIDATOR_HOME --chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block --output json)
code_id=$(echo $store_ret | jq -r '.logs[0].events[1].attributes[] | select(.key | contains("code_id")).value')

oraid tx wasm instantiate $code_id '{}' --label 'cw clock contract' --from validator1 --home $VALIDATOR_HOME -b block --admin $(oraid keys show validator1 --keyring-backend test --home $VALIDATOR_HOME -a) --chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block
contract_address=$(oraid query wasm list-contract-by-code $code_id --output json | jq -r '.contracts[0]')
echo "cw-clock contract address: $contract_address"

contract_gas_limit="123000000"
title="add contract to clock"
initial_deposit=200000orai
description="add contract to clock"

oraid tx clock add-contract $contract_address $contract_gas_limit $title $initial_deposit $description --from validator1 --chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 --home $VALIDATOR_HOME

oraid tx gov vote 1 yes --from validator1 --home "$HOME/.oraid/validator1" $ARGS && oraid tx gov vote 1 yes --from validator2 --home "$HOME/.oraid/validator2" $ARGS

# sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."
sleep 12

