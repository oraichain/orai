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
LOG_LEVEL=${LOG_LEVEL:-error}

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
    printBoldColor $BROWN "      - 'addGenAccount' - add a genesis account with balance"
    printBoldColor $BLUE  "          fn addGenAccount --address value --amount value"           
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

requirePass(){
  if [ -z "$PASS" ]
  then 
    echo "Enter passphrase: "  
    read PASS  
  fi 
}

helloFn() {  
  requirePass
  echo "passphrase: $PASS"
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
}

oraidFn(){
    pkill oraid
    pkill oraivisor
    if [[ -d "$PWD/.oraid/" ]] 
    then
      oraivisor start --rpc.laddr tcp://0.0.0.0:26657 --log_level $(getArgument "log_level" $LOG_LEVEL) --p2p.seeds $(getArgument "seeds" $SEEDS) --p2p.seed_mode $(getArgument "seed_mode" false)
    else
      tail -f /dev/null 
    fi
}

checkExpectProgram() {
  if [ ! `command -v expect` ]; then  
    echo "Installing expect ..."
    apk add expect
  fi  
}

addGenAccount(){
  address=$(getArgument "address")
  if [ -z "$address" ] 
  then 
    echo "Address is empty"
    exit 1
  fi 

  amount=$(getArgument "amount" "100000000000000")
  genesis=.oraid/config/genesis.json
  account='{"@type":"/cosmos.auth.v1beta1.BaseAccount","address":"'$address'","pub_key":null,"account_number":"0","sequence":"0"}'
  balance='{"address":"'$address'","coins":[{"denom":"orai","amount":"'$amount'"}]}'

  exsited_account=$(jq '.app_state.auth.accounts[] | select(.address == "'$address'")' $genesis)
  exsited_balance=$(jq '.app_state.bank.balances[] | select(.address == "'$address'")' $genesis)
  if [ ! -z "$exsited_account" ] || [ ! -z "$exsited_balance" ]
  then 
    echo "$address is already existed"
  else 
    result=$(jq --argjson account $account --argjson balance $balance '.app_state.auth.accounts += [$account] | .app_state.bank.balances += [$balance]' $genesis)
    echo $result | jq > $genesis
  fi
}

initFn(){ 

  checkExpectProgram
  
  ### Check if a directory does not exist ###
  if [[ ! -d "$PWD/.oraid/" ]] 
  then
    echo "Directory /path/to/dir DOES NOT exists."

    oraid init $(getArgument "moniker" $MONIKER) --chain-id Oraichain

    requirePass

    echo "passphrase: $PASS"

    sleep 10

    if [ -z "$MNEMONIC" ]
    then 
      echo "Enter mnemonic: "  
      read MNEMONIC
    fi 

    echo "mnemonic: $MNEMONIC"

    sleep 10

    # Configure your CLI to eliminate need to declare them as flags

    expect << EOF

        spawn oraid keys add $USER --recover
        expect {
          "override the existing name*" {send -- "y\r"}
        }

        expect "Enter your bip39 mnemonic*"
        send -- "$MNEMONIC\r"

        expect {
          "*passphrase:" { send -- "$PASS\r" }
        }
        expect {
          "*passphrase:" { send -- "$PASS\r" }
        }
        expect eof
EOF

    # download genesis json file
  
    wget -O .oraid/config/genesis.json $GENESIS_URL
    
    # rm -f .oraid/config/genesis.json && wget https://raw.githubusercontent.com/oraichain/oraichain-static-files/ducphamle2-test/genesis.json -q -P .oraid/config/

    # add persistent peers to listen to blocks
    local persistentPeers=$(getArgument "persistent_peers" "$SEEDS")
    # [ ! -z $persistentPeers ] && sed -i 's/seeds *= *".*"/seeds = "'"$persistentPeers"'"/g' .oraid/config/config.toml 

    # sed -i 's/persistent_peers *= *".*"/persistent_peers = "25e3dd0839fa44a89735b38b7b749acdfac8438e@164.90.180.95:26656,e07a89a185c538820258b977b01b44a806dfcece@157.230.22.169:26656,db13b4e2d1fd922640904590d6c9b5ae698de85c@165.232.118.44:26656,b46c45fdbb59ef0509d93e89e574b2080a146b14@178.128.61.252:26656,2a8c59cfdeccd2ed30471b90f626da09adcf3342@178.128.57.195:26656,b495da1980d3cd7c3686044e800412af53ae4be4@159.89.206.139:26656,addb91a1dbc48ffb7ddba30964ae649343179822@178.128.220.155:26656"/g' .oraid/config/config.toml

    oraid validate-genesis
    # done init
  fi
}

createValidatorFn() {
  local user=$(getArgument "user" $USER)
  # run at background without websocket
  # # 30 seconds timeout to check if the node is alive or not, the '&' symbol allows to run below commands while still having the process running
  # oraid start &
  # 30 seconds timeout  
  timeout 30 bash -c 'while [[ -z `wget --spider -S localhost:26657/health 2>&1 | grep "HTTP/.*200"` ]]; do sleep 1; done' || false

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
    gasPrices="0.000000000025orai"
  fi
  local securityContract=$(getArgument "security_contract" $SECURITY_CONTRACT)
  local identity=$(getArgument "identity" $IDENTITY)
  local website=$(getArgument "website" $WEBSITE)
  local details=$(getArgument "details" $DETAILS)

  echo "start creating validator..."
  requirePass
  sleep 5

  (echo "$PASS"; echo "$PASS") | oraid tx staking create-validator --amount $amount --pubkey $pubkey --moniker $moniker --chain-id Oraichain --commission-rate $commissionRate --commission-max-rate $commissionMaxRate --commission-max-change-rate $commissionMaxChangeRate --min-self-delegation $minDelegation --gas $gas --gas-prices $gasPrices --security-contact $securityContract --identity $identity --website $website --details $details --from $user -y
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
  start)
    oraidFn
  ;;
  initScript)
    initScriptFn
  ;;
  addGenAccount)
    addGenAccount  
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
