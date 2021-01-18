package wasm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	resp, err := http.Get(url)

	if err != nil {
		return json.Marshal(map[string]string{"error": err.Error()})
	}

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)

	return json.Marshal(contents)

}

func CreateQueryPlugins() QueryPlugins {
	return QueryPlugins{
		Custom: RunCustomQuerier,
	}
}
