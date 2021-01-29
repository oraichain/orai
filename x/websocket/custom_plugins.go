package websocket

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/oraichain/orai/x/websocket/types"
)

type OracleQueryPlugin struct {
	client  *http.Client
	bank    bankkeeper.ViewKeeper
	staking stakingkeeper.Keeper
	re      *regexp.Regexp
}

// Custom run custom command
func (oracleQueryPlugin OracleQueryPlugin) Custom(ctx sdk.Context, query json.RawMessage) ([]byte, error) {
	var request types.Request

	// also support proto
	err := ModuleCdc.UnmarshalJSON(query, &request)
	if err != nil {
		return nil, err
	}

	if request.Fetch.Method == "" {
		request.Fetch.Method = "GET"
	}

	r := strings.NewReader(request.Fetch.Body)
	req, err := http.NewRequest(request.Fetch.Method, request.Fetch.Url, r)

	// authorization header
	if request.Fetch.Authorization != "" {
		req.Header.Add("Authorization", request.Fetch.Authorization)
	}

	// call request
	resp, err := oracleQueryPlugin.client.Do(req)

	if err != nil {
		return ModuleCdc.LegacyAmino.MarshalJSON(map[string]string{"error": err.Error()})
	}

	defer resp.Body.Close()

	// we treat content as json, and fix float problem for wasm contract to help deterministically deserialization
	contents, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		contents = oracleQueryPlugin.re.ReplaceAll(contents, []byte(`$1"$2"$3`))
	}

	return contents, err
}

// CreateQueryPlugins create custom query
func CreateQueryPlugins(bank bankkeeper.ViewKeeper, staking stakingkeeper.Keeper) *wasm.QueryPlugins {

	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	oracleQueryPlugin := OracleQueryPlugin{
		client:  client,
		bank:    bank,
		staking: staking,
		re:      regexp.MustCompile(`([:,\[])\s*(\d+\.\d+)\s*([},\]])`),
	}

	return &wasm.QueryPlugins{
		Custom: oracleQueryPlugin.Custom,
	}
}
