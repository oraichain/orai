# PRINT EVERY COMMAND
set -ux

rm -rf .oraid/

moniker="NODE_SYNC"

# make orai state sync directories
mkdir .oraid
mkdir .oraid/state_sync

SNAP_RPC2="http://127.0.0.1:26651"
SNAP_RPC1="http://127.0.0.1:26657"
SNAP_RPC="http://127.0.0.1:26654"
CHAIN_ID="testing"

# persistent_peers
PEER="0a89582c5138a51a23f46175e42bb4a7312f2e7e@127.0.0.1:26653"

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
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26648|g' $STATESYNC_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26647|g' $STATESYNC_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26643|g' $STATESYNC_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26640|g' $STATESYNC_CONFIG

sed -i -E 's|localhost:6060|localhost:6070|g' $STATESYNC_CONFIG

# Change config files (set the node name, add persistent peers, set indexer = "null")
sed -i -e "s%^moniker *=.*%moniker = \"$moniker\"%; " $STATESYNC_CONFIG
sed -i -e "s%^indexer *=.*%indexer = \"null\"%; " $STATESYNC_CONFIG

# GET TRUST HASH AND TRUST HEIGHT
LATEST_HEIGHT=$(curl -s $SNAP_RPC/block | jq -r .result.block.header.height); \
BLOCK_HEIGHT=$((LATEST_HEIGHT - 2000)); \
TRUST_HASH=$(curl -s "$SNAP_RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)

# pruning
pruning="custom"
pruning_keep_recent="5"
pruning_keep_every="10"
pruning_interval="10"

sed -i -e "s%^pruning *=.*%pruning = \"$pruning\"%; " $STATESYNC_APP_TOML
sed -i -e "s%^pruning-keep-recent *=.*%pruning-keep-recent = \"$pruning_keep_recent\"%; " $STATESYNC_APP_TOML
sed -i -e "s%^pruning-keep-every *=.*%pruning-keep-every = \"$pruning_keep_every\"%; " $STATESYNC_APP_TOML
sed -i -e "s%^pruning-interval *=.*%pruning-interval = \"$pruning_interval\"%; " $STATESYNC_APP_TOML

# state sync
snapshot_interval="10"
snapshot_keep_recent="2"

sed -i -e "s%^snapshot-interval *=.*%snapshot-interval = \"$snapshot_interval\"%; " $STATESYNC_APP_TOML
sed -i -e "s%^snapshot-keep-recent *=.*%snapshot-keep-recent = \"$snapshot_keep_recent\"%; " $STATESYNC_APP_TOML


# TELL USER WHAT WE ARE DOING
echo "LATEST HEIGHT: $LATEST_HEIGHT"
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"

sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true| ; \

s|^(addr_book_strict[[:space:]]+=[[:space:]]+).*$|\1false| ; \

s|^(persistent_peers[[:space:]]+=[[:space:]]+).*$|\1\"$PEER,$PEER\"| ; \

s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"$SNAP_RPC,$SNAP_RPC1,$SNAP_RPC2\"| ; \

s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1$BLOCK_HEIGHT| ; \

s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"$TRUST_HASH\"| ; \

s|^(seeds[[:space:]]+=[[:space:]]+).*$|\1\"\"|" $STATESYNC_CONFIG

echo "Waiting 1 seconds to start state sync"
sleep 1

# THERE, NOW IT'S SYNCED AND YOU CAN PLAY
screen -S state_sync -d -m oraid start --home=.oraid/state_sync --minimum-gas-prices=0.00001orai