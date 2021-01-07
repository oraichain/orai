#!/bin/bash

#
# Copyright Oraichain All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# export so other script can access

# colors
BROWN='\033[0;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'
BOLD=$(tput bold)
NORMAL=$(tput sgr0)

# environment
BASE_DIR=$PWD
SCRIPT_NAME=`basename "$0"`

# verify the result of the end-to-end test
verifyResult() {  
  if [ $1 -ne 0 ] ; then    
    printBoldColor $RED  "========= $2 ==========="
    echo
      exit 1
  fi
}

printCommand(){
  echo -e ""
  printBoldColor $BROWN "Command:"
  printBoldColor $BLUE "\t$1"  
}

printBoldColor(){
  echo -e "$1${BOLD}$2${NC}${NORMAL}"
}

# Print the usage message
printHelp () {

  echo $BOLD "Usage: "  
  echo "  $SCRIPT_NAME -h|--help (Show help)"  
  echo $NORMAL

  if [[ ! -z $2 ]]; then
    res=$(printHelp 0 | grep -A2 "\- '$2' \-")
    echo "$res"    
  else      
    printBoldColor $BROWN "      - 'start' - Run the full node"
    printBoldColor $BLUE  "          fn start"           
    echo
    printBoldColor $BROWN "      - 'broadcast' - broadcast transaction"
    printBoldColor $BLUE  "          fn broadcast --key value"           
    echo
    printBoldColor $BROWN "      - 'init' - Init the orai node"
    printBoldColor $BLUE  "          fn init"           
    echo
    printBoldColor $BROWN "      - 'sign' - sign transaction"
    printBoldColor $BLUE  "          fn sign --key value"           
    echo
    printBoldColor $BROWN "      - 'initScript' - init AI request script"
    printBoldColor $BLUE  "          fn initScript --key value"           
    echo
    printBoldColor $BROWN "      - 'clear' - Clear all existing data"
    printBoldColor $BLUE  "          fn clear"           
    echo
  fi

  echo
  echo "  $SCRIPT_NAME method --argument=value"
  
  # default exit as 0
  exit ${1:-0}
}


# Get a value:
getArgument() {     
  local key="args_${1/-/_}"  
  echo ${!key:-$2}  
}


# check first param is method
if [[ $1 =~ ^[a-z] ]]; then 
  METHOD=$1
  shift
fi

# use [[ ]] we dont have to quote string
args=()
case "$METHOD" in
  bash)        
    while [[ ! -z $2 ]];do         
      if [[ ${1:0:2} == '--' ]]; then
        KEY=${1/--/}            
        if [[ $KEY =~ ^([a-zA-Z_-]+)=(.+) ]]; then                
          declare "args_${BASH_REMATCH[1]/-/_}=${BASH_REMATCH[2]}"
        else          
          declare "args_${KEY/-/_}=$2" 
          shift
        fi    
      else 
        args+=($1)
      fi
      shift
    done
    QUERY="$@"            
  ;;
  config)        
    while [[ $# -gt 0 ]] ; do            
      if [[ ${1:0:2} == '--' ]]; then
        KEY=${1/--/}            
        if [[ $KEY =~ ^([a-zA-Z_-]+)=(.+) ]]; then                
          declare "args_${BASH_REMATCH[1]/-/_}=${BASH_REMATCH[2]}"
        else          
          declare "args_${KEY/-/_}=$2" 
          shift
        fi    
      else 
        args+=($1)
      fi
      shift
    done     
  ;;
  *) 
    # normal processing
    while [[ $# -gt 0 ]] ; do                
      if [[ ${1:0:2} == '--' ]]; then
        KEY=${1/--/}                
        if [[ $KEY =~ ^([a-zA-Z_-]+)=(.+) ]]; then         
          declare "args_${BASH_REMATCH[1]/-/_}=${BASH_REMATCH[2]}"
        else
          declare "args_${KEY/-/_}=$2"        
          shift
        fi    
      else 
        case "$1" in
          -h|\?)            
            printHelp 0 $2
          ;;
          *)  
            args+=($1)            
          ;;  
        esac    
      fi 
      shift
    done 
  ;; 
esac


clear(){
    rm -rf .oraid/
    rm -rf .oraicli/
    rm -rf .oraifiles/    
}

oraidFn(){
    # oraid start
    orai start --chain-id $CHAIN_ID --laddr tcp://0.0.0.0:1317 --node tcp://0.0.0.0:26657 # --trust-node
}


initFn(){
  clear    
  oraid init $MONIKER --chain-id Oraichain
  res=$?        
  verifyResult $res "can not run oraid init"  

  # Configure your CLI to eliminate need to declare them as flags
  oraicli config chain-id Oraichain
  oraicli config output json
  oraicli config indent true
  oraicli config trust-node true
  oraicli config keyring-backend test

  oraicli keys add $USER
  res=$?        
  verifyResult $res "can not add $USER"  

  oraid add-genesis-account $(oraicli keys show $USER -a) 9000000000000000orai
  res=$?        
  verifyResult $res "can not add-genesis-account $USER"  

  oraid gentx --keyring-backend test --amount 900000000000orai --name $USER --min-self-delegation $MIN_SELF_DELEGATION
  res=$?        
  verifyResult $res "can not gentx $USER"  

  # put the validators into the genesis file so the chain is aware of the validators
  oraid collect-gentxs

  oraid validate-genesis

  # run at background without websocket
  oraid start --minimum-gas-prices $GAS_PRICES &  

  # 30 seconds timeout
  websocketInitFn    
  
  sleep 10
  pkill orai
  pkill oraid
  pkill oraicli
  pkill websocket
  sleep 2
}

websocketInitFn() {
  # run at background without websocket
  # # 30 seconds timeout to check if the node is alive or not
  timeout 30 bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:26657/health)" != "200" ]]; do sleep 1; done' || false
  local reporter="${USER}_reporter"
  # for i in $(eval echo {1..$2})
  # do
    # add reporter key

  ###################### init websocket for the validator

  HOME=$PWD/.oraid
  # rm -rf ~/.websocket
  WEBSOCKET="websocket --home $HOME"
  #$WEBSOCKET keys delete-all
  $WEBSOCKET keys add $reporter  

  # config chain id
  $WEBSOCKET config chain-id Oraichain

  # add validator to websocket config
  $WEBSOCKET config validator $(oraicli keys show $USER -a --bech val --keyring-backend test)

  # setup broadcast-timeout to websocket config
  $WEBSOCKET config broadcast-timeout "30s"

  # setup rpc-poll-interval to websocket config
  $WEBSOCKET config rpc-poll-interval "1s"

  # setup max-try to websocket config
  $WEBSOCKET config max-try 5

  # config log type
  $WEBSOCKET config log-level debug

  sleep 10

  # send orai tokens to reporters
  echo "y" | oraicli tx send $(oraicli keys show $USER -a) $($WEBSOCKET keys show $reporter) 10000000orai --from $(oraicli keys show $USER -a) --fees 5000orai

  sleep 6

  #wait for sending orai tokens transaction success

  # add reporter to oraichain
  echo "y" | oraicli tx websocket add-reporters $($WEBSOCKET keys list -a) --from $USER --fees 5000orai --keyring-backend test
  sleep 8
  pkill oraid
}


