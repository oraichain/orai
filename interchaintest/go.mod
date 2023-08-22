module github.com/oraichain/orai/tests/interchaintest

go 1.21.0

require (
	github.com/CosmWasm/wasmd v0.31.0
	github.com/cosmos/cosmos-sdk v0.45.16-ics
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ibc-go/v4 v4.3.1
	github.com/cosmos/interchain-accounts v0.2.6
	github.com/gorilla/mux v1.8.0
	github.com/osmosis-labs/osmosis/x/ibc-hooks v0.0.0-20230201151635-ef43e092d196
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.14.0
	github.com/rakyll/statik v0.1.7
	github.com/rs/zerolog v1.27.0
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.1
	github.com/strangelove-ventures/packet-forward-middleware/v4 v4.0.6
	github.com/stretchr/testify v1.8.2
	github.com/tendermint/tendermint v0.34.29
	github.com/tendermint/tm-db v0.6.8-0.20220506192307-f628bb5dc95b
)

replace (
	// Use cosmos keyring
	github.com/99designs/keyring => github.com/cosmos/keyring v1.2.0
	// fork wasmd so that we have legacy wasm types
	github.com/CosmWasm/wasmd => github.com/oraichain/wasmd v0.30.2-0.20230704072512-1f776e9a4dcf
	// go list -m -json github.com/oraichain/wasmvm@main | jq '.|"\(.Path) \(.Version)"' -r
	github.com/CosmWasm/wasmvm => github.com/oraichain/wasmvm v1.2.4
	// same version as cosmos-sdk
	github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.22.2
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

	// use Oraichain fork of cometbft
	// go list -m -json github.com/oraichain/cometbft@482cde0c4e04479d4164b1b4f7a0b90071c58b7f | jq '.|"\(.Path) \(.Version)"' -r
	github.com/tendermint/tendermint => github.com/oraichain/cometbft v0.34.30-0.20230711110635-482cde0c4e04
	google.golang.org/grpc => google.golang.org/grpc v1.33.2

// Fix upstream GHSA-h395-qcrw-5vmq vulnerability.
// TODO Remove it: https://github.com/cosmos/cosmos-sdk/issues/10409
)
