#!/bin/sh
set -o errexit -o nounset -o pipefail

CHAIN_ID=${CHAIN_ID:-Oraichain}
PASS=${1:-123456789}

sh ./scripts/deploy-contract-store-addr.sh smart-contracts/datasource-eth/artifacts/datasource_eth.wasm "datasource-eth 1" '' $PASS

echo $PASS | oraid tx provider set-datasource eth_price $(cat ../address.txt) "test eth price" --from $USER --chain-id $CHAIN_ID -y

sh ./scripts/deploy-contract-store-addr.sh smart-contracts/testcase-price/artifacts/testcase_price.wasm "testcase-price 1" '' $PASS

echo $PASS | oraid tx provider set-testcase eth_price_testcase $(cat ../address.txt) "test eth price testcase" --from $USER --chain-id $CHAIN_ID -y

sh ./scripts/deploy-contract-store-addr.sh smart-contracts/oscript-price/artifacts/oscript_price.wasm "oscript-price 1" '{"ai_data_source":"datasource_eth","testcase":"testcase_price"}' $PASS

echo $PASS | oraid tx provider set-oscript oscript_eth $(cat ../address.txt) "test eth price oracle script" --ds eth_price --tc eth_price_testcase --from $USER --chain-id $CHAIN_ID -y