#!/bin/bash

input=$(echo "$1" |base64 --decode)

res=$(curl --location --request POST "https://ocr.v-chain.vn/v1/ocr" --form "image=@$input")

code=$(jq '.code' <<< "$res")
data=$(jq '.data' <<< "$res")

# check status code of the request
if [[ $code = 200 ]]
  then 
    echo "$data" | tr -d '"' # return result of the image classification and remove double quotes
else
    echo null
fi