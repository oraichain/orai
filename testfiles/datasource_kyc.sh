#!/bin/bash
res=$(curl -s -X POST -H "Content-Type: multipart/form-data" -F "image=@$1" "http://104.248.99.206:8080/v1/identify")

status=$(jq '.status' <<< "$res")
user_id=$(jq '.data.user_id' <<< "$res")

echo $status
status_trim= echo "$status" | tr -d '"'
status_correct="success"
echo $status_correct
echo $status_trim

# check status code of the request
if [[ "$status_trim" = "$status_correct" ]]; then
    echo "$user_id" | tr -d '"' # return result of the image classification and remove double quotes
else
    echo null
fi