## Step by step guide to participate in Oraichain as a full node:

### 1) Remove files and directories from the previous configured node (if has):

```
rm -rf .oraid/

rm -rf .oraicli/

rm -rf .oraifiles/

# rm -rf .images/

rm go.sum

go get ./...

make install

```

### 2. Initiate your full node with your moniker and the chain id of the Oraichain network:

```
oraid init <moniker-name> --chain-id Oraichain
oraid init ducpl --chain-id Oraichain
```

### 3. Configure your CLI to eliminate need to declare them as flags
```
oraicli config chain-id Oraichain
oraicli config output json
oraicli config indent true
oraicli config trust-node true
```

### 4. We'll use the "test" keyring backend which save keys unencrypted in the configuration directory of your project (defaults to ~/.nsd). You should **never** use the "test" keyring backend in production. For more information about other options for keyring-backend take a look at https://docs.cosmos.network/master/interfaces/keyring.html
```
oraicli config keyring-backend test
```

### 5. Copy the genesis file of the network to .oraid/config/genesis.json

### 6. Add at least one persistent peer to retrieve all the blocks that have been committed from that peer.

```
persistent_peers = "<node-id>@<node-ip>:<node-port>"
persistent_peers = "91b3aa80c7900b9d7fea28d30d7c6e35fe0c9b02@157.230.22.169:26656"
```
In order to collect node-id: oraid tendermint show-node-id

### 7. Run the full node:

```
oraid start
```

### 8. Create a local account to interact with the network:

```
oraicli keys add <account-name>
oraicli keys add duc
```

### 9. Use a genesis accout to send some ORAI to the newly created account. This step must be done in order to broadcast the account onto the network.

```
oraicli tx send <sender address> <receiver address> amount --from <account-name> --fees <orai>
oraicli tx send orai1xrmqwy0ssk7uxz4p5egkvj3lct958qjyfn5raq orai1mmj588ap4z948gz5mxuwwju3nwsn9hkhep7fz7 100000000000orai --from ducbean --fees 5000orai
```

### 10. Create a new validator for your full node to receive ORAI tokens when committing a block:

```
oraicli tx staking create-validator --amount <orai> --pubkey <validator-pubkey> --moniker <your-moniker> --chain-id <network chain-id> --commission-rate <decimal> --commission-max-rate <decimal> --commission-max-change-rate <decimal> --min-self-delegation <integer> --gas <auto, integer> --gas-adjustment <decimal> --gas-prices <orai> --from duc
oraicli tx staking create-validator --amount 10000orai --pubkey oraivalconspub1addwnpepqvydmv22mkzc9rc92g43unew08cmj4q46dhk7vz0a9fj2xjsjn2lvqj0dfr --moniker ducphamle --chain-id Oraichain --commission-rate 0.10 --commission-max-rate 0.20 --commission-max-change-rate 0.01 --min-self-delegation 100 --gas auto --gas-adjustment 1.15 --gas-prices 0.025orai --from duc
```