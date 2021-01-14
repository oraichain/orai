
package oracle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateQueryPlugins() QueryPlugins {
	return QueryPlugins {
		Custom:  CustomQuerier,
	}
}

func CustomQuerier(_ sdk.Context, query json.RawMessage) ([]byte, error) {
	fmt.Printf("xyuzzzzzzzzz: %s", query);
	resp, _ := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")	

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body);

	n := map[string]string{"msg": string(contents)}

	return json.Marshal(n);
	// return nil, wasmvmtypes.UnsupportedRequest{Kind: "custom"}
}