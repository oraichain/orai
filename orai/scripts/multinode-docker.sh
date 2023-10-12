#!/bin/bash
set -ux

# # always returns true so set -e doesn't exit if it is not running.
rm -rf $PWD/.oraid/

# make four orai directories
mkdir $PWD/.oraid
mkdir $PWD/.oraid/validator1
mkdir $PWD/.oraid/validator2
mkdir $PWD/.oraid/validator3
docker volume rm orai_oraivisor
docker-compose -f docker-compose-e2e-upgrade.yml up -d --force-recreate
docker_command="docker-compose -f docker-compose-e2e-upgrade.yml exec"
validator1_command="$docker_command validator1 bash -c"
validator2_command="$docker_command validator2 bash -c"
validator3_command="$docker_command validator3 bash -c"
working_dir=/workspace/.oraid

# init all three validators
$validator1_command "oraid init --chain-id=testing validator1 --home=$working_dir/validator1"
$validator2_command "oraid init --chain-id=testing validator2 --home=$working_dir/validator2"
$validator3_command "oraid init --chain-id=testing validator3 --home=$working_dir/validator3"

# # create keys for all three validators
$validator1_command "oraid keys add validator1 --keyring-backend=test --home=$working_dir/validator1"
$validator1_command "oraid keys add validator2 --keyring-backend=test --home=$working_dir/validator2"
$validator1_command "oraid keys add validator3 --keyring-backend=test --home=$working_dir/validator3"

update_genesis () {    
    cat $PWD/.oraid/validator1/config/genesis.json | jq "$1" > $PWD/.oraid/validator1/config/tmp_genesis.json && mv $PWD/.oraid/validator1/config/tmp_genesis.json $PWD/.oraid/validator1/config/genesis.json
}

# change staking denom to orai
update_genesis '.app_state["staking"]["params"]["bond_denom"]="orai"'

# create validator node 1
validator1_key=`$validator1_command "oraid keys list --output json --home $working_dir/validator1 --keyring-backend test" | jq .[0].address`
$validator1_command "oraid add-genesis-account $validator1_key 1000000000000orai,1000000000000stake --home=$working_dir/validator1"
$validator1_command "oraid gentx validator1 500000000orai --keyring-backend=test --home=$working_dir/validator1 --chain-id=testing"
$validator1_command "oraid collect-gentxs --home=$working_dir/validator1"
$validator1_command "oraid validate-genesis --home=$working_dir/validator1"

# update staking genesis
update_genesis '.app_state["staking"]["params"]["unbonding_time"]="240s"'

# update crisis variable to orai
update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="orai"'

# udpate gov genesis
update_genesis '.app_state["gov"]["voting_params"]["voting_period"]="15s"'
update_genesis '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="orai"'

# update mint genesis
update_genesis '.app_state["mint"]["params"]["mint_denom"]="orai"'

# port key (validator1 uses default ports)
# validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060
# validator2 1316, 9088, 9089, 26655, 26654, 26653, 6061
# validator3 1315, 9086, 9087, 26652, 26651, 26650, 6062


# change app.toml values
VALIDATOR1_APP_TOML=$PWD/.oraid/validator1/config/app.toml
VALIDATOR2_APP_TOML=$PWD/.oraid/validator2/config/app.toml
VALIDATOR3_APP_TOML=$PWD/.oraid/validator3/config/app.toml

# change config.toml values
VALIDATOR1_CONFIG=$PWD/.oraid/validator1/config/config.toml
VALIDATOR2_CONFIG=$PWD/.oraid/validator2/config/config.toml
VALIDATOR3_CONFIG=$PWD/.oraid/validator3/config/config.toml

# Pruning - comment this configuration if you want to run upgrade script
pruning="custom"
pruning_keep_recent="5"
pruning_keep_every="10"
pruning_interval="10"

sed -i -e "s%^pruning *=.*%pruning = \"$pruning\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-keep-recent *=.*%pruning-keep-recent = \"$pruning_keep_recent\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-keep-every *=.*%pruning-keep-every = \"$pruning_keep_every\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^pruning-interval *=.*%pruning-interval = \"$pruning_interval\"%; " $VALIDATOR1_APP_TOML

sed -i -e "s%^pruning *=.*%pruning = \"$pruning\"%; " $VALIDATOR2_APP_TOML
sed -i -e "s%^pruning-keep-recent *=.*%pruning-keep-recent = \"$pruning_keep_recent\"%; " $VALIDATOR2_APP_TOML
sed -i -e "s%^pruning-keep-every *=.*%pruning-keep-every = \"$pruning_keep_every\"%; " $VALIDATOR2_APP_TOML
sed -i -e "s%^pruning-interval *=.*%pruning-interval = \"$pruning_interval\"%; " $VALIDATOR2_APP_TOML

sed -i -e "s%^pruning *=.*%pruning = \"$pruning\"%; " $VALIDATOR3_APP_TOML
sed -i -e "s%^pruning-keep-recent *=.*%pruning-keep-recent = \"$pruning_keep_recent\"%; " $VALIDATOR3_APP_TOML
sed -i -e "s%^pruning-keep-every *=.*%pruning-keep-every = \"$pruning_keep_every\"%; " $VALIDATOR3_APP_TOML
sed -i -e "s%^pruning-interval *=.*%pruning-interval = \"$pruning_interval\"%; " $VALIDATOR3_APP_TOML

# state sync  - comment this configuration if you want to run upgrade script
snapshot_interval="10"
snapshot_keep_recent="2"

sed -i -e "s%^snapshot-interval *=.*%snapshot-interval = \"$snapshot_interval\"%; " $VALIDATOR1_APP_TOML
sed -i -e "s%^snapshot-keep-recent *=.*%snapshot-keep-recent = \"$snapshot_keep_recent\"%; " $VALIDATOR1_APP_TOML

sed -i -e "s%^snapshot-interval *=.*%snapshot-interval = \"$snapshot_interval\"%; " $VALIDATOR2_APP_TOML
sed -i -e "s%^snapshot-keep-recent *=.*%snapshot-keep-recent = \"$snapshot_keep_recent\"%; " $VALIDATOR2_APP_TOML

sed -i -e "s%^snapshot-interval *=.*%snapshot-interval = \"$snapshot_interval\"%; " $VALIDATOR3_APP_TOML
sed -i -e "s%^snapshot-keep-recent *=.*%snapshot-keep-recent = \"$snapshot_keep_recent\"%; " $VALIDATOR3_APP_TOML

# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR1_CONFIG

# validator2
sed -i -E 's|tcp://127.0.0.1:26658|tcp://0.0.0.0:26658|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://0.0.0.0:26657|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26656|g' $VALIDATOR2_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR2_CONFIG

# validator3
sed -i -E 's|tcp://127.0.0.1:26658|tcp://0.0.0.0:26658|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://0.0.0.0:26657|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26656|g' $VALIDATOR3_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR3_CONFIG

# copy validator1 genesis file to validator2-3
cp $PWD/.oraid/validator1/config/genesis.json $PWD/.oraid/validator2/config/genesis.json
cp $PWD/.oraid/validator1/config/genesis.json $PWD/.oraid/validator3/config/genesis.json

validator1_id_res=`$validator1_command "oraid tendermint show-node-id --home=$working_dir/validator1 --log_format json"`
validator1_id=$(echo "$validator1_id_res" | tr -d -c '[:alnum:]') # remove all special characters because of the command's result
# copy tendermint node id of validator1 to persistent peers of validator2-3
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$validator1_id@validator1:26656\"|g" $VALIDATOR2_CONFIG
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$validator1_id@validator1:26656\"|g" $VALIDATOR3_CONFIG

# start all three validators
$docker_command -d validator1 bash -c "oraivisor start --home=$working_dir/validator1 --minimum-gas-prices=0.00001orai"
$docker_command -d validator2 bash -c "oraivisor start --home=$working_dir/validator2 --minimum-gas-prices=0.00001orai"
$docker_command -d validator3 bash -c "oraivisor start --home=$working_dir/validator3 --minimum-gas-prices=0.00001orai"

# send orai from first validator to second validator
echo "Waiting 7 seconds to send funds to validators 2 and 3..."
sleep 7

validator2_key_res=`$validator2_command "oraivisor keys show validator2 -a --keyring-backend=test --home=$working_dir/validator2"`
validator2_key=$(echo "$validator2_key_res" | tr -d -c '[:alnum:]') # remove all special characters because of the command's result
validator3_key_res=`$validator3_command "oraivisor keys show validator3 -a --keyring-backend=test --home=$working_dir/validator3"`
validator3_key=$(echo "$validator3_key_res" | tr -d -c '[:alnum:]') # remove all special characters because of the command's result

$validator1_command "oraid tx send validator1 $validator2_key 5000000000orai --keyring-backend=test --home=$working_dir/validator1 --chain-id=testing --broadcast-mode block --gas 200000 --fees 2orai --yes"
$validator1_command "oraid tx send validator1 $validator3_key 4000000000orai --keyring-backend=test --home=$working_dir/validator1 --chain-id=testing --broadcast-mode block --gas 200000 --fees 2orai --yes"
# send test orai to a test account
# oraid tx send $(oraid keys show validator1 -a --keyring-backend=test --home=$PWD/.oraid/validator1) orai14n3tx8s5ftzhlxvq0w5962v60vd82h30rha573 5000000000orai --keyring-backend=test --home=$PWD/.oraid/validator1 --chain-id=testing --broadcast-mode block --gas 200000 --fees 2orai --node http://localhost:26657 --yes

# create second & third validator
validator2_pubkey_res=`$validator2_command "oraid tendermint show-validator --home=$working_dir/validator2"`
validator2_pubkey=$(echo "$validator2_pubkey_res" | jq '@json') # remove all special characters because of the command's result
validator3_pubkey_res=`$validator3_command "oraid tendermint show-validator --home=$working_dir/validator3"`
validator3_pubkey=$(echo "$validator3_pubkey_res" | jq '@json') # remove all special characters because of the command's result
$validator2_command "oraid tx staking create-validator --amount=500000000orai --from=validator2 --pubkey=$validator2_pubkey --moniker=validator2 --chain-id=testing --commission-rate=0.1 --commission-max-rate=0.2 --commission-max-change-rate=0.05 --min-self-delegation=500000000 --keyring-backend=test --home=$working_dir/validator2 --broadcast-mode block --gas 200000 --fees 2orai --yes"
$validator3_command "oraid tx staking create-validator --amount=400000000orai --from=validator3 --pubkey=$validator3_pubkey --moniker=validator3 --chain-id=testing --commission-rate=0.1 --commission-max-rate=0.2 --commission-max-change-rate=0.05 --min-self-delegation=400000000 --keyring-backend=test --home=$working_dir/validator3 --broadcast-mode block --gas 200000 --fees 2orai --yes"

echo "All 3 Validators are up and running!"