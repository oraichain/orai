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

# query a data source through cli
oraid query wasm contract-state smart $CONTRACT '{"get":{"input":""}}'

# query wasm through lcd
curl <url>/wasm/v1beta1/contract/<contract-address>/smart/<json-string-encoded-in-base64>

oraid query wasm contract-state smart orai16at0lzgx3slnqlgvcc7r79056f5wkuczenn09k '{"test":{"input":"{\"image\":\"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSfx__RoRYzLDgXDiJxYGxLihJC4zoqV3V0xg&usqp=CAU\",\"model\":\"inception_v3\",\"name\":\"test_image\"}","output":"a","contract":"orai1aysde07zjurpp99jgl4xa7vskr8xnlcfkedkd9"}}'

```

## Some basic commands to test with the node

Run websocket as background process

```bash
echo <passphrase> | oraid tx websocket subscribe --max-try 10 --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=Oraichain -y
```

Init smart contracts and create an AI request

```bash

./scripts/basic.sh <passphrase>

# open another terminal and run
oraid tx airequest set-aireq oscript_eth "5" "6" 30000orai 1 --from $USER --chain-id Oraichain -y

# interact with the AI services 
oraid tx airequest set-aireq classification_oscript '{"image":"https://encrypted-tbn0.gstati
c.com/images?q=tbn:ANd9GcSfx__RoRYzLDgXDiJxYGxLihJC4zoqV3V0xg&usqp=CAU","model":"inception_v3","name":"te
st_image"}' "6" 30000orai 1 --from $USER --chain-id Oraichain -y

# Check if the AI request has finished executing
oraid query airesult fullreq <request-id>

```

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