package websocket

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/oraichain/orai/x/wasm"
)

type Request struct {
	Fetch struct {
		Method        string `json:"method,omitempty"`
		Authorization string `json:"authorization,omitempty"`
		Body          string `json:"body,omitempty"`
		Url           string `json:"url"`
	} `json:"fetch"`
}

type OracleQueryPlugin struct {
	client  *http.Client
	bank    bankkeeper.ViewKeeper
	staking stakingkeeper.Keeper
}

func (oracleQueryPlugin OracleQueryPlugin) Custom(ctx sdk.Context, query json.RawMessage) ([]byte, error) {
	var request Request
	json.Unmarshal(query, &request)

	if request.Fetch.Method == "" {
		request.Fetch.Method = "GET"
	}

	// fmt.Printf("Request :%v\n", request.Fetch)

	r := strings.NewReader(request.Fetch.Body)
	req, err := http.NewRequest(request.Fetch.Method, request.Fetch.Url, r)

	// authorization header
	if request.Fetch.Authorization != "" {
		req.Header.Add("Authorization", request.Fetch.Authorization)
	}

	// call request
	resp, err := oracleQueryPlugin.client.Do(req)

	if err != nil {
		return json.Marshal(map[string]string{"error": err.Error()})
	}

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)

	return json.Marshal(contents)

}

func CreateQueryPlugins(bank bankkeeper.ViewKeeper, staking stakingkeeper.Keeper) *wasm.QueryPlugins {

	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	oracleQueryPlugin := OracleQueryPlugin{
		client:  client,
		bank:    bank,
		staking: staking,
	}

	return &wasm.QueryPlugins{
		Custom: oracleQueryPlugin.Custom,
	}
}
