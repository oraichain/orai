# Tutorial to participate in the genesis phase

### 1. Download the setup file.

```bash
https://raw.githubusercontent.com/oraichain/oraichain-static-files/master/setup.sh
```

### 2. Run the setup file.

```bash
sudo chmod -R +x setup.sh && ./setup.sh
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

The next two files are <moniker_name>_accounts.txt & <moniker_name>_validators.txt. Please copy the contents of these two files and paste them to the following google form: 

### 8. Wait for the team to connect to other nodes

### 9. After the team has finished setting up, type the following commands:

```bash
sed -i 's/persistent_peers *= *".*"/persistent_peers = "'"<list-of-private-ips-here>"'"/g' .oraid/config/config.toml 
```

```bash
curl $GENESIS_URL > .oraid/config/genesis.json
```

Remember to replace the <list-of-private-ips-here> by the input given by the team.

### 10. Enter the container and type the following command to start the node:

```
oraid start --rpc.laddr tcp://0.0.0.0:26657 --log_level info
```