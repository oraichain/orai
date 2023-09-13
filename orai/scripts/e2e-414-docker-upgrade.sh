#!/bin/bash
set -ux

# setup the network using the old binary

OLD_VERSION=${OLD_VERSION:-"v0.41.3"}
WASM_PATH=${WASM_PATH:-"./scripts/wasm_file/swapmap.wasm"}
ARGS="--chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block"
NEW_VERSION=${NEW_VERSION:-"v0.41.4"}
MIGRATE_MSG=${MIGRATE_MSG:-'{}'}
EXECUTE_MSG=${EXECUTE_MSG:-'{"ping":{}}'}
docker_command="docker-compose -f docker-compose-e2e-upgrade.yml exec"
validator1_command="$docker_command validator1 bash -c"
validator2_command="$docker_command validator2 bash -c"
working_dir=/workspace/.oraid
VALIDATOR_HOME=${VALIDATOR_HOME:-"$working_dir/validator1"}

# setup local network
sh $PWD/multinode-docker.sh

# # deploy new contract
# store_ret=$(oraid tx wasm store $WASM_PATH --from validator1 --home $VALIDATOR_HOME $ARGS --output json)
# code_id=$(echo $store_ret | jq -r '.logs[0].events[1].attributes[] | select(.key | contains("code_id")).value')
# oraid tx wasm instantiate $code_id '{}' --label 'testing' --from validator1 --home $VALIDATOR_HOME -b block --admin $(oraid keys show validator1 --keyring-backend test --home $VALIDATOR_HOME -a) $ARGS
# contract_address=$(oraid query wasm list-contract-by-code $code_id --output json | jq -r '.contracts[0]')

# echo "contract address: $contract_address"

# # create new upgrade proposal
UPGRADE_HEIGHT=${UPGRADE_HEIGHT:-19}
$validator1_command "oraid tx gov submit-proposal software-upgrade $NEW_VERSION --title 'foobar' --description 'foobar' --from validator1 --upgrade-height $UPGRADE_HEIGHT --upgrade-info 'https://github.com/oraichain/orai/releases/download/v0.41.4-rc0/manifest.json' --deposit 10000000orai $ARGS --home $VALIDATOR_HOME"
$validator1_command "oraid tx gov vote 1 yes --from validator1 --home $VALIDATOR_HOME $ARGS"
$validator1_command "oraid tx gov vote 1 yes --from validator2 --home $working_dir/validator2 $ARGS"

# sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."
sleep 1m

height_before=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

re='^[0-9]+$'
if ! [[ $height_before =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

sleep 30s

height_after=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

re='^[0-9]+$'
if ! [[ $height_after =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

if [ $height_after -gt $height_before ]
then
echo "Test Passed"
else
echo "Test Failed"
fi