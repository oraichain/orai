package subscribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MinGasPrices is the struct for querying the minimum gas prices of a node
type MinGasPrices struct {
	Price string `json:"minFees"`
}

func getMinGasPrices() (string, error) {
	resp, err := http.Get("http://localhost:1317/provider/minfees?OracleScriptName=min_gas_prices&ValNum=0")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var minGasPrices MinGasPrices
	// collect response from the request
	json.Unmarshal(body, &minGasPrices)
	fmt.Println("min gas prices: ", minGasPrices.Price)
	return minGasPrices.Price, nil
}
