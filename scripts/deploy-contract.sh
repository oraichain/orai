#!/bin/sh

store_ret=oraid tx wasm store ../oraiwasm/package/plus/swapmap/artifacts/swapmap.wasm --from validator1 --chain-id testing -y --home ~/.oraid/validator1/ -y --keyring-backend test --fees 200orai --gas 20000000 -b block

code_id=$(echo $store_ret | jq -r '.logs[0].events[0].attributes[] | select(.key | contains("code_id")).value')

oraid tx wasm instantiate $code_id '{}' --label 'testing' --from validator1 --gas auto --gas-adjustment 1.2 --chain-id testing -y --keyring-backend test --home ~/.oraid/validator1/ -b block --admin $(oraid keys show validator1 --keyring-backend test --home ~/.oraid/validator1/ -a)  --fees 200orai