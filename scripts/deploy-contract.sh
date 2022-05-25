#!/bin/sh
set -o errexit -o nounset -o pipefail

contract_path=$1
label=$2
init=${3:-{\}}
code_id=${4:-}
CHAIN_ID=${CHAIN_ID:-Oraichain}

echo "Enter passphrase:"
read -s passphrase

if [ -z $code_id ]
then 
    store_ret=$(echo $passphrase | oraid tx wasm store $contract_path --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=$CHAIN_ID -y)
    echo $store_ret
    code_id=$(echo $store_ret | jq -r '.logs[0].events[0].attributes[] | select(.key | contains("code_id")).value')
fi 

# echo "oraid tx wasm instantiate $code_id '$init' --from $USER --label '$label' --gas auto --gas-adjustment 1.2 --chain-id=$CHAIN_ID -y"
# quote string with "" with escape content inside which contains " characters
(echo $passphrase;echo $passphrase) | oraid tx wasm instantiate $code_id "$init" --from $USER --label "$label" --gas auto --gas-adjustment 1.2 --chain-id=$CHAIN_ID -y
contract_address=$(oraid query wasm list-contract-by-code $code_id | grep address | awk '{print $(NF)}')
