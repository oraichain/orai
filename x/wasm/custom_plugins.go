



package wasm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func RunCustomQuerier(_ sdk.Context, query json.RawMessage) ([]byte, error) {
	resp, _ := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")	

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body);

	n := map[string]string{"msg": string(contents)}

	return json.Marshal(n);
	// return nil, wasmvmtypes.UnsupportedRequest{Kind: "custom"}
}

func CreateQueryPlugins() QueryPlugins {
	return QueryPlugins{
		Custom:  RunCustomQuerier,
	}
}