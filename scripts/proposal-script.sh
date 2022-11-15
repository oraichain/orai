#!/bin/sh

PASS=${PASS:-12345678}
VERSION=${VERSION:-"v0.44"}

(echo "$PASS") |oraid tx gov submit-proposal software-upgrade $VERSION --title "Upgrade Oraichain mainnet to v0.40.2 to pump Cosmos SDK to v0.42.11 & Tendermint to v0.34.14" --description "This upgrade will massively improve the network's performance because Tendermint has now implemented concurrent read & write to the database stored within each node"  --from $USER --upgrade-height $HEIGHT --upgrade-info "x" --deposit 10000000orai --chain-id $CHAIN_ID -y && (echo "$PASS") | oraid tx gov vote $ID yes --from $USER --chain-id $CHAIN_ID -y

