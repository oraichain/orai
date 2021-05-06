package subscribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func (subscriber *Subscriber) submitReport(msgReportBytes []byte) (err error) {
	fmt.Println("msg report string: ", string(msgReportBytes))
	msgReport := &types.MsgCreateReport{BaseReport: &types.BaseReport{}}
	err = json.Unmarshal(msgReportBytes, &msgReport)
	if err != nil {
		subscriber.log.Error(":exploding_head: Failed to parse json with error: %s", err.Error())
		return err
	}
	fmt.Println("message report: ", msgReport.BaseReport)

	if err = msgReport.ValidateBasic(); err != nil {
		subscriber.log.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return err
	}
	minGasPrices, err := subscriber.calculateGasPrices()
	if err != nil {
		subscriber.log.Error(":exploding_head: Failed to collect gas prices: %s", err.Error())
		return err
	}
	txf := subscriber.newTxFactory("websocket")
	// add gas prices to pay for the report
	txf = txf.WithGasPrices(minGasPrices.String())
	txf = txf.WithGasAdjustment(subscriber.config.Txf.GasAdjustment())
	clientCtx := *subscriber.cliCtx
	//prepare factory to calculate gas
	txfTemp, err := tx.PrepareFactory(clientCtx, txf)
	if err != nil {
		subscriber.log.Error(":exploding_head: Failed to prepare factory: %s", err.Error())
		return err
	}
	// calculate estimate gas used to get fees
	_, estimatedGas, err := tx.CalculateGas(clientCtx.QueryWithData, txfTemp, msgReport)
	if err != nil {
		subscriber.log.Error(":exploding_head: Cannot estimate gas with error %s", err.Error())
		return err
	}
	msgReport.BaseReport.Fees, _ = minGasPrices.MulDec(sdk.NewDec(int64(estimatedGas))).TruncateDecimal()
	fmt.Println("client address: ", subscriber.cliCtx.GetFromName())
	// update report fees
	for try := uint64(1); try <= subscriber.config.MaxTry; try++ {
		subscriber.log.Info(":e-mail: Try to broadcast report transaction(%d/%d)", try, subscriber.config.MaxTry)
		err = tx.BroadcastTx(*subscriber.cliCtx, txf, msgReport)
		if err == nil {
			break
		}
		subscriber.log.Error(":warning: Failed to broadcast tx with error: %s", err.Error())
		time.Sleep(subscriber.config.RPCPollInterval)
	}

	// the last error
	return err
}

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func (subscriber *Subscriber) submitTestCaseReport(msgReportBytes []byte) (err error) {
	msgReport := &types.MsgCreateTestCaseReport{BaseReport: &types.BaseReport{}}
	err = json.Unmarshal(msgReportBytes, msgReport)
	if err != nil {
		return err
	}

	if err = msgReport.ValidateBasic(); err != nil {
		subscriber.log.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return err
	}
	minGasPrices, err := subscriber.calculateGasPrices()
	if err != nil {
		return err
	}
	txf := subscriber.newTxFactory("websocket")
	// add gas prices to pay for the report
	txf = txf.WithGasPrices(minGasPrices.String())
	txf = txf.WithGasAdjustment(subscriber.config.Txf.GasAdjustment())
	// prepare factory to calculate gas
	// calculate estimate gas used to get fees
	clientCtx := *subscriber.cliCtx
	//prepare factory to calculate gas
	txfTemp, err := tx.PrepareFactory(clientCtx, txf)
	if err != nil {
		subscriber.log.Error(":exploding_head: Failed to prepare factory: %s", err.Error())
		return err
	}
	// calculate estimate gas used to get fees
	_, estimatedGas, err := tx.CalculateGas(clientCtx.QueryWithData, txfTemp, msgReport)
	if err != nil {
		subscriber.log.Error(":exploding_head: Cannot estimate gas with error %s", err.Error())
		return err
	}
	msgReport.BaseReport.Fees, _ = minGasPrices.MulDec(sdk.NewDec(int64(estimatedGas))).TruncateDecimal()
	// update report fees
	for try := uint64(1); try <= subscriber.config.MaxTry; try++ {
		subscriber.log.Info(":e-mail: Try to broadcast report transaction(%d/%d)", try, subscriber.config.MaxTry)
		err = tx.BroadcastTx(*subscriber.cliCtx, txf, msgReport)
		if err == nil {
			break
		}
		subscriber.log.Error(":warning: Failed to broadcast tx with error: %s", err.Error())
		time.Sleep(subscriber.config.RPCPollInterval)
	}

	// the last error
	return err
}

// MinGasPrices is the struct for querying the minimum gas prices of a node
type MinGasPrices struct {
	Price string `json:"minFees"`
}

func getMinGasPrices() (string, error) {
	resp, err := http.Get("http://localhost:1317/aioracle/min-gas-prices")
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
