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
NC='\033[0m'
BOLD=$(tput bold)
NORMAL=$(tput sgr0)

# environment
BASE_DIR=$PWD
SCRIPT_NAME=`basename "$0"`

# verify the result of the end-to-end test
verifyResult() {  
  if [ $1 -ne 0 ] ; then
    echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute End-2-End Scenario ==========="
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
    printBoldColor $BROWN "      - 'hello' - hello"
    printBoldColor $BLUE  "          fn hello --key value"           
    echo
    printBoldColor $BROWN "      - 'broadcast' - broadcast"
    printBoldColor $BLUE  "          fn broadcast --key value"           
    echo
    printBoldColor $BROWN "      - 'restServer' - restServer"
    printBoldColor $BLUE  "          fn restServer --key value"           
    echo
    printBoldColor $BROWN "      - 'sign' - sign"
    printBoldColor $BLUE  "          fn sign --key value"           
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


helloFn(){
    local value=$(getArgument "value" world) 
    echo "$value"
}

oraidFn(){
    oraid start
}

websocketInitFn(){
    local reporter="${USER}_reporter"
    ./websocket-init.sh $USER $reporter
}

websocketRunFn(){
    websocket run
}

restServerFn(){
    oraicli rest-server --chain-id Oraichain --laddr tcp://0.0.0.0:1317  --trust-node
}

initScriptFn(){
  oraicli tx provider set-datasource coingecko_eth ./testfiles/coingecko_eth.py "A data source that fetches the ETH price from Coingecko API" --from $USER --fees 5000orai

  sleep 5

  oraicli tx provider set-datasource crypto_compare_eth ./testfiles/crypto_compare_eth.py "A data source that collects ETH price from crypto compare" --from $USER --fees 5000orai

  sleep 5

  oraicli tx provider set-testcase testcase_price ./testfiles/testcase_price.py "A sample test case that uses the expected output of users provided to verify the bitcoin price from the datasource" --from $USER --fees 5000orai

  sleep 5

  oraicli tx provider set-oscript oscript_eth ./testfiles/oscript_eth.py "An oracle script that fetches and aggregates ETH price from different sources" --from $USER --fees 5000orai
}

unsignedFn(){
  local id=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.address" -r)
  local unsigned=$(curl --location --request POST 'http://localhost:1317/airequest/aireq/pricereq' \
--header 'Content-Type: application/json' \
--data-raw '{
    "base_req":{
        "from":"'$id'",
        "chain_id":"Oraichain"
    },
    "oracle_script_name":"oscript_eth",
    "input":"",
    "expected_output":"NTAwMA==",
    "fees":"60000orai",
    "validator_count": "1"
}' > tmp/unsignedTx.json)
}

signFn(){     
    # $1 is account number
    local sequence=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.sequence" -r)
    local acc_num=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.account_number" -r)
    oraicli tx sign tmp/unsignedTx.json --from $USER --offline --chain-id Oraichain --sequence $sequence --account-number $acc_num > tmp/signedTx.json
}

broadcastFn(){
    oraicli tx broadcast tmp/signedTx.json
}

USER=$(getArgument "user" duc)

# processing
case "${METHOD}" in     
  hello)
    helloFn
  ;;
  oraid)
    oraidFn
  ;;
  websocketInit)
    websocketInitFn
  ;;
  websocketRun)
    websocketRunFn
  ;;
  unsign)
    unsignedFn
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
  restServer)
    restServerFn
  ;;
  *) 
    printHelp 1 ${args[0]}
  ;;
esac