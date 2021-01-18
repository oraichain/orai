# Installation

```bash
docker-compose up -d && docker-compose exec orai ash
cp $PWD/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
make build-oraid
ln -s $PWD/build/wasmd /usr/bin/wasmd

# setup blockchain and run
./docker/setup_wasmd.sh
./docker/run_wasmd.sh
```

# Build smart contract and interact with it

```bash
# go to rust-optimizer container
docker-compose exec rust bash
cd play-smartc
optimize.sh .

# run unit-test
RUST_BACKTRACE=1 cargo unit-test -- --exact contract::tests::increment --show-output

# go to blockchain node container

# step1: store smart contract (will overide by smart contract name)
wasmd tx wasm store smart-contracts/play-smartc/target/wasm32-unknown-unknown/release/play_smartc.wasm --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y

# step 2: get code id from response and instantiate it
wasmd tx wasm instantiate $CODE_ID '{"count":10}' --from $USER --label "oracle 1" --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y

# step 3: get smart contract address and query it
wasmd query wasm list-contract-by-code $CODE_ID
wasmd query wasm contract-state smart $CONTRACT '{"fetch":{"url":"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"}}'

# step 4: test execute and query the updated state
wasmd tx wasm execute $CONTRACT '{"increment":{}}' --from $USER --gas="auto" --gas-adjustment="1.2" --chain-id=testing -y
wasmd query wasm contract-state smart $CONTRACT '{"get_count":{}}'
```
