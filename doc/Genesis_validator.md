# Tutorial to participate in the genesis phase

### 1. Download the setup file.

```bash
https://raw.githubusercontent.com/oraichain/oraichain-static-files/master/setup.sh
```

### 2. Run the setup file.

```bash
sudo chmod -R +x setup.sh && ./setup.sh
```

### 3. Create a new wallet account on the explorer.

### 4. Edit your orai.env accordingly to your wallet account.

### 5. Run the following command:

```
docker-compose -f docker-compose.genesis.yml up -d
```

# TODO: How to init the node without stoping it ?

Wait for about 1 minutes then continue the next steps

### 7. Wait for the team to connect to other nodes

### 8. After the team has finished setting up, type the following command:

```bash
sed -i 's/persistent_peers *= *".*"/persistent_peers = "'"<list-of-private-ips-here>"'"/g' .oraid/config/config.toml 
```

Remember to replace the <list-of-private-ips-here> by the input given by the team.

### 9. Enter the container and type the following command to start the node:

```
oraid start --rpc.laddr tcp://0.0.0.0:26657 --log_level info
```