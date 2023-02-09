# PRINT EVERY COMMAND
set -ux

rm -rf .oraid/

moniker="NODE_SYNC"

# make orai state sync directories
mkdir .oraid
mkdir .oraid/state_sync

SNAP_RPC2="http://0.0.0.0:26651"
SNAP_RPC1="http://0.0.0.0:26657"
SNAP_RPC="http://0.0.0.0:26654"
CHAIN_ID="testing"
TRUST_HEIGHT_RANGE=${TRUST_HEIGHT_RANGE:-10}

PEER_RPC_PORT=26657
PEER_P2P_PORT=26656
PEER_ID=$(curl --no-progress-meter http://0.0.0.0:$PEER_RPC_PORT/status | jq -r '.result.node_info.id')

echo "peer id: $PEER_ID"

# persistent_peers
PEER="$PEER_ID@0.0.0.0:$PEER_P2P_PORT"

# MAKE HOME FOLDER AND GET GENESIS
oraid init $moniker --chain-id $CHAIN_ID --home=.oraid/state_sync
cp ~/.oraid/validator1/config/genesis.json .oraid/state_sync/config
# cp -R ~/.oraid/validator1/wasm .oraid/state_sync/

# reset the node
oraid tendermint unsafe-reset-all --home=.oraid/state_sync

# change app.toml values
STATESYNC_APP_TOML=.oraid/state_sync/config/app.toml

# state_sync
sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1350|g' $STATESYNC_APP_TOML
sed -i -E 's|0.0.0.0:9090|0.0.0.0:9080|g' $STATESYNC_APP_TOML
sed -i -E 's|0.0.0.0:9091|0.0.0.0:9081|g' $STATESYNC_APP_TOML

# change config.toml values
STATESYNC_CONFIG=.oraid/state_sync/config/config.toml

# state sync node
sed -i -E 's|tcp://127.0.0.1:26658|tcp://0.0.0.0:26648|g' $STATESYNC_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://0.0.0.0:26647|g' $STATESYNC_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26643|g' $STATESYNC_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26640|g' $STATESYNC_CONFIG

sed -i -E 's|localhost:6060|localhost:6070|g' $STATESYNC_CONFIG

# Change config files (set the node name, add persistent peers, set indexer = "null")
sed -i -e "s%^moniker *=.*%moniker = \"$moniker\"%; " $STATESYNC_CONFIG
sed -i -e "s%^indexer *=.*%indexer = \"null\"%; " $STATESYNC_CONFIG

# GET TRUST HASH AND TRUST HEIGHT
LATEST_HEIGHT=$(curl -s $SNAP_RPC/block | jq -r .result.block.header.height); \
BLOCK_HEIGHT=$((LATEST_HEIGHT - $TRUST_HEIGHT_RANGE)); \
TRUST_HASH=$(curl -s "$SNAP_RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

# TELL USER WHAT WE ARE DOING
echo "LATEST HEIGHT: $LATEST_HEIGHT"
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"

sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true| ; \

s|^(addr_book_strict[[:space:]]+=[[:space:]]+).*$|\1false| ; \

s|^(persistent_peers[[:space:]]+=[[:space:]]+).*$|\1\"$PEER\"| ; \

s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"$SNAP_RPC,$SNAP_RPC1,$SNAP_RPC2\"| ; \

s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1$BLOCK_HEIGHT| ; \

s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"$TRUST_HASH\"| ; \

s|^(seeds[[:space:]]+=[[:space:]]+).*$|\1\"\"|" $STATESYNC_CONFIG

sed -i "s/allow_duplicate_ip *= *.*/allow_duplicate_ip = true/g" .oraid/state_sync/config/config.toml

echo "Waiting 1 seconds to start state sync"
sleep 1

# THERE, NOW IT'S SYNCED AND YOU CAN PLAY
screen -S state_sync -d -m oraid start --home=.oraid/state_sync --minimum-gas-prices=0.00001orai