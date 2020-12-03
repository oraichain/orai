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
  echo "  $SCRIPT_NAME [-m|--method=] hello|sign"
  echo "  $SCRIPT_NAME -h|--help (print this message)"  
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

websocketFn(){    
    websocket run
}

broadcastFn(){
    oraicli tx broadcast tmp/signedTx.json
}

restServerFn(){
    oraicli rest-server --chain-id Oraichain --laddr tcp://0.0.0.0:1317  --trust-node
}

signFn(){     
    # $1 is account number
    local sequence=$(curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show $USER -a)" | jq ".result.value.sequence" -r)
    oraicli tx sign tmp/unsignedTx.json --from $USER --offline --chain-id Oraichain --sequence $sequence --account-number $1 > tmp/signedTx.json
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
  websocket)
    websocketFn
  ;;
  sign)
    signFn 3
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