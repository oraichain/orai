#!/bin/bash
set -u

# setup the network using the old binary

ARGS="--chain-id testing -y --keyring-backend test --gas 20000000 --gas-adjustment 1.5 -b block"
NEW_VERSION=${NEW_VERSION:-"v0.42.2"}
UPGRADE_INFO_VERSION=${UPGRADE_INFO_VERSION:-"v0.42.2"}
MIGRATE_MSG=${MIGRATE_MSG:-'{}'}
EXECUTE_MSG=${EXECUTE_MSG:-'{"ping":{}}'}
docker_command="docker-compose -f $PWD/docker-compose-e2e-upgrade.yml exec"
validator1_command="$docker_command validator1 bash -c"
validator2_command="$docker_command validator2 bash -c"
working_dir=/workspace
oraid_dir=$working_dir/.oraid
VALIDATOR_HOME=${VALIDATOR_HOME:-"$oraid_dir/validator1"}
WASM_PATH=${WASM_PATH:-"$working_dir/scripts/wasm_file/swapmap.wasm"}
HIDE_LOGS="/dev/null"

# setup local network
sh $PWD/scripts/multinode-docker.sh

# create new upgrade proposal
UPGRADE_HEIGHT=${UPGRADE_HEIGHT:-85}
$validator1_command "oraid tx gov submit-proposal software-upgrade $NEW_VERSION --title 'foobar' --description 'foobar' --from validator1 --upgrade-height $UPGRADE_HEIGHT --upgrade-info 'https://github.com/oraichain/orai/releases/download/$UPGRADE_INFO_VERSION/manifest.json' --deposit 10000000orai $ARGS --home $VALIDATOR_HOME > $HIDE_LOGS"
$validator1_command "oraid tx gov vote 1 yes --from validator1 --home $VALIDATOR_HOME $ARGS > $HIDE_LOGS"
$validator1_command "oraid tx gov vote 1 yes --from validator2 --home $oraid_dir/validator2 $ARGS > $HIDE_LOGS"

# sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."

# Check if latest height is less than the upgrade height
latest_height=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')
echo $latest_height
while [ $latest_height -lt $UPGRADE_HEIGHT ];
do
   sleep 5
   ((latest_height=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')))
   echo $latest_height
done

# sleep about 5 secs to wait for the rest & json rpc server to be u
echo "Waiting for the REST & JSONRPC servers to be up ..."
sleep 19

oraid_version=$($validator1_command "oraid version")
if ! [[ $oraid_version =~ $NEW_VERSION ]] ; then
   echo "The chain has not upgraded yet. There's something wrong!"; exit 1
fi

height_before=$(curl --no-progress-meter http://localhost:1317/blocks/latest | jq '.block.header.height | tonumber')

re='^[0-9]+([.][0-9]+)?$'
if ! [[ $height_before =~ $re ]] ; then
   echo "error: Not a number" >&2; exit 1
fi

sleep 5

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

result=$(curl --no-progress-meter http://localhost:8545/ -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"net_listening","params":[],"id":1}' | jq '.result')

if ! [[ $result =~ true ]] ; then
   echo "Error: Cannot query JSONRPC" >&2;
   echo "Tests Failed"; exit 1
fi

evm_denom=$(curl --no-progress-meter http://localhost:1317/ethermint/evm/v1/params | jq '.params.evm_denom')
if ! [[ $evm_denom =~ "aorai" ]] ; then
   echo "Error: EVM denom is not correct. The upgraded version is not the latest!" >&2;
   echo "Tests Failed"; exit 1
fi

NODE_HOME=$VALIDATOR_HOME USER=validator1 WASM_PATH=$WASM_PATH bash $PWD/scripts/tests-0.42.1/test-gasless-docker.sh
NODE_HOME=$VALIDATOR_HOME USER=validator1 WASM_PATH=$WASM_PATH bash $PWD/scripts/tests-0.42.1/test-tokenfactory-docker.sh
NODE_HOME=$VALIDATOR_HOME USER=validator1 WASM_PATH=$WASM_PATH bash $PWD/scripts/tests-0.42.1/test-evm-cosmos-mapping-docker.sh

echo "Tests Passed"
