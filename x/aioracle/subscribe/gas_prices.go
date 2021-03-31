package subscribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (subscriber *Subscriber) calculateGasPrices() (sdk.DecCoins, error) {
	// collect gas prices to create a new report
	minGasPricesObj := subscriber.config.Txf.GasPrices()
	minGasPrices := minGasPricesObj.String()
	// if the validator does not specify gas prices then we collect from the api get min gas prices
	if minGasPrices == "" {
		minGasPrices, err := getMinGasPrices()
		if err != nil {
			subscriber.log.Error(":exploding_head: Failed to collect the minimum gas prices of your node to create a new report with error: %s", err.Error())
			return nil, err
		}
		// test parsing the gas prices to prevent panic
		minGasPricesObj, err = sdk.ParseDecCoins(minGasPrices)
		if err != nil {
			subscriber.log.Error(":exploding_head: Invalid syntax for the minimum gas prices. Expected orai, got error: %s", err.Error())
			return nil, err
		}
	}
	return minGasPricesObj, nil
}
