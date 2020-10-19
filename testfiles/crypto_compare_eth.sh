#!/bin/bash
curl -s -X GET "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD" -H "accept: application/json" | jq -r ".USD"