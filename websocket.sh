#!/bin/bash

# ./websocket.sh duc 5
HOME=$PWD/.oraid
# rm -rf ~/.websocket
WEBSOCKET="websocket --home $HOME"

$WEBSOCKET keys delete-all 

# config chain id
$WEBSOCKET config chain-id Oraichain

# add validator to websocket config
$WEBSOCKET config validator $(oraicli keys show $1 -a --bech val --keyring-backend test)

# setup broadcast-timeout to websocket config
$WEBSOCKET config broadcast-timeout "30s"

# setup rpc-poll-interval to websocket config
$WEBSOCKET config rpc-poll-interval "1s"

# setup max-try to websocket config
$WEBSOCKET config max-try 5

# config log type
$WEBSOCKET config log-level debug

# echo "y" | oraicli tx oracle activate --from $1 --keyring-backend test

# wait for activation transaction success
sleep 2

# for i in $(eval echo {1..$2})
# do
  # add reporter key
  $WEBSOCKET keys add $2
# done

oraid start

sleep 2

# send orai tokens to reporters
echo "y" | oraicli tx send $(oraicli keys show $1 -a) $($WEBSOCKET keys show $2) 10000000orai --from $(oraicli keys show $1 -a) --fees 5000orai

sleep 6

#wait for sending orai tokens transaction success

# add reporter to oraichain
echo "y" | oraicli tx websocket add-reporters $($WEBSOCKET keys list -a) --from $1 --fees 5000orai --keyring-backend test

sleep 2

pkill oraid

# run websocket
# $WEBSOCKET run