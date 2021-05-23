package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

func (subscriber *Subscriber) handleReport(queryClient types.QueryClient, msgReportBytes []byte) error {
	msgReport := &types.MsgCreateReport{BaseReport: &types.BaseReport{}}
	msgTcReport := &types.MsgCreateTestCaseReport{BaseReport: &types.BaseReport{}}
	errReport := json.Unmarshal(msgReportBytes, msgReport)
	errTcReport := json.Unmarshal(msgReportBytes, msgTcReport)
	if errReport != nil {
		if errTcReport != nil {
			subscriber.log.Error(":exploding_head: Failed to unmarshal test case report and report with err: %v, %v", errReport, errTcReport)
			return fmt.Errorf("Cannot unmarshal report")
		}
		subscriber.log.Info(":eyes: unmarshal test case report successfully: %v", msgTcReport)
		subscriber.submitTestCaseReport(queryClient, msgTcReport)
		return nil
	}
	subscriber.submitReport(queryClient, msgReport)
	return nil
}

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func (subscriber *Subscriber) submitReport(queryClient types.QueryClient, msgReport *types.MsgCreateReport) (err error) {
	fmt.Println("message report: ", msgReport.BaseReport)

	if err = msgReport.ValidateBasic(); err != nil {
		subscriber.log.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return err
	}
	minGasPrices, err := subscriber.calculateGasPrices(queryClient)
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
func (subscriber *Subscriber) submitTestCaseReport(queryClient types.QueryClient, msgReport *types.MsgCreateTestCaseReport) (err error) {
	if err = msgReport.ValidateBasic(); err != nil {
		subscriber.log.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return err
	}
	minGasPrices, err := subscriber.calculateGasPrices(queryClient)
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

func (subscriber *Subscriber) calculateGasPrices(queryClient types.QueryClient) (sdk.DecCoins, error) {
	// collect gas prices to create a new report
	minGasPricesObj := subscriber.config.Txf.GasPrices()
	minGasPrices := minGasPricesObj.String()
	// if the validator does not specify gas prices then we collect from the api get min gas prices
	if minGasPrices == "" {
		minGasPrices, err := queryClient.QueryMinGasPrices(context.Background(), &types.MinGasPricesReq{})
		if err != nil {
			subscriber.log.Error(":exploding_head: Failed to collect the minimum gas prices of your node to create a new report with error: %s", err.Error())
			return nil, err
		}
		// test parsing the gas prices to prevent panic
		minGasPricesObj, err = sdk.ParseDecCoins(minGasPrices.MinGasPrices)
		if err != nil {
			subscriber.log.Error(":exploding_head: Invalid syntax for the minimum gas prices. Expected orai, got error: %s", err.Error())
			return nil, err
		}
	}
	return minGasPricesObj, nil
}
