#/bin/bash

res=$(curl --location --request GET 'https://api.scan.orai.io/v1/txs-account/orai1pu9qqke50krgj970hmwmzj6s6ke4gzw4zge7ll?TxType=cosmos-sdk/MsgSend&page_id=2')

list=$(jq -r '.data[] | "\(.memo)"' <<< "$res")
printf "%s\n" "$list" > memo.txt