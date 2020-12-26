# Initialize configuration files and genesis file
# moniker is the name of your node
#nsd init <moniker> --chain-id namechain
USER=${1:-duc}
MONIKER=${2:-"$USER"_"Oraichain"_$(($RANDOM%10000000000))}
MIN_SELF_DELEGATION=${3:-100}

rm -rf .oraid/

rm -rf .oraicli/

rm -rf .oraifiles/

# rm -rf .images/

# rm -rf .websocket/

# rm go.sum

# go get ./...

# make install

oraid init $MONIKER --chain-id Oraichain

# Configure your CLI to eliminate need to declare them as flags
oraicli config chain-id Oraichain
oraicli config output json
oraicli config indent true
oraicli config trust-node true

# Copy the `Address` output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
# Note: In order for a new full node to join the network, after creating a local 
oraicli keys add $USER

# Copy the `Address` output here and save it for later use
# oraicli keys add hongeinh

# Add both accounts, with coins to the genesis file
oraid add-genesis-account $(oraicli keys show $USER -a) 9000000000000000orai
# oraid add-genesis-account $(oraicli keys show hongeinh -a) 1000000000000orai

# oraicli tx staking create-validator --amount 10000orai --pubkey oraivalconspub1addwnpepqvydmv22mkzc9rc92g43unew08cmj4q46dhk7vz0a9fj2xjsjn2lvqj0dfr --moniker ducphamle --chain-id Oraichain --commission-rate 0.10 --commission-max-rate 0.20 --commission-max-change-rate 0.01 --min-self-delegation 100 --gas auto --gas-adjustment 1.15 --gas-prices 0.025orai --from duc

# Note that for creating a validator, the gas may not be enough so it should be set to auto with a gas-adjustment value

# The "nscli config" command saves configuration for the "nscli" command but not for "nsd" so we have to 
# declare the keyring-backend with a flag here
#nsd gentx --name jack <or your key_name> --keyring-backend test
oraid gentx --amount 900000000000orai --name $USER --min-self-delegation $MIN_SELF_DELEGATION

# put the validators into the genesis file so the chain is aware of the validators
oraid collect-gentxs

oraid validate-genesis

# oraid start --minimum-gas-prices 0.025orai

#oraid start

############################################# CLI requests

# oraicli tx provider set-oscript test ./testfiles/oscript_kyc.sh "" --from duc --fees 5000orai

# oraicli tx provider set-datasource test ./testfiles/datasource_kyc.sh "" --from duc --fees 5000orai

# oraicli tx provider set-testcase test ./testfiles/testcase.sh "" --from duc --fees 5000orai

# oraicli query provider onames

## oraicli tx provider set-datasource test1 'curl -s -X GET "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd" -H "accept: application/json" | jq -r ".[\"bitcoin\"].usd"' --from duc

## oraicli tx provider set-aireq test --from duc

## oraicli query provider aireqs

## oraicli query provider aireq 1gRw15C6LF2I85NStROMA1Ggmkh (remember to remove prefix key of the request or name of oscript / data source to query)

###################################### REST requests

# oraicli rest-server --chain-id Oraichain --trust-node

# curl -XPOST -s http://localhost:1317/provider/oscript --data-binary '{"base_req":{"from":"'$(oraicli keys show duc -a)'","chain_id":"Oraichain"},"name":"testing","code":"./testfiles/oscript.sh"}' > unsignedTx.json

# curl -s -X POST -H "Content-Type: multipart/form-data" -F "image=@images/sample.png" -F "oracle_script_name=oscript_classification" -F "fees=45000orai" -F "from=$(oraicli keys show duc -a)" -F "chain_id=Oraichain" -F "input=''" -F "expected_output=5000" -F "validator_count=1" "http://localhost:1317/airequest/aireq/kycreq" > unsignedTx.json

# curl -s "http://localhost:1317/auth/accounts/$(oraicli keys show duc -a)"

