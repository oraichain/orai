#!/bin/sh
# download the docker-compose & orai.env file

curl -OL https://raw.githubusercontent.com/oraichain/oraichain-static-files/master/docker-compose.genesis.yml

curl -OL https://raw.githubusercontent.com/oraichain/oraichain-static-files/master/orai.env

curl -OL https://raw.githubusercontent.com/oraichain/oraichain-static-files/master/init_genesis.sh

# modify the orai.env name & content