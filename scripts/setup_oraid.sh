#!/bin/sh
set -o errexit -o nounset -o pipefail

echo -n "Enter passphrase:"
read -s PASSWORD

CHAIN_ID=${CHAIN_ID:-Oraichain}
USER=${USER:-tupt}
MONIKER=${MONIKER:-node001}
# PASSWORD=${PASSWORD:-$1}
rm -rf .oraid/

oraid init --chain-id "$CHAIN_ID" "$MONIKER"

(echo "$PASSWORD"; echo "$PASSWORD") | oraid keys add $USER

# hardcode the validator account for this instance
(echo "$PASSWORD") | oraid add-genesis-account $USER "100000000000000orai"

# submit a genesis validator tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
(echo "$PASSWORD"; echo "$PASSWORD") | oraid gentx $USER "250000000orai" --chain-id="$CHAIN_ID" --amount="250000000orai" -y

oraid collect-gentxs


