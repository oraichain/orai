#!/bin/sh
set -o errexit -o nounset -o pipefail

CHAIN_ID=${CHAIN_ID:-Oraichain}
PASS=${1:-123456789}

sh ./scripts/deploy-contract-store-addr.sh smart-contracts/classification/artifacts/classification.wasm "classification 1" '' $PASS

echo $PASS | oraid tx provider set-datasource classification_ds $(cat ../address.txt) "test classification" --from $USER --chain-id $CHAIN_ID -y

sh ./scripts/deploy-contract-store-addr.sh smart-contracts/cv009/artifacts/cv009.wasm "cv009 1" '' $PASS

echo $PASS | oraid tx provider set-datasource cv009 $(cat ../address.txt) "test cv009" --from $USER --chain-id $CHAIN_ID -y

sh ./scripts/deploy-contract-store-addr.sh smart-contracts/classification-testcase/artifacts/classification_testcase.wasm "classification-testcase 1" '' $PASS

echo $PASS | oraid tx provider set-testcase classification_tc $(cat ../address.txt) "test classification test case" --from $USER --chain-id $CHAIN_ID -y

sh ./scripts/deploy-contract-store-addr.sh smart-contracts/classification-oscript/artifacts/classification_oscript.wasm "classification-oscript 1" '{"ai_data_source":["classification_ds","cv009"],"testcase":["classification_tc"]}' $PASS

echo $PASS | oraid tx provider set-oscript classification_oscript $(cat ../address.txt) "test classifcation oscript" --ds classification_ds,cv009 --tc classification_tc --from $USER --chain-id $CHAIN_ID -y