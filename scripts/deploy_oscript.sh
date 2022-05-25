#!/bin/bash

CHAIN_ID=${CHAIN_ID:-Oraichain}
IFS=',' read -r -a DS <<< "$1"
DS=${DS:-classification}
DS_RAW=${1:-classification}
IFS=',' read -r -a TC <<< "$2"
TC=${TC:-classification_testcase}
TC_RAW=${2:-classification_testcase}
OS=${3:-classification_oscript}
OS_DS=''
OS_TC=''
NONCE=${5:-1}
DIR_PATH=${6:-$PWD}
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

OS_INPUT=${4:-'{"ai_data_source":['$OS_DS'],"testcase":['$OS_TC']}'}

sh $PWD/scripts/deploy-contract-store-addr.sh $DIR_PATH/smart-contracts/$OS/artifacts/$OS.wasm "$OS $NONCE" "$OS_INPUT" $PASS

# check if the oracle script exists or not
oraid query provider oscript $OS 2> is_exist.txt
description="test $OS"
address=$(cat ../address.txt)
# if the file is empty, then the oracle script does not exist. We create new

echo $DS_RAW
echo $TC_RAW

if [ -s is_exist.txt ]
then
    echo $PASS | oraid tx provider set-oscript $OS $address "$description" --ds ${DS:-classification} --tc ${TC:-classification_testcase} --from $USER --chain-id $CHAIN_ID -y
else
    echo $PASS | oraid tx provider edit-oscript $OS $OS $address "$description" --ds $DS_RAW --tc $TC_RAW --from $USER --chain-id $CHAIN_ID -y
fi