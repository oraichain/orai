#!/bin/bash

CHAIN_ID=${CHAIN_ID:-Oraichain}
IFS=',' read -r -a DS <<< "$1"
DS=${DS:-classification}
DS_INPUT=${2:-''}
NONCE=${3:-1}
DIR_PATH=${4:-$PWD}
PASS=${5:-123456789}

echo "data sources: ${DS[@]}"

# deploy smart contract data source and create data source
for i in "${DS[@]}"
do
    sh $PWD/scripts/deploy-contract-store-addr.sh $DIR_PATH/smart-contracts/$i/artifacts/$i.wasm "$i $NONCE" "$DS_INPUT" $PASS

    # check if the data source exists or not
    oraid query provider dsource $i 2> is_exist.txt
    description="test $i"
    address=$(cat ../address.txt)
    echo "address: $address"
    echo "description: $description"
    echo "data source file: $i"

    # if the file is not empty, then the data source does not exist. We create new
    if [ -s is_exist.txt ]
    then
        echo $PASS | oraid tx provider set-datasource $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    else
        # if it exists already, we update the contract
        echo $PASS | oraid tx provider edit-datasource $i $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    fi

done