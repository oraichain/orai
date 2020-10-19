#!/bin/bash
curl -s -X GET "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd" -H "accept: application/json" | jq -r ".[\"bitcoin\"].usd"