



package wasm

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Request struct {
	Fetch struct {
		Url string `json:"url"`
	} `json:"fetch"`
}

func RunCustomQuerier(_ sdk.Context, query json.RawMessage) ([]byte, error) {
	var request Request
	json.Unmarshal(query, &request)

	url := request.Fetch.Url
	resp, _ := http.Get(url)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body);

	// n := map[string][]byte{"result": contents}
	// fmt.Printf("resxxx %s", n);
	return json.Marshal(contents);

	// return nil, wasmvmtypes.UnsupportedRequest{Kind: "custom"}
}

func CreateQueryPlugins() QueryPlugins {
	return QueryPlugins {
		Custom:  RunCustomQuerier,
	}
}