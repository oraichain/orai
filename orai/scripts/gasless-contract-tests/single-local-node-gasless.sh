#!/bin/bash

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
MONIKER=${MONIKER:-node001}
WASM_PATH=${WASM_PATH:-"$PWD/scripts/wasm_file/swapmap.wasm"}
EXECUTE_MSG=${EXECUTE_MSG:-'{"ping":{}}'}
ARGS="--chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block"

pkill oraid
rm -rf .oraid/

update_genesis () {    
    cat $PWD/.oraid/config/genesis.json | jq "$1" > $PWD/.oraid/config/tmp_genesis.json && mv $PWD/.oraid/config/tmp_genesis.json $PWD/.oraid/config/genesis.json
}

oraid init --chain-id "$CHAIN_ID" "$MONIKER"

# 2s for fast test
update_genesis '.app_state["gov"]["voting_params"]["voting_period"]="2s"'

oraid keys add $USER --keyring-backend test 2>&1 | tee account.txt
oraid keys add $USER-eth --keyring-backend test --eth 2>&1 | tee account-eth.txt
oraid keys unsafe-export-eth-key $USER-eth --keyring-backend test 2>&1 | tee priv-eth.txt

# hardcode the validator account for this instance
oraid add-genesis-account $USER "100000000000000orai" --keyring-backend test
oraid add-genesis-account $USER-eth "100000000000000orai" --keyring-backend test
oraid add-genesis-account orai1kzkf6gttxqar9yrkxfe34ye4vg5v4m588ew7c9 "100000000000000orai" --keyring-backend test

# submit a genesis validator tx
oraid gentx $USER "250000000orai" --chain-id="$CHAIN_ID" --amount="250000000orai" -y --keyring-backend test

oraid collect-gentxs

screen -S test-gasless -d -m oraid start --json-rpc.address="0.0.0.0:8545" --json-rpc.ws-address="0.0.0.0:8546" --json-rpc.api="eth,web3,net,txpool,debug" --json-rpc.enable

# wait for the node to start
sleep 2

# prepare a new contract for gasless
store_ret=$(oraid tx wasm store $WASM_PATH --from $USER $ARGS --output json)
code_id=$(echo $store_ret | jq -r '.logs[0].events[1].attributes[] | select(.key | contains("code_id")).value')
oraid tx wasm instantiate $code_id '{}' --label 'testing' --from $USER --admin $(oraid keys show $USER --keyring-backend test -a) $ARGS
contract_address=$(oraid query wasm list-contract-by-code $code_id --output json | jq -r '.contracts[0]')
echo $contract_address

# try executing something, gas should equal 0
gas_used_before=$(oraid tx wasm execute $contract_address $EXECUTE_MSG --from $USER $ARGS --output json --gas 20000000 | jq '.gas_used | tonumber')
echo "gas used before gasless: $gas_used_before"

# set gasless proposal
oraid tx gov submit-proposal set-gasless $contract_address --title "gasless" --description "gasless" --deposit 10000000orai --from $USER $ARGS
oraid tx gov vote 1 yes --from $USER $ARGS

# proposal takes 2s
sleep 3
proposal_status=$(oraid query gov proposal 1 --output json | jq .status)
if ! [[ $proposal_status =~ "PROPOSAL_STATUS_PASSED" ]] ; then
   echo "The proposal has not passed yet"; exit 1
fi

# try executing something, gas should equal 0
gas_used_after=$(oraid tx wasm execute $contract_address $EXECUTE_MSG --from $USER $ARGS --output json --gas 20000000 | jq '.gas_used | tonumber')
echo "gas used after gasless: $gas_used_after"

# 1.9 is a magic number chosen to check that if the gas used after gasless has dropped significantly or not
gas_used_compare=$(echo "$gas_used_before / 1.9" | bc -l)
if [[ $gas_used_compare < $gas_used_after ]] ; then
   echo "Gas used after is not small enough!"; exit 1
fi