initScriptFn(){
  echo "y" | oraicli tx provider set-datasource coingecko_eth ./testfiles/coingecko_eth.py "A data source that fetches the ETH price from Coingecko API" --from $USER --fees 5000orai

  sleep 5

  echo "y" | oraicli tx provider set-datasource crypto_compare_eth ./testfiles/crypto_compare_eth.py "A data source that collects ETH price from crypto compare" --from $USER --fees 5000orai

  sleep 5

  echo "y" | oraicli tx provider set-testcase testcase_price ./testfiles/testcase_price.py "A sample test case that uses the expected output of users provided to verify the bitcoin price from the datasource" --from $USER --fees 5000orai

  sleep 5

  echo "y" | oraicli tx provider set-oscript oscript_eth ./testfiles/oscript_eth.py "An oracle script that fetches and aggregates ETH price from different sources" --ds coingecko_eth,crypto_compare_eth --tc testcase_price --from $USER --fees 5000orai
}

unsignedFn(){
  local id=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.address" -r)
  local unsigned=$(curl --location --request POST 'http://localhost:1317/airequest/aireq' \
--header 'Content-Type: application/json' \
--data-raw '{
    "base_req":{
        "from":"'$id'",
        "chain_id":"'$CHAIN_ID'"
    },
    "oracle_script_name":"oscript_eth",
    "input":"",
    "expected_output":{"price":"5000"},
    "fees":"60000orai",
    "validator_count": "1"
}' > tmp/unsignedTx.json)

    res=$?  
    verifyResult $res "Unsigned failed"
}

unsignedSetDsFn(){
  local id=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.address" -r)
  local unsigned=$(curl --location --request POST 'http://localhost:1317/provider/datasource' \
--header 'Content-Type: application/json' \
--data-raw '{
    "base_req":{
        "from":"'$id'",
        "chain_id":"Oraichain"
    },
    "name":"coingecko_eth",
    "code_path":"/workspace/testfiles/coingecko_eth.py",
    "description":"NTAwMA==",
    "fees":"60000orai",
    "test":["abc","efgh"]
}' > tmp/unsignedTx.json)
}

clear(){
    rm -rf .oraid/
    rm -rf .oraicli/
    rm -rf .oraifiles/    
}

signFn(){     
    # $1 is account number
    local sequence=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.sequence" -r)
    local acc_num=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.account_number" -r)
    oraicli tx sign tmp/unsignedTx.json --from $USER --offline --chain-id $CHAIN_ID --sequence $sequence --account-number $acc_num > tmp/signedTx.json
    oraicli tx broadcast tmp/signedTx.json

    res=$?  
    verifyResult $res "Signed failed"
}


USER=$(getArgument "user" $USER)
CHAIN_ID=$(getArgument "chain-id" Oraichain)

# processing
case "${METHOD}" in     
  hello)
    helloFn
  ;;
  init)
    initFn
  ;;
  initDev)
    initDevFn
  ;;
  start)
    oraidFn
  ;;  
  unsign)
    unsignedFn
  ;;
  unsignedSetDs)
    unsignedSetDsFn
  ;;
  initScript)
    initScriptFn
  ;;
  sign)
    signFn
  ;;
  broadcast)
    broadcastFn
  ;;  
  clear)
    clear
  ;; 
  *) 
    printHelp 1 ${args[0]}
  ;;
esac