#!/bin/bash
curl -s -X POST -H "Content-Type: multipart/form-data" -F "image=@$1" "http://104.248.99.206:8080/v1/identify" | jq -r ".[\"data\"].user_id"