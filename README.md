## Installation

```bash
docker-compose up -d

# enter protoc container and generate the proto files
docker-compose exec protoc ash

# first time
go get ./...

# build protobuf templates
make proto-gen

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
oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.5" --chain-id=$CHAIN_ID -y

# run as a background process
docker-compose exec -d orai ash -c "echo $KEYRING_PASS | oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.5" --chain-id=$CHAIN_ID -y"
```

## Build smart contract and interact with it

```bash
# go to rust-optimizer container
docker-compose exec rust bash
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
./scripts/deploy-contract.sh smart-contracts/testcase-price/artifacts/testcase_price.wasm "testcase-price 1" '{"ai_data_source":["datasource_eth"],"testcase":["testcase_price"]}' [code_id]
# if not, then don't add the [code-id] field, it will give an error because the smart contract has not had a code id yet.

# query a data source through cli
oraid query wasm contract-state smart $CONTRACT '{"get":{"input":"{\"image\":\"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSfx__RoRYzLDgXDiJxYGxLihJC4zoqV3V0xg&usqp=CAU\",\"model\":\"inception_v3\",\"name\":\"test_image\"}"}}'

# query wasm through lcd
curl <url>/wasm/v1beta1/contract/<contract-address>/smart/<json-string-encoded-in-base64>

oraid query wasm contract-state smart orai16at0lzgx3slnqlgvcc7r79056f5wkuczenn09k '{"test":{"input":"{\"image\":\"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSfx__RoRYzLDgXDiJxYGxLihJC4zoqV3V0xg&usqp=CAU\",\"model\":\"inception_v3\",\"name\":\"test_image\"}","output":"a","contract":"orai1aysde07zjurpp99jgl4xa7vskr8xnlcfkedkd9"}}'

```

## Some basic commands to test with the node

Run websocket as background process

```bash
echo <passphrase> | oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.5" --chain-id=$CHAIN_ID -y
```

Init smart contracts and create an AI request. To run the script, your current dir must contain the smart-contracts/ directory that already have wasm files built. The directory name with the wasm file should also be consistent. Eg: dir name: classification, then the wasm file is classification.wasm

```bash

./scripts/deploy_ai_services.sh <list-of-datasource-names> <list-of-testcase-names> <oscript-name> <datasource-init-input> <testcase-input> <script-indexing> <path to the oraiwasm directory> <passphrase>

Eg: ./scripts/deploy_ai_services.sh classification,cv009 classification_testcase classification_oscript '' '' '{"ai_data_source":["classification","cv009"],"testcase":["classification_testcase"]}' 1 /workspace/oraiwasm 123456789

# open another terminal and run
oraid tx airequest set-aireq oscript_price "5" "6" 30000orai 1 --from $USER --chain-id $CHAIN_ID -y

# interact with the AI services 
oraid tx airequest set-aireq classification_oscript '{"image":"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSfx__RoRYzLDgXDiJxYGxLihJC4zoqV3V0xg&usqp=CAU","model":"inception_v3","name":"test_image"}' "6" 30000orai 1 --from $USER --chain-id $CHAIN_ID -y

# Check if the AI request has finished executing
oraid query airesult fullreq <request-id>

```

Most of the time, the initial inputs for data sources and test cases are unecessary. However, you must set the input json for the oracle script with data source and test case information.

## Run test
`make test-method PACKAGE=github.com/oraichain/orai/x/airequest/keeper METHOD=TestCalucateMol`

## Build docker image

development `docker build -t orai/orai:alpine-wasm .`  
production `docker build -t orai/orai:0.18-alpine -f Dockerfile.prod .`  
oraivisor-upgrade `docker build -t orai/orai:mainnet-alpine-0.1 -f Dockerfile.oraivisor .`  

## Development with oraivisor

```bash
ln -s /workspace/oraivisor/build/oraivisor /usr/bin/oraivisor
mkdir -p /oraivisor/genesis/bin
ln -s /workspace/build/oraid /oraivisor/genesis/bin/oraid
DAEMON_NAME=oraid DAEMON_HOME=/ oraivisor start
```

## Create swagger documentation

```bash
# go to proto
docker-compose exec proto bash
make proto-swagger
# then create static file
go get github.com/rakyll/statik
statik -src doc/swagger-ui/ -dest doc -f
```