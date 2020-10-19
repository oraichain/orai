#!/bin/bash
curl -s -X GET "https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD" -H "accept: application/json" | jq -r ".USD"