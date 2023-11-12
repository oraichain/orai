#!/bin/bash

CHAIN_ID=${CHAIN_ID:-testing}
USER=${USER:-tupt}
MONIKER=${MONIKER:-node001}
# PASSWORD=${PASSWORD:-$1}
rm -rf .oraid/

oraid init --chain-id "$CHAIN_ID" "$MONIKER"

(cat .env) | oraid keys add $USER --recover --keyring-backend test

# hardcode the validator account for this instance
oraid add-genesis-account $USER "100000000000000orai"

# submit a genesis validator tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
oraid gentx $USER "250000000orai" --chain-id="$CHAIN_ID" --amount="250000000orai" -y

oraid collect-gentxs

oraid start --rpc.laddr tcp://0.0.0.0:26657 --grpc.address 0.0.0.0:9090

