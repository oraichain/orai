package websocket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
}

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
		oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("response error: %v\n", err))
		return []byte{}, fmt.Errorf("cannot connect to the given URL to retrieve the oracle response")
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("cannot read the oracle response content: %v\n", err))
		// return empty bytes to show that the response content has error
		return []byte{}, err
	}

	// // convert the content to struct to verify the status code
	// response := OracleResponsePlugin{}
	// err = json.Unmarshal(contents, &response)
	// if err != nil {
	// 	oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("the response content does not match the required structure: %v\n", err))
	// 	// return empty bytes to show that the response content has error
	// 	return []byte{}, err
	// }

	// // check response code and data length
	// codeInt, err := strconv.Atoi(response.ResponseCode)
	// if err != nil || (response.ResponseCode != "200" && codeInt != 200 && response.ResponseCode != "sucess") || len(response.Data) == 0 {

	// 	// bad case we dont care about the error, just return empty content, as long it is Binary compatible
	// 	oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("the status code is not correct, or the response does not contain any data\n"))
	// 	return []byte{}, fmt.Errorf("the status code is not correct, or the response does not contain any data")
	// }
	// remove double quotes from the response data
	// responseUnquote, _ := strconv.Unquote(string(response.Data))
	// return the actual content of the oracle response, like type:bytes
	responseBytes, err := ModuleCdc.Amino.MarshalJSON(contents)
	if err != nil {
		oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("cannot marshal the response data with error: %v\n", err))
		return []byte(fmt.Sprintf("cannot marshal the response data with error: %v\n", err)), err
	}
	return responseBytes, nil
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
