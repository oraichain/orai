#!/bin/sh
#set -o errexit -o nounset -o pipefail

PASSWORD=${PASSWORD:-12345678}
STAKE=${STAKE_TOKEN:-ustake}
FEE=${FEE_TOKEN:-ucosm}
CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
MONIKER=${MONIKER:-node001}

rm -rf "$HOME"/.wasmd

wasmd init --chain-id "$CHAIN_ID" "$MONIKER"
# staking/governance token is hardcoded in config, change this
sed -i "s/\"stake\"/\"$STAKE\"/" "$HOME"/.wasmd/config/genesis.json

(echo "$PASSWORD"; echo "$PASSWORD") | wasmd keys add $USER

# hardcode the validator account for this instance
echo "$PASSWORD" | wasmd add-genesis-account $USER "1000000000$STAKE,1000000000$FEE"

# (optionally) add a few more genesis accounts
for addr in "$@"; do
  echo $addr
  wasmd add-genesis-account "$addr" "1000000000$STAKE,1000000000$FEE"
done

# (optionally) add smart contract
if [ ! -z $LOCAL ];then     
  echo "## Genesis CosmWasm contract"
  wasmd add-wasm-genesis-message store x/wasm/internal/keeper/testdata/play_smartc.wasm --instantiate-everybody false --run-as $USER

  echo "-----------------------"
  echo "## Genesis CosmWasm instance"
  INIT='{"count":10}'
  BASE_ACCOUNT=$(wasmd keys show $USER -a)
  wasmd add-wasm-genesis-message instantiate-contract 1 $INIT --run-as $USER --label=oracle --amount=100ustake --admin $BASE_ACCOUNT

  # if need execute
  CONTRACT=$(wasmd add-wasm-genesis-message list-contracts | jq '.[0].contract_address' -r)
  echo "-----------------------"
  echo "## List Genesis CosmWasm codes"
  wasmd add-wasm-genesis-message list-codes

  echo "-----------------------"
  echo "## List Genesis CosmWasm contracts"
  wasmd add-wasm-genesis-message list-contracts
fi

# submit a genesis validator tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
(echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | wasmd gentx $USER "250000000$STAKE" --chain-id="$CHAIN_ID" --amount="250000000$STAKE"
## should be:
# (echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | wasmd gentx validator "250000000$STAKE" --chain-id="$CHAIN_ID"
wasmd collect-gentxs


