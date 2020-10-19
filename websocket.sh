#!/bin/bash

# ./websocket.sh duc 5

rm -rf ~/.websocket

websocket keys delete-all

# config chain id
websocket config chain-id Oraichain

# add validator to websocket config
websocket config validator $(oraicli keys show $1 -a --bech val --keyring-backend test)

# setup broadcast-timeout to websocket config
websocket config broadcast-timeout "30s"

# setup rpc-poll-interval to yoda config
websocket config rpc-poll-interval "1s"

# setup max-try to yoda config
websocket config max-try 5

# config log type
websocket config log-level debug

# echo "y" | oraicli tx oracle activate --from $1 --keyring-backend test

# wait for activation transaction success
sleep 2

# for i in $(eval echo {1..$2})
# do
  # add reporter key
  websocket keys add $2
# done

sleep 2

# send orai tokens to reporters
echo "y" | oraicli tx send $(oraicli keys show $1 -a) $(websocket keys show $2) 10000000orai --from $(oraicli keys show $1 -a) --fees 5000orai

sleep 6

#wait for sending orai tokens transaction success

# add reporter to oraichain
echo "y" | oraicli tx provider add-reporters $(websocket keys list -a) --from $1 --fees 5000orai --keyring-backend test

sleep 2

# run websocket
websocket run