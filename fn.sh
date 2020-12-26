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
    orai start --chain-id Oraichain --laddr tcp://0.0.0.0:1317 --node tcp://0.0.0.0:26657 # --trust-node
}


initFn(){ 

  ### Check if a directory does not exist ###
  if [[ ! -d "$PWD/.oraid/" || ! -d "$PWD/.oraicli/" ]] 
  then
    echo "Directory /path/to/dir DOES NOT exists."

    oraid init $(getArgument "moniker" $MONIKER) --chain-id Oraichain

    # Configure your CLI to eliminate need to declare them as flags
    oraicli config chain-id Oraichain
    oraicli config output json
    oraicli config indent true
    oraicli config trust-node true

    expect -c "

    spawn oraicli keys add $USER --recover
    expect {
        \"override the existing name*\" {send -- \"y\r\"}
    }

    expect \"*bip39 mnemonic\"

    send -- \"$MNEMONIC\r\"

    expect {
        \"Enter keyring passphrase:\" send -- {\"$PASS\r\"; exp_continue }
        \"Re-enter keyring passphrase:\" send -- {\"$PASS\r\"; exp_continue }
        eof
    }"
#EOF
#     expect << EOF

#     spawn oraicli keys add $USER --recover
#     expect {
#         "override the existing name*" {send -- "y\r"}
#     }

#     expect "*bip39 mnemonic"

#     send -- "$MNEMONIC\r"

#     expect {
#         "Enter keyring passphrase:" send -- {"$PASS\r"; exp_continue }
#         "Re-enter keyring passphrase:" send -- {"$PASS\r"; exp_continue }
#         eof
#     }
# EOF

    # download genesis json file
  
    curl https://raw.githubusercontent.com/oraichain/oraichain-static-files/ducphamle2-test/genesis.json > .oraid/config/genesis.json
    
    # rm -f .oraid/config/genesis.json && wget https://raw.githubusercontent.com/oraichain/oraichain-static-files/ducphamle2-test/genesis.json -q -P .oraid/config/

    # add persistent peers to listen to blocks
    # local persistentPeers=$(getArgument "--persistent_peers" "$PERSISTENT_PEERS")
    # [ ! -z $persistentPeers ] && sed -i 's/persistent_peers *= *".*"/persistent_peers = "$PERSISTENT_PEERS"/g' .oraid/config/config.toml 

    # add persistent peers to listen to blocks
    sed -i 's/persistent_peers *= *".*"/persistent_peers = "8bd5a1689706f72a595ad2cb3ea09ad7f9a47eb9@209.97.154.247:26656"/g' .oraid/config/config.toml

    oraid validate-genesis
    # done init
  fi
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

  sleep 2

  # send orai tokens to reporters
  echo "y" | oraicli tx send $(oraicli keys show $USER -a) $($WEBSOCKET keys show $reporter) 10000000orai --from $(oraicli keys show $USER -a) --fees 5000orai

  sleep 6

  #wait for sending orai tokens transaction success

  # add reporter to oraichain
  echo "y" | oraicli tx websocket add-reporters $($WEBSOCKET keys list -a) --from $USER --fees 5000orai --keyring-backend test
  sleep 8
  pkill oraid
}

