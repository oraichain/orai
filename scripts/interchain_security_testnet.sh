#!/bin/bash
set -ux

PROV_NODE_MONIKER=coordinator
PROV_NODE_DIR=~/provider-coordinator
PROV_CHAIN_ID=provider

rm -rf $PROV_NODE_DIR
pkill -f oraid

mkdir -p $PROV_NODE_DIR
oraid init $PROV_NODE_MONIKER --chain-id $PROV_CHAIN_ID --home $PROV_NODE_DIR

jq ".app_state.gov.voting_params.voting_period = \"180s\"" \
    ${PROV_NODE_DIR}/config/genesis.json > ${PROV_NODE_DIR}/edited_genesis.json
mv ${PROV_NODE_DIR}/edited_genesis.json ${PROV_NODE_DIR}/config/genesis.json

PROV_KEY=provider-key
oraid keys add $PROV_KEY --home $PROV_NODE_DIR \
    --keyring-backend test --output json > ${PROV_NODE_DIR}/${PROV_KEY}.json 2>&1

# Get local account address
PROV_ACCOUNT_ADDR=$(jq -r .address ${PROV_NODE_DIR}/${PROV_KEY}.json)

# Add tokens
oraid add-genesis-account $PROV_ACCOUNT_ADDR 10000000000000orai \
    --keyring-backend test --home $PROV_NODE_DIR

oraid gentx $PROV_KEY 10000000000orai \
--keyring-backend test \
--moniker $PROV_NODE_MONIKER \
--chain-id $PROV_CHAIN_ID \
--home $PROV_NODE_DIR

oraid collect-gentxs --home $PROV_NODE_DIR \
    --gentx-dir ${PROV_NODE_DIR}/config/gentx/

MY_IP="0.0.0.0"
sed -i -r "/node =/ s/= .*/= \"tcp:\/\/${MY_IP}:26658\"/" \
    ${PROV_NODE_DIR}/config/client.toml

oraid start --home $PROV_NODE_DIR \
        --rpc.laddr tcp://${MY_IP}:26658 \
        --grpc.address ${MY_IP}:9091 \
        --address tcp://${MY_IP}:26655 \
        --p2p.laddr tcp://${MY_IP}:26656 \
        --grpc-web.enable=false \
        &> ${PROV_NODE_DIR}/logs &

tail -f ${PROV_NODE_DIR}/logs

oraid q staking validators --home $PROV_NODE_DIR

tee ${PROV_NODE_DIR}/consumer-proposal.json<<EOF
{
    "title": "Propose the addition of a new chain",
    "description": "Gonna be a great chain",
    "chain_id": "consumer",
    "initial_height": {
        "revision_height": 1
    },
    "genesis_hash": "Z2VuX2hhc2g=",
    "binary_hash": "YmluX2hhc2g=",
    "spawn_time": "2023-03-17T18:03:04.292959+07:00",
    "consumer_redistribution_fraction": "0.75",
    "blocks_per_distribution_transmission": 1000,
    "historical_entries": 10000,
    "ccv_timeout_period": 2419200000000000,
    "transfer_timeout_period": 3600000000000,
    "unbonding_period": 1728000000000000,
    "deposit": "10000001orai"
}
EOF



#create proposal
oraid tx gov submit-proposal \
       consumer-addition ${PROV_NODE_DIR}/consumer-proposal.json \
       --keyring-backend test \
       --chain-id $PROV_CHAIN_ID \
       --from $PROV_KEY \
       --home $PROV_NODE_DIR \
       -b block

#vote yes
oraid tx gov vote 1  yes --from $PROV_KEY \
       --keyring-backend test --chain-id $PROV_CHAIN_ID --home $PROV_NODE_DIR -b block

#Verify that the proposal status is now `PROPOSAL_STATUS_PASSED`
oraid q gov proposal 1 --home $PROV_NODE_DIR

# Consumer chain
CONS_NODE_DIR=~/consumer-coordinator
rm -rf $CONS_NODE_DIR
pkill -f interchain-security-cd

mkdir -p $CONS_NODE_DIR
CONS_NODE_MONIKER=coordinator
CONS_CHAIN_ID=consumer

interchain-security-cd init $CONS_NODE_MONIKER --chain-id $CONS_CHAIN_ID --home $CONS_NODE_DIR

CONS_KEY=consumer-key
interchain-security-cd keys add $CONS_KEY --home $CONS_NODE_DIR \
    --keyring-backend test --output json > ${CONS_NODE_DIR}/${CONS_KEY}.json 2>&1

