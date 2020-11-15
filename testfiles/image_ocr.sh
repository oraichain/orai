#!/bin/bash

#input=$(echo "$1" |base64 --decode)

res=$(curl --location --request POST "https://ocr.v-chain.vn/v1/ocr" --form "image=@$1")

code=$(jq '.code' <<< "$res")
data=$(jq '.data' <<< "$res")

# check status code of the request
if [[ $code = 200 ]]
  then
    trim_quotes=$(echo "$data" | tr -d '"')
    # trim last \f character
    trim=$(echo $trim_quotes | rev | cut -c3- | rev)
    # repace all \n string with actual \n
    echo $trim | sed 's/\\n/ /g'
else
    echo null
fi