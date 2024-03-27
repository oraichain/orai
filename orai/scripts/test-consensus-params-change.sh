#!/bin/bash
set -ux

ARGS="--chain-id testing -y --keyring-backend test --fees 200orai --gas auto --gas-adjustment 1.5 -b block"
VALIDATOR_HOME=${VALIDATOR_HOME:-"$HOME/.oraid/validator1"}

# kill all running binaries
pkill oraid

# setup local network
sh $PWD/scripts/multinode-local-testnet.sh

time_iota_ms_before=$(curl "http://localhost:26657/consensus_params" | jq '.result.consensus_params.block.time_iota_ms | tonumber')
max_age_num_blocks_before=$(curl "http://localhost:26657/consensus_params" | jq '.result.consensus_params.evidence.max_age_num_blocks | tonumber')
max_bytes_before=$(curl "http://localhost:26657/consensus_params" | jq '.result.consensus_params.block.max_bytes | tonumber')

oraid tx gov submit-proposal param-change ../evidence.json --from validator1 --home "$HOME/.oraid/validator1" $ARGS
oraid tx gov submit-proposal param-change ../time_iota_ms.json --from validator1 --home "$HOME/.oraid/validator1" $ARGS
oraid tx gov vote 1 yes --from validator1 --home "$HOME/.oraid/validator1" $ARGS && oraid tx gov vote 1 yes --from validator2 --home "$HOME/.oraid/validator2" $ARGS
oraid tx gov vote 2 yes --from validator1 --home "$HOME/.oraid/validator1" $ARGS && oraid tx gov vote 2 yes --from validator2 --home "$HOME/.oraid/validator2" $ARGS

# sleep to wait til the proposal passes
echo "Sleep til the proposal passes..."
sleep 31

time_iota_ms=$(curl "http://localhost:26657/consensus_params" | jq '.result.consensus_params.block.time_iota_ms | tonumber')
max_age_num_blocks=$(curl "http://localhost:26657/consensus_params" | jq '.result.consensus_params.evidence.max_age_num_blocks | tonumber')
max_bytes=$(curl "http://localhost:26657/consensus_params" | jq '.result.consensus_params.block.max_bytes | tonumber')

if [[ $max_bytes_before == $max_bytes ]] ; then
   echo "Could not update max_bytes through param change" >&2; exit 1
fi

if [[ $max_age_num_blocks == $max_age_num_blocks_before ]] ; then
   echo "Could not update max_age_num_blocks through param change" >&2; exit 1
fi

if [[ $time_iota_ms == $time_iota_ms_before ]] ; then
   echo "Could not update time_iota_ms through param change" >&2; exit 1
fi

echo "Test Passed"