#!/bin/bash

CHAIN_ID=${CHAIN_ID:-Oraichain}
IFS=',' read -r -a DS <<< "$1"
DS=${DS:-classification}
DS_RAW=${1:-classification}
IFS=',' read -r -a TC <<< "$2"
TC=${TC:-classification_testcase}
TC_RAW=${2:-classification_testcase}
OS=${3:-classification_oscript}
DS_INPUT=${4:-''}
TC_INPUT=${5:-''}
OS_DS=''
OS_TC=''
NONCE=${6:-1}
PASS=${7:-123456789}

# add double quotes in the list of data sources
for ((i=0; i<${#DS[@]}; i++));
do
    OS_DS+=\"${DS[$i]}\",
done
# remove the last character (comma)
OS_DS=${OS_DS::-1}

# add double quotes in the list of test cases
for ((i=0; i<${#TC[@]}; i++));
do
    OS_TC+=\"${TC[$i]}\",
done
# remove the last character (comma)
OS_TC=${OS_TC::-1}

OS_INPUT=${6:-'{"ai_data_source":['$OS_DS'],"testcase":['$OS_TC']}'}

echo $OS_INPUT

# deploy smart contract data source and create data source
for i in "${DS[@]}"
do
    sh $PWD/scripts/deploy-contract-store-addr.sh $PWD/smart-contracts/$i/artifacts/$i.wasm "$i $NONCE" "$DS_INPUT" $PASS

    # check if the data source exists or not
    oraid query provider dsource $i > is_exist.txt
    description="test $i"
    address=$(cat ../address.txt)
    echo "address: $address"
    echo "description: $description"
    echo "data source file: $i"

    # if the file is empty, then the data source does not exist. We create new
    if [ ! -s is_exist.txt ]
    then
        echo $PASS | oraid tx provider set-datasource $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    else
        # if it exists already, we update the contract
        echo $PASS | oraid tx provider edit-datasource $i $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    fi

done

# deploy smart contract test case and create test case
for i in "${TC[@]}"
do
    sh $PWD/scripts/deploy-contract-store-addr.sh $PWD/smart-contracts/$i/artifacts/$i.wasm "$i $NONCE" "$TC_INPUT" $PASS

    # check if the test case exists or not
    oraid query provider tcase $i > is_exist.txt
    description="test $i"
    address=$(cat ../address.txt)

    # if the file is empty, then the test case does not exist. We create new
    if [ ! -s is_exist.txt ]
    then
        echo $PASS | oraid tx provider set-testcase $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    else
        # if it exists already, we update the contract
        echo $PASS | oraid tx provider edit-testcase $i $i $address "$description" --from $USER --chain-id $CHAIN_ID -y
    fi
done

sh  $PWD/scripts/deploy-contract-store-addr.sh $PWD/smart-contracts/$OS/artifacts/$OS.wasm "$OS $NONCE" "$OS_INPUT" $PASS

# check if the oracle script exists or not
oraid query provider oscript $OS > is_exist.txt
description="test $OS"
address=$(cat ../address.txt)
# if the file is empty, then the oracle script does not exist. We create new
if [ ! -s is_exist.txt ]
then
    echo $PASS | oraid tx provider set-oscript $OS $address "$description" --ds ${DS:-classification} --tc ${TC:-classification_testcase} --from $USER --chain-id $CHAIN_ID -y
else
    echo $PASS | oraid tx provider edit-oscript $OS $OS $address "$description" --ds $DS_RAW --tc $TC_RAW --from $USER --chain-id $CHAIN_ID -y
fi
