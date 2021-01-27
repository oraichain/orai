## Installation

```bash
docker-compose up -d && docker-compose exec orai ash
# wget https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a -O /lib/libwasmvm_muslc.a
make build GOMOD_FLAGS=
ln -s $PWD/build/oraid /usr/bin/oraid

# setup blockchain and run
./scripts/setup_oraid.sh 12345678

# start node
oraid start --rpc.laddr tcp://0.0.0.0:26657 --log_level error

# start websocket subscribe for processing event log
oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y
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

# step 6: test testcase contract call datasource contract
# install contract and get the CODE_ID
oraid tx wasm store smart-contracts/datasource-price/artifacts/datasource_price.wasm --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y
TESTCASE_CONTRACT=$(oraid query wasm list-contract-by-code $CODE_ID | grep address | awk '{print $(NF)}')
# then query it with datasource contract address
oraid query wasm contract-state smart $TESTCASE_CONTRACT "{\"get_price\":{\"contract\":\"$CONTRACT\",\"token\":\"bitcoin\"}}"

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
make proto-gen PROTO_DIR=x/websocket/types/
```

## Build docker image
`docker build -t orai/orai:0.15-alpine -f Dockerfile.prod .`