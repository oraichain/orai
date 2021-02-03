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
oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=$CHAIN_ID -y
```

## Build smart contract and interact with it

```bash
# go to rust-optimizer container
docker-compose exec rust bash
cd play-smartc
optimize.sh .

# can run using simulate environment
docker-compose exec simulate bash
cosmwasm-simulate oscript-price/artifacts/oscript_price.wasm

# run unit-test
RUST_BACKTRACE=1 cargo unit-test -- --exact contract::tests::increment --show-output

# can using automated deployment
./scripts/deploy-contract.sh smart-contracts/testcase-price/artifacts/testcase_price.wasm "testcase-price 1" '{"ai_data_source":"datasource_eth","testcase":"testcase_price"}' [code_id]

```

## Some basic commands to test with the node

```bash

oraid tx provider set-datasource bitcoin_price $CONTRACT "test bitcoin price" --from duc --chain-id $CHAIN_ID -y

oraid tx provider set-testcase bitcoin_price_testcase $CONTRACT "test bitcoin price testcase" --from duc --chain-id $CHAIN_ID -y

oraid tx provider set-oscript oscript_btc $CONTRACT "test bitcoin price oracle script" --ds bitcoin_price --tc bitcoin_price_testcase --from duc --chain-id $CHAIN_ID -y

curl -X POST -i http://localhost:1317/provider/datasource -d '{"name":"abc"}'

```


## Build protobuf and do lint check
```bash
docker-compose exec protoc ash

# first time
go get ./...

# check protobuf lint
make proto-lint

# build protobuf templates
make proto-gen PROTO_DIR=x/websocket/types/
```

## Run test
`make test-method PACKAGE=github.com/oraichain/orai/x/airequest/keeper METHOD=TestCalucateMol`

## Build docker image
`docker build -t orai/orai:0.15-alpine -f Dockerfile.prod .`