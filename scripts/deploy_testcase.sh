#!/bin/bash

CHAIN_ID=${CHAIN_ID:-Oraichain}
IFS=',' read -r -a TC <<< "$1"
TC=${TC:-classification_testcase}
TC_INPUT=${2:-''}
NONCE=${3:-1}
DIR_PATH=${4:-$PWD}
PASS=${5:-123456789}

# deploy smart contract test case and create test case
for i in "${TC[@]}"
do
    sh $PWD/scripts/deploy-contract-store-addr.sh $DIR_PATH/smart-contracts/$i/artifacts/$i.wasm "$i $NONCE" "$TC_INPUT" $PASS

    # check if the test case exists or not
    oraid query provider tcase $i 2> is_exist.txt
    description="test $i"
    address=$(cat ../address.txt)

    # if the file is empty, then the test case does not exist. We create new
    if [ -s is_exist.txt ]
    then
        echo $PASS | oraid tx provider set-testcase $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    else
        # if it exists already, we update the contract
        echo $PASS | oraid tx provider edit-testcase $i $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    fi
done