package subscribe

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/types"
)

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func (subscriber *Subscriber) submitReport(msgReport *types.MsgCreateReport) (err error) {

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
	clientCtx := *subscriber.cliCtx
	// prepare factory to calculate gas
	txfTemp, err := tx.PrepareFactory(clientCtx, txf)
	if err != nil {
		return err
	}
	// calculate estimate gas used to get fees
	_, estimatedGas, err := tx.CalculateGas(clientCtx.QueryWithData, txfTemp, msgReport)
	if err != nil {
		subscriber.log.Error(":exploding_head: Cannot estimate gas with error %s", err.Error())
		return err
	}
	msgReport.Fees, _ = minGasPrices.MulDec(sdk.NewDec(int64(estimatedGas))).TruncateDecimal()
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
