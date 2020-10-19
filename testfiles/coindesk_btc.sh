#!/bin/bash
curl -s -X GET "https://api.coindesk.com/v1/bpi/currentprice.json" -H "accept: application/json" | jq -r ".[\"bpi\"].USD.rate_float"