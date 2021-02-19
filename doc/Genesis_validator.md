# Tutorial to participate in the genesis phase

### 1. Download the setup file.

```bash
curl -OL https://raw.githubusercontent.com/oraichain/oraichain-static-files/master/setup_genesis.sh
sudo chmod +x setup_genesis.sh
```

### 2. Run the setup file.

```bash
./setup_genesis.sh
```

### 3. Edit wallet name and moniker you prefer to create a new wallet and validator.

### 4. Run the following command:

```
docker-compose -f docker-compose.genesis.yml up -d
```

### 5. Enter the container through the command:

```
docker-compose exec orai ash
```

### 6. Type the following command to initiate your genesis node:

```bash
setup <your passphrase>
```

After running, there will be three .txt files generated.

One file, which is account.txt, stores your genesis account information as well as its mnemonic. Please keep it safe, and remove the file when you finish storing your account infromation.

Please upload the .oraid/config/gentx/gentx-\*.json file to the following google form:

### 8. Wait for the team to setup the genesis file

On the same vpn facilitating docker-machine, it can be done using command:
`docker-machine scp genesis1:"/home/ubuntu/oraichain/.oraid/config/gentx/gentx-*.json" genesis2:/home/ubuntu/oraichain/.oraid/config/gentx/`

then add genesis accounts:
`./fn.sh addGenAccount --address address --amount 10000000`

and finally call `oraid collect-gentxs` to complete the genesis.json file and use `docker-machine scp genesis1:/home/ubuntu/oraichain/.oraid/config/genesis.json genesis2:/home/ubuntu/oraichain/.oraid/config` to override other node's genesis file.

then following [Network](./network.md) for secured network setup.

### 9. After the team has finished setting up, type the following commands:

Download the new genesis file containing all the information of the genesis nodes.

```bash
curl $GENESIS_URL > .oraid/config/genesis.json
```

After downloading, please check if it contains your account and validator information. If it does not, please inform us so we can add your information.

### 10. Restart the container to start your node:

```
docker-compose -f docker-compose.genesis.yml restart orai
```

to check if the node has run successfully, you can make a simple http request as follows:

```bash
curl  -X GET http://localhost:1317/cosmos/staking/v1beta1/validators
```

if you see your validator information as well as others, then your node is running well

```json
{
  "validators": [
    {
      "operator_address": "oraivaloper13fw6fhmcnllp4c4u584rjsnuun2stddjgngk4y",
      "consensus_pubkey": {
        "@type": "/cosmos.crypto.ed25519.PubKey",
        "key": "B5zXxXtJ3fGOp9Ngxn5GtemEuX7JrAZL/ysayZSU2V4="
      },
      "jailed": false,
      "status": "BOND_STATUS_BONDED",
      "tokens": "250000000",
      "delegator_shares": "250000000.000000000000000000",
      "description": {
        "moniker": "phamthanhtu",
        "identity": "",
        "website": "",
        "security_contact": "",
        "details": ""
      },
      "unbonding_height": "0",
      "unbonding_time": "1970-01-01T00:00:00Z",
      "commission": {
        "commission_rates": {
          "rate": "0.100000000000000000",
          "max_rate": "0.200000000000000000",
          "max_change_rate": "0.010000000000000000"
        },
        "update_time": "2021-01-27T07:46:51.048265860Z"
      },
      "min_self_delegation": "1"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```