createValidatorFn() {
  local user=$(getArgument "user" $USER)
  # run at background without websocket
  # # 30 seconds timeout to check if the node is alive or not, the '&' symbol allows to run below commands while still having the process running
  oraid start &
    # 30 seconds timeout
  timeout 30 bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:26657/health)" != "200" ]]; do sleep 1; done' || false

  # loop to query the account, since the account may not exist until a specific block, so we need to constantly check
  local acc=$(oraicli query auth account $(oraicli keys show $user -a) | jq .type)
  while [[ "$acc" != \""cosmos-sdk/Account"\" ]];
  do 
    # reset the value for the loop condition
    acc=$(oraicli query auth account $(oraicli keys show tester -a) | jq .type)
    sleep 10
  done

  local amount=$(getArgument "amount" $AMOUNT)
  local pubkey=$(oraid tendermint show-validator)
  local moniker=$(getArgument "moniker" $MONIKER)
  if [[ $moniker == "" ]]; then
    moniker="$USER"_"Oraichain"_$(($RANDOM%100000000000))
  fi
  local floatRe='^[+-]?[0-9]+\.?[0-9]*$'
  local commissionRate=$(getArgument "commission_rate" $COMMISSION_RATE)
  if [[ $commissionRate > 1 || !$commissionRate =~ $floatRe || $commissionRate == "" ]]; then
    commissionRate=0.10
  fi
  local commissionMaxRate=$(getArgument "commission_max_rate" $COMMISSION_MAX_RATE)
  if [[ $commissionMaxRate > 1 || !$commissionMaxRate =~ $floatRe || commissionMaxRate == "" ]]; then
    commissionMaxRate=0.20
  fi
  local commissionMaxChangeRate=$(getArgument "commission_max_change_rate" $COMMISSION_MAX_CHANGE_RATE)
  if [[ $commissionMaxChangeRate > 1 || !$commissionMaxChangeRate =~ $floatRe || $commissionMaxChangeRate == "" ]]; then
    commissionMaxChangeRate=0.01
  fi
  local minDelegation=$(getArgument "min_self_delegation" $MIN_SELF_DELEGATION)

  # verify env from user, regex for checking number
  local re='^[0-9]+$'
  if [[ $minDelegation < 1 || !$minDelegation =~ $re || $minDelegation == "" ]]; then
    minDelegation=1
  fi

  local gas=$(getArgument "gas" $GAS)
  if [[ $gas != "auto" && !$gas =~ $re ]]; then
    gas=200000
  fi

  # workaround, since auto gas in this case is not good, sometimes get out of gas
  if [[ $gas == "auto" || $gas < 200000 ]]; then
    gas=200000
  fi

  local gasPrices=$(getArgument "gas_prices" $GAS_PRICES)
  if [[ $gasPrices == "" ]]; then
    gasPrices="0.025orai"
  fi
  local securityContract=$(getArgument "security_contract" $SECURITY_CONTRACT)
  local identity=$(getArgument "identity" $IDENTITY)
  local website=$(getArgument "website" $WEBSITE)
  local details=$(getArgument "details" $DETAILS)

  sleep 10

  echo "y" | oraicli tx staking create-validator --amount $amount --pubkey $pubkey --moniker $moniker --chain-id Oraichain --commission-rate $commissionRate --commission-max-rate $commissionMaxRate --commission-max-change-rate $commissionMaxChangeRate --min-self-delegation $minDelegation --gas $gas --gas-prices $gasPrices --security-contact $securityContract --identity $identity --website $website --details $details --from $user

  local reporter="${user}_reporter"
  # # for i in $(eval echo {1..$2})
  # # do
  #   # add reporter key

  # ###################### init websocket for the validator

  HOME=$PWD/.oraid
  # rm -rf ~/.websocket
  WEBSOCKET="websocket --home $HOME"
  #$WEBSOCKET keys delete-all
  $WEBSOCKET keys add $reporter

  # config chain id
  $WEBSOCKET config chain-id Oraichain

  # add validator to websocket config
  $WEBSOCKET config validator $(oraicli keys show $user -a --bech val)

  # setup broadcast-timeout to websocket config
  $WEBSOCKET config broadcast-timeout "30s"

  # setup rpc-poll-interval to websocket config
  $WEBSOCKET config rpc-poll-interval "1s"

  # setup max-try to websocket config
  $WEBSOCKET config max-try 5

  # config log type
  $WEBSOCKET config log-level debug

  $WEBSOCKET config gas-prices $gasPrices

  $WEBSOCKET config gas $gas

  sleep 10

  # send orai tokens to reporters
  echo "y" | oraicli tx send $(oraicli keys show $user -a) $($WEBSOCKET keys show $reporter) 10000000orai --from $(oraicli keys show $user -a) --gas-prices $gasPrices

  sleep 10

  #wait for sending orai tokens transaction success

  # add reporter to oraichain
  echo "y" | oraicli tx websocket add-reporters $($WEBSOCKET keys list -a) --from $user --gas-prices $gasPrices
  sleep 8

  pkill oraid
}

USER=$(getArgument "user" $USER)

# processing
case "${METHOD}" in     
  hello)
    helloFn
  ;;
  init)
    initFn
  ;;
  websocketInit)
  websocketInitFn
  ;;
  start)
    oraidFn
  ;;
  initScript)
    initScriptFn
  ;; 
  clear)
    clear
  ;; 
  createValidator)
    createValidatorFn
  ;;
  *) 
    printHelp 1 ${args[0]}
  ;;
esac