#!/bin/bash
res=$(curl --location --request POST 'https://image-classification.v-chain.vn/v1/short-classification' --form 'image=@images/photo_2020-11-12_11-38-52.jpg' --form 'model=resnet18')

code=$(jq '.code' <<< "$res")
data=$(jq '.data' <<< "$res")

# check status code of the request
if [[ $code = 200 ]]
  then 
    echo "$data" | tr -d '"' # return result of the image classification and remove double quotes
else
    echo null
fi