## Installation

```bash
docker-compose up -d

# enter protoc container and generate the proto files
docker-compose exec protoc ash

# first time
go get ./...

# check protobuf lint
make proto-lint

# build protobuf templates
make proto-gen PROTO_DIR=x/provider/types/
make proto-gen PROTO_DIR=x/airequest/types/
make proto-gen PROTO_DIR=x/websocket/types/
make proto-gen PROTO_DIR=x/airesult/types/

# exit the container
exit

docker-compose exec orai ash

# wget https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a -O /lib/libwasmvm_muslc.a
make build GOMOD_FLAGS=
ln -s $PWD/build/oraid /usr/bin/oraid

# setup blockchain and run
./scripts/setup_oraid.sh 12345678

# start node
oraid start --rpc.laddr tcp://0.0.0.0:26657 --log_level error

# replace oraid with oraivisor for auto-upgrade
oraivisor start --rpc.laddr tcp://0.0.0.0:26657 --log_level error

# start websocket subscribe for processing event log in another terminal
oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=Oraichain -y

# run as a background process
docker-compose exec -d orai ash -c "echo $KEYRING_PASS | oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=Oraichain -y"
```

## Build smart contract and interact with it

```bash
# go to rust-optimizer container
docker-compose exec rust bash
cd play-smartc
optimize.sh .
# similarly, build other smart contracts
cd datasource-eth
optimize.sh .

# can run using simulate environment
docker-compose exec simulate bash
cosmwasm-simulate oscript-price/artifacts/oscript_price.wasm

# run unit-test
RUST_BACKTRACE=1 cargo unit-test -- --exact contract::tests::increment --show-output

# can using automated deployment
# if the smart contract has been stored using oraid tx wasm store, then use the below command with suitable code id
./scripts/deploy-contract.sh smart-contracts/testcase-price/artifacts/testcase_price.wasm "testcase-price 1" '{"ai_data_source":"datasource_eth","testcase":"testcase_price"}' [code_id]
# if not, then don't add the [code-id] field, it will give an error because the smart contract has not had a code id yet.

```

## Some basic commands to test with the node

```bash

./scripts/deploy-contract.sh smart-contracts/datasource-eth/artifacts/datasource_eth.wasm "datasource-eth 1" ''

oraid tx provider set-datasource eth_price $CONTRACT "test eth price" --from $USER --chain-id Oraichain -y --fees 5000orai

./scripts/deploy-contract.sh smart-contracts/testcase-price/artifacts/testcase_price.wasm "testcase-price 1" ''

oraid tx provider set-testcase eth_price_testcase $CONTRACT "test eth price testcase" --from $USER --chain-id Oraichain -y --fees 5000orai

./scripts/deploy-contract.sh smart-contracts/oscript-price/artifacts/oscript_price.wasm "oscript-price 1" '{"ai_data_source":"datasource_eth","testcase":"testcase_price"}'

oraid tx provider set-oscript oscript_eth $CONTRACT "test eth price oracle script" --ds eth_price --tc eth_price_testcase --from $USER --chain-id Oraichain -y

# open another terminal and run
oraid tx airequest set-aireq oscript_eth "5" "6" 30000orai 1 --from $USER --chain-id Oraichain -y

# Check if the AI request has finished executing
oraid query airesult fullreq <request-id>

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

## Development with oraivisor

```bash
ln -s /workspace/oraivisor/build/oraivisor /usr/bin/oraivisor
mkdir -p /workspace/.oraid/oraivisor/genesis/bin
ln -s /workspace/build/oraid /workspace/.oraid/oraivisor/genesis/bin/oraid
DAEMON_NAME=oraid DAEMON_HOME=/workspace/.oraid oraivisor start
```
