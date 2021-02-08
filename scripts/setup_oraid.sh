#!/bin/sh
set -o errexit -o nounset -o pipefail

echo -n "Enter passphrase:"
read -s PASSWORD

CHAIN_ID=${CHAIN_ID:-Oraichain}
USER=${USER:-tupt}
MONIKER=${MONIKER:-node001}
# PASSWORD=${PASSWORD:-$1}
rm -rf .oraid/

oraid init --chain-id "$CHAIN_ID" "$MONIKER"

(echo "$PASSWORD"; echo "$PASSWORD") | oraid keys add $USER

# hardcode the validator account for this instance
(echo "$PASSWORD") | oraid add-genesis-account $USER "100000000000000orai"

# (optionally) add smart contract
CONTRACT_CODE=${CONTRACT_CODE:-smart-contracts/play-smartc/target/wasm32-unknown-unknown/release/play_smartc.wasm}

if [ -f $CONTRACT_CODE ];then     
  echo "## Genesis CosmWasm contract"
  (echo "$PASSWORD") | oraid add-wasm-genesis-message store $CONTRACT_CODE --instantiate-everybody false --run-as $USER

  echo "-----------------------"
  echo "## Genesis CosmWasm instance"
  INIT='{"count":10}'
  BASE_ACCOUNT=$(echo "$PASSWORD" | oraid keys show $USER -a)
  (echo "$PASSWORD") | oraid add-wasm-genesis-message instantiate-contract 1 $INIT --run-as $USER --label=oracle --amount=100orai --admin $BASE_ACCOUNT

  # if need execute
  CONTRACT=$(oraid add-wasm-genesis-message list-contracts | jq '.[0].contract_address' -r)
  echo "-----------------------"
  echo "## List Genesis CosmWasm codes"
  oraid add-wasm-genesis-message list-codes

  echo "-----------------------"
  echo "## List Genesis CosmWasm contracts"
  oraid add-wasm-genesis-message list-contracts
fi

# submit a genesis validator tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
(echo "$PASSWORD"; echo "$PASSWORD") | oraid gentx $USER "250000000orai" --chain-id="$CHAIN_ID" --amount="250000000orai" -y

oraid collect-gentxs


