#!/bin/sh
set -o errexit -o nounset -o pipefail

contract_path=$1
label="$2"
echo "Enter passphrase:"
read -s passphrase

store_ret=$(echo $passphrase | oraid tx wasm store $contract_path --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y)
echo $store_ret
code_id=$(echo $store_ret | jq -r '.logs[0].events[0].attributes[] | select(.key | contains("code_id")).value')
echo $passphrase | oraid tx wasm instantiate $code_id '{}' --from $USER --label "$label" --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y
contract_address=$(oraid query wasm list-contract-by-code $code_id | grep address | awk '{print $(NF)}')

echo "Contract address: $contract_address"

oraid query wasm contract-state smart $contract_address '{"test_price":{"contract":"orai1wgnpjk9s5dgyashvf4xl9hy3kzkfqynjjmnt8y","output":"1333.65"}}'