#Get local account address
CONS_ACCOUNT_ADDR=$(jq -r .address ${CONS_NODE_DIR}/${CONS_KEY}.json)

#Add account address to genesis
interchain-security-cd add-genesis-account $CONS_ACCOUNT_ADDR 1000000000orai \
    --keyring-backend test --home $CONS_NODE_DIR

oraid query provider consumer-genesis $CONS_CHAIN_ID --home $PROV_NODE_DIR -o json > ccvconsumer_genesis.json

jq -s '.[0].app_state.ccvconsumer = .[1] | .[0]' ${CONS_NODE_DIR}/config/genesis.json ccvconsumer_genesis.json > \
      ${CONS_NODE_DIR}/edited_genesis.json 

mv ${CONS_NODE_DIR}/edited_genesis.json ${CONS_NODE_DIR}/config/genesis.json &&
    rm ccvconsumer_genesis.json

echo '{"height": "0","round": 0,"step": 0}' > ${CONS_NODE_DIR}/data/priv_validator_state.json  
cp ${PROV_NODE_DIR}/config/priv_validator_key.json ${CONS_NODE_DIR}/config/priv_validator_key.json  
cp ${PROV_NODE_DIR}/config/node_key.json ${CONS_NODE_DIR}/config/node_key.json

sed -i -r "/node =/ s/= .*/= \"tcp:\/\/${MY_IP}:26648\"/" ${CONS_NODE_DIR}/config/client.toml

# consumer local node use the following command
interchain-security-cd start --home $CONS_NODE_DIR \
        --rpc.laddr tcp://${MY_IP}:26648 \
        --grpc.address ${MY_IP}:9081 \
        --address tcp://${MY_IP}:26645 \
        --p2p.laddr tcp://${MY_IP}:26646 \
        --grpc-web.enable=false \
        &> ${CONS_NODE_DIR}/logs &

tail -f ${CONS_NODE_DIR}/logs

tee ~/.hermes/config.toml<<EOF
[global]
 log_level = "info"

[[chains]]
account_prefix = "cosmos"
clock_drift = "5s"
gas_multiplier = 1.1
grpc_addr = "tcp://${MY_IP}:9081"
id = "$CONS_CHAIN_ID"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://${MY_IP}:26648"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "14days"
websocket_addr = "ws://${MY_IP}:26648/websocket"

[chains.gas_price]
       denom = "orai"
       price = 0.00

[chains.trust_threshold]
       denominator = "3"
       numerator = "1"

[[chains]]
account_prefix = "orai"
clock_drift = "5s"
gas_multiplier = 1.1
grpc_addr = "tcp://${MY_IP}:9091"
id = "$PROV_CHAIN_ID"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://${MY_IP}:26658"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "14days"
websocket_addr = "ws://${MY_IP}:26658/websocket"

[chains.gas_price]
       denom = "orai"
       price = 0.00

[chains.trust_threshold]
       denominator = "3"
       numerator = "1"
EOF

#Delete all previous keys in relayer
hermes keys delete --chain consumer --all
hermes keys delete --chain provider --all

#Import accounts key
hermes keys add --key-file  ${CONS_NODE_DIR}/${CONS_KEY}.json --chain consumer
hermes keys add --key-file  ${PROV_NODE_DIR}/${PROV_KEY}.json --chain provider

hermes create connection \
     --a-chain consumer \
    --a-client 07-tendermint-0 \
    --b-client 07-tendermint-0

hermes create channel \
    --a-chain consumer \
    --a-port consumer \
    --b-port provider \
    --order ordered \
    --channel-version 1 \
    --a-connection connection-0

pkill -f hermes  
hermes --json start &> ~/.hermes/logs &

tail -f ~/.hermes/logs


# Get validator delegations
DELEGATIONS=$(oraid q staking delegations \
  $(jq -r .address ${PROV_NODE_DIR}/${PROV_KEY}.json) --home $PROV_NODE_DIR -o json)

# Get validator operator address
OPERATOR_ADDR=$(echo $DELEGATIONS | jq -r '.delegation_responses[0].delegation.validator_address')

# Delegate tokens
oraid tx staking delegate $OPERATOR_ADDR 6900000orai \
                --from $PROV_KEY \
                --keyring-backend test \
                --home $PROV_NODE_DIR \
                --chain-id $PROV_CHAIN_ID \
                -y -b block

# Query provider chain valset
oraid q tendermint-validator-set --home $PROV_NODE_DIR
    
# Query consumer chain valset    
interchain-security-cd q tendermint-validator-set --home $CONS_NODE_DIR