## Installation

```bash
docker-compose up -d && docker-compose exec orai ash
# wget https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a -O /lib/libwasmvm_muslc.a
make build
ln -s $PWD/build/oraid /usr/bin/oraid

# setup blockchain and run
./scripts/setup_oraid.sh 12345678
./scripts/run_oraid.sh
```

## Build smart contract and interact with it

```bash
# go to rust-optimizer container
docker-compose exec rust bash
cd play-smartc
optimize.sh .

# run unit-test
RUST_BACKTRACE=1 cargo unit-test -- --exact contract::tests::increment --show-output

# go to blockchain node container

# step1: store smart contract (will overide by smart contract name)
oraid tx wasm store smart-contracts/play-smartc/target/wasm32-unknown-unknown/release/play_smartc.wasm --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y

# step 2: get code id from response and instantiate it
oraid tx wasm instantiate $CODE_ID '{"count":10}' --from $USER --label "oracle 1" --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y

# if using genesis smart contract, the address will not changed if deployed on same block height
CONTRACT=$(oraid add-wasm-genesis-message list-contracts | jq '.[0].contract_address' -r)

# step 3: get smart contract address and query it
oraid query wasm list-contract-by-code $CODE_ID
oraid query wasm contract-state smart $CONTRACT '{"fetch":{"url":"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"}}'
# can change method to POST and optional authorization header
oraid query wasm contract-state smart $CONTRACT '{"fetch":{"url":"https://my-json-server.typicode.com/typicode/demo/posts","method":"POST"}}'

# step 4: test execute and query the updated state
oraid tx wasm execute $CONTRACT '{"increment":{}}' --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y
oraid query wasm contract-state smart $CONTRACT '{"get_count":{}}'

# step 5: migrate contract to a new code
# TODO: this is for oracle aggregation only in the future with latest version of wasmvm

```

## Some basic commands to test with the node

```bash

oraid tx provider set-datasource bitcoin_price $CONTRACT "test bitcoin price" --from duc --chain-id testing -y

oraid tx provider set-testcase bitcoin_price_testcase $CONTRACT "test bitcoin price testcase" --from duc --chain-id testing -y

oraid tx provider set-oscript oscript_btc $CONTRACT "test bitcoin price oracle script" --ds bitcoin_price --tc bitcoin_price_testcase --from duc --chain-id testing -y

curl -X POST -i http://localhost:1317/provider/datasource -d '{"name":"abc"}'

```


## Build protobuf and do lint check
```bash
docker-compose exec protoc ash

# check protobuf lint
make proto-lint

# build protobuf templates
PROTO_DIR=x/websocket/types/ make proto-gen
```