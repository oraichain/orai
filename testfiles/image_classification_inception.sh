#!/bin/bash

#input=$(echo "$1" |base64 --decode)

res=$(curl --location --request POST "https://image-classification.v-chain.vn/v1/short-classification" --form "image=@$1" --form "model=inception_v3")

code=$(jq '.code' <<< "$res")
data=$(jq '.data' <<< "$res")

# check status code of the request
if [[ $code = 200 ]]
  then 
    echo "$data" | tr -d '"' # return result of the image classification and remove double quotes
else
    echo null
fi