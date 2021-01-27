#!/bin/sh
#set -o errexit -o nounset -o pipefail

pass="$1"
if [ ! -z $pass ] 
then   
  # 3 times send passphrase
  expect << EOF
    set timeout $timeout
    spawn $0
    expect {
        "*passphrase:" { send -- "$pass\r" }
    }
    expect {
        "*passphrase:" { send -- "$pass\r" }
    }
    expect {
        "*passphrase:" { send -- "$pass\r" }
    }
    expect {
        "*passphrase:" { send -- "$pass\r" }
    }
    expect {
        "*passphrase:" { send -- "$pass\r" }
    }
    expect {
        "*passphrase:" { send -- "$pass\r" }
    }
    expect {
        "*passphrase:" { send -- "$pass\r" }
    }
    expect eof
EOF

  exit 0
fi 

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
MONIKER=${MONIKER:-node001}

rm -rf "$HOME"/.oraid

oraid init --chain-id "$CHAIN_ID" "$MONIKER"
# staking/governance token is hardcoded in config, change this
sed -i "s/\"stake\"/\"$STAKE\"/" "$HOME"/.oraid/config/genesis.json

oraid keys add $USER

# hardcode the validator account for this instance
oraid add-genesis-account $USER "100000000000000orai"

# (optionally) add a few more genesis accounts
for addr in "$@"; do
  echo $addr
  oraid add-genesis-account "$addr" "1000000000orai"
done

# (optionally) add smart contract
CONTRACT_CODE=${CONTRACT_CODE:-smart-contracts/play-smartc/target/wasm32-unknown-unknown/release/play_smartc.wasm}
echo $CONTRACT_CODE
if [ -f $CONTRACT_CODE ];then     
  echo "## Genesis CosmWasm contract"
  oraid add-wasm-genesis-message store $CONTRACT_CODE --instantiate-everybody false --run-as $USER

  echo "-----------------------"
  echo "## Genesis CosmWasm instance"
  INIT='{"count":10}'
  BASE_ACCOUNT=$(oraid keys show $USER -a)
  oraid add-wasm-genesis-message instantiate-contract 1 $INIT --run-as $USER --label=oracle --amount=100orai --admin $BASE_ACCOUNT

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
oraid gentx $USER "250000000orai" --chain-id="$CHAIN_ID" --amount="250000000orai"

oraid collect-gentxs