# oraicli tx sign unsignedTx.json --from duc --offline --chain-id Oraichain --sequence 4 --account-number 0 > signedTx.json

# oraicli tx broadcast signedTx.json

# curl -s "http://localhost:1317/provider/oscript/test"

# curl -s "http://localhost:1317/provider/oscripts?page=1&limit=5"

# curl -s -X POST -H "Content-Type: multipart/form-data" -F "image=@images/Screenshot from 2020-08-25 16-04-46.png" -F "oracle_script_name=test" -F "fees=35000orai" -F "from=$(oraicli keys show duc -a)" -F "chain_id=Oraichain" "http://localhost:1317/provider/aireq/kycreq" > unsignedTx.json

# curl -s -X POST -H "Content-Type: multipart/form-data" -F "oracle_script_name=oscript_price" -F "fees=35000orai" -F "from=$(oraicli keys show duc -a)" -F "chain_id=Oraichain" -F "input=''" -F "expected_output=NTAwMA==" "http://localhost:1317/airequest/aireq/kycreq" > unsignedTx.json

# curl -XPOST -s http://localhost:1317/airequest/aireq/kycreq --data-binary '{"base_req":{"from":"'$(oraicli keys show duc -a)'","chain_id":"Oraichain"},"oracle_script_name":"oscript_price","input":"","expected_output":"NTAwMA==","fees":"35000orai"}' > unsignedTx.json

# curl -X POST -F "image=@images/Screenshot from 2020-08-25 16-04-46.png" "http://164.90.180.95:5001/api/v0/add"

# curl -o image.png -X POST "http://164.90.180.95:5001/api/v0/cat?arg=QmXipiMWjkm9ggFHbAVsLyixSGKHRRKhSTbtYWqDVDwZrf"

# curl -XPOST -s http://165.232.118.44:1317/provider/aireq/pricereq --data-binary '{"base_req":{"from":"'$(oraicli keys show duc -a)'","chain_id":"Oraichain"},"oracle_script_name":"oscript_price","input":"","expected_output":"NTAwMA==","fees":"35000orai"}' > unsignedTx.json

# curl -s "http://165.232.118.44:1317/auth/accounts/$(oraicli keys show duc -a)"

# oraicli tx sign unsignedTx.json --from duc --offline --chain-id Oraichain --sequence 14 --account-number 3 > signedTx.json

# oraicli tx send orai15tmx7qke90puyx2y63jkr8fl4fcjw06hm7alam orai10merxsp7n7wn5k0h5ux5e3acsz8rpg80679zky 10000000000orai --from orai15tmx7qke90puyx2y63jkr8fl4fcjw06hm7alam --fees 5000orai

# oraicli tx provider edit-datasource datasource coingecko_btc ./testfiles/coingecko_btc.sh "A data source that fetches the BTC price from Coingecko API" --from duc --fees 5000orai



# oraicli tx provider set-datasource coindesk_btc ./testfiles/coindesk_btc.sh "A data source that collects BTC price from coindesk" --from duc --fees 5000orai

## preparation
# oraicli tx provider set-datasource coingecko_eth ./testfiles/coingecko_eth.sh "A data source that fetches the ETH price from Coingecko API" --from duc --fees 5000orai

# oraicli tx provider set-datasource crypto_compare_eth ./testfiles/crypto_compare_eth.sh "A data source that collects ETH price from crypto compare" --from duc --fees 5000orai

# oraicli tx provider edit-testcase testcase_price testcase_price ./testfiles/testcase_price.sh "A sample test case that uses the expected output of users provided to verify the bitcoin price from the datasource" --from duc --fees 5000orai

# oraicli tx provider set-oscript oscript_eth ./testfiles/oscript_eth.sh "An oracle script that fetches and aggregates ETH price from different sources" --from duc --fees 5000orai

# curl -XPOST -s http://localhost:8000/api/txs/req_price -H "Content-Type: application/json" --data '{"oscript_name": "oscript_btc","price": "MA==","expected_price": "MA==","fees": "35000"}'

