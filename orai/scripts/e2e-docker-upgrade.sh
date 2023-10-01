#!/bin/bash
set -ux

# setup the network using the old binary

OLD_VERSION=${OLD_VERSION:-"v0.41.4"}
ARGS="--chain-id testing -y --keyring-backend test --fees 200orai --gas 20000000 --gas-adjustment 1.5 -b block"
NEW_VERSION=${NEW_VERSION:-"v0.41.5"}
UPGRADE_INFO_VERSION=${UPGRADE_INFO_VERSION:-"v0.41.5"}
MIGRATE_MSG=${MIGRATE_MSG:-'{}'}
EXECUTE_MSG=${EXECUTE_MSG:-'{"ping":{}}'}
docker_command="docker-compose -f $PWD/docker-compose-e2e-upgrade.yml exec"
validator1_command="$docker_command validator1 bash -c"
validator2_command="$docker_command validator2 bash -c"
working_dir=/workspace
oraid_dir=$working_dir/.oraid
VALIDATOR_HOME=${VALIDATOR_HOME:-"$oraid_dir/validator1"}
WASM_PATH=${WASM_PATH:-"$working_dir/scripts/wasm_file/swapmap.wasm"}

# setup local network
sh $PWD/scripts/multinode-docker.sh

# # deploy new contract
store_ret=`$validator1_command "oraid tx wasm store $WASM_PATH --from validator1 --home $VALIDATOR_HOME $ARGS --output json"`
code_id=$(echo $store_ret | jq -r '.logs[0].events[1].attributes[] | select(.key | contains("code_id")).value')
$validator1_command "oraid tx wasm instantiate $code_id \"{}\" --label 'testing' --from validator1 --home $VALIDATOR_HOME -b block --no-admin $ARGS"
contract_address_res=`$validator1_command "oraid query wasm list-contract-by-code $code_id --output json | jq -r '.contracts[0]'"`
contract_address=$(echo "$contract_address_res" | tr -d -c '[:alnum:]') # remove all special characters because of the command's result

echo "contract address: $contract_address"

# create new upgrade proposal
UPGRADE_HEIGHT=${UPGRADE_HEIGHT:-19}
$validator1_command "oraid tx gov submit-proposal software-upgrade $NEW_VERSION --title 'foobar' --description 'foobar' --from validator1 --upgrade-height $UPGRADE_HEIGHT --upgrade-info 'https://github.com/oraichain/orai/releases/download/$UPGRADE_INFO_VERSION/manifest.json' --deposit 10000000orai $ARGS --home $VALIDATOR_HOME"
$validator1_command "oraid tx gov vote 1 yes --from validator1 --home $VALIDATOR_HOME $ARGS"
$validator1_command "oraid tx gov vote 1 yes --from validator2 --home $oraid_dir/validator2 $ARGS"

# sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."
sleep 12

# Check if latest height is less than the upgrade height
latest_height=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')
echo $latest_height
while [ $latest_height -lt $UPGRADE_HEIGHT ];
do
   sleep 7
   ((latest_height=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')))
   echo $latest_height
done

$validator1_command "oraid tx wasm execute $contract_address $(echo $EXECUTE_MSG | jq '@json') --from validator1 $ARGS --home $VALIDATOR_HOME"

height_before=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

re='^[0-9]+([.][0-9]+)?$'
if ! [[ $height_before =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

sleep 30

height_after=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

if ! [[ $height_after =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

if [ $height_after -gt $height_before ]
then
echo "Upgrade Passed"
else
echo "Upgarde Failed"
fi

inflation=$(curl --no-progress-meter http://localhost:1317/cosmos/mint/v1beta1/inflation | jq '.inflation | tonumber')
if ! [[ $inflation =~ $re ]] ; then
   echo "Error: Cannot query inflation => Potentially missing Go GRPC backport" >&2;
   echo "Tests Failed"; exit 1
fi

echo "Tests Passed"