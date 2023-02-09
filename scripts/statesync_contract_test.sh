#!/bin/bash
set -o errexit -o nounset -o pipefail -x

echo "-----------------------"
echo "## Add new wallet to state sync node"
oraid keys add alice --keyring-backend=test --home=.oraid/state_sync
oraid keys add bob --keyring-backend=test --home=.oraid/state_sync

echo "-----------------------"
echo "## Send fund to state sync account"
oraid tx send $(oraid keys show validator1 -a --keyring-backend=test --home=$HOME/.oraid/validator1) $(oraid keys show alice -a --keyring-backend=test --home=.oraid/state_sync) 500000orai --keyring-backend=test --home=$HOME/.oraid/validator1 --chain-id=testing --broadcast-mode block --gas 200000 --fees 2orai --node http://localhost:26657 --yes

oraid tx send $(oraid keys show validator1 -a --keyring-backend=test --home=$HOME/.oraid/validator1) $(oraid keys show bob -a --keyring-backend=test --home=.oraid/state_sync) 500000orai --keyring-backend=test --home=$HOME/.oraid/validator1 --chain-id=testing --broadcast-mode block --gas 200000 --fees 2orai --node http://localhost:26657 --yes

echo "-----------------------"
echo "## Add new CosmWasm contract"
RESP=$(oraid tx wasm store scripts/wasm_file/cw_nameservice-aarch64.wasm --from=validator1 --keyring-backend=test --home=$HOME/.oraid/validator1 --gas 1500000 --fees 150orai --chain-id="testing" -y --node=http://localhost:26657 -b block -o json)

CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
CODE_HASH=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-2].value')
echo "* Code id: $CODE_ID"
echo "* Code checksum: $CODE_HASH"

echo "-----------------------"
echo "## Create new contract instance"
INIT='{"purchase_price":{"amount":"100","denom":"orai"},"transfer_price":{"amount":"999","denom":"orai"}}'
TXFLAG=(--node "tcp://localhost:26647" --chain-id=testing --gas-prices 0.0001orai --gas auto --gas-adjustment 1.3)

oraid tx wasm instantiate "$CODE_ID" "$INIT" --from=alice --admin="$(oraid keys show alice -a --keyring-backend=test --home=.oraid/state_sync)" --keyring-backend=test --home=.oraid/state_sync --label "name service" $TXFLAG -y

CONTRACT=$(oraid query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT"

echo "### Query all"
RESP=$(oraid query wasm contract-state all "$CONTRACT" -o json)
echo "$RESP" | jq

echo "-----------------------"
echo "## Execute contract $CONTRACT"
# Register a name for the wallet address
REGISTER='{"register":{"name":"tony"}}'
oraid tx wasm execute $CONTRACT "$REGISTER" --amount 100orai --from=alice --keyring-backend=test --home=.oraid/state_sync $TXFLAG -y -b block -o json | jq

# Query the owner of the name record
NAME_QUERY='{"resolve_record": {"name": "tony"}}'
oraid query wasm contract-state smart $CONTRACT "$NAME_QUERY" --node "tcp://localhost:26647" --output json
# Owner is the alice's address

# Transfer the ownership of the name record to bob (change the "to" address to bob generated during environment setup)
# get bob's address
oraid keys show bob -a --keyring-backend=test --home=.oraid/state_sync
# bob's address is orai1gyjdry7uspdmpj4pvgx3rkn80tlaqnhy4d876h
TRANSFER='{"transfer":{"name":"tony","to":"orai1gyjdry7uspdmpj4pvgx3rkn80tlaqnhy4d876h"}}'
oraid tx wasm execute $CONTRACT "$TRANSFER" --amount 999orai --from=alice --keyring-backend=test --home=.oraid/state_sync $TXFLAG -y

# Query the record owner again to see the new owner address:
NAME_QUERY='{"resolve_record": {"name": "tony"}}'
oraid query wasm contract-state smart $CONTRACT "$NAME_QUERY" --node "tcp://localhost:26647" --output json
# Owner is the bob's address = orai1gyjdry7uspdmpj4pvgx3rkn80tlaqnhy4d876h

echo "-----------------------"
echo "## Set new admin"
echo "### Query old admin: $(oraid q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
echo "### Update contract"
oraid tx wasm set-contract-admin "$CONTRACT" "$(oraid keys show bob -a --keyring-backend=test --home=.oraid/state_sync)" \
  --from alice -y --keyring-backend=test --home=.oraid/state_sync --chain-id=testing --gas 200000 --fees 2orai -b block -o json | jq
echo "### Query new admin: $(oraid q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"

echo "-----------------------"
echo "## Migrate contract"
echo "### Upload new code"
RESP=$(oraid tx wasm store scripts/wasm_file//burner.wasm \
  --from=alice --keyring-backend=test --home=.oraid/state_sync --gas 1500000 --fees 150orai --chain-id="testing" -y -b block -o json)

BURNER_CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
echo "### Migrate to code id: $BURNER_CODE_ID"

DEST_ACCOUNT=$(oraid keys show alice -a --keyring-backend=test --home=.oraid/state_sync)
oraid tx wasm migrate "$CONTRACT" "$BURNER_CODE_ID" "{\"payout\": \"$DEST_ACCOUNT\"}" --from alice \
  --chain-id=testing --keyring-backend=test --home=.oraid/state_sync --gas 1500000 --fees 150orai -b block -y -o json | jq

echo "### Query destination account: $BURNER_CODE_ID"
oraid q bank balances "$DEST_ACCOUNT" -o json | jq

echo "### Query contract meta data: $CONTRACT"
oraid q wasm contract "$CONTRACT" -o json | jq

echo "### Query contract meta history: $CONTRACT"
oraid q wasm contract-history "$CONTRACT" -o json | jq
