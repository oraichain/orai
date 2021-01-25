package subscribe

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	utils "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/oraichain/orai/x/websocket/types"
)

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func SubmitReport(cliCtx client.Context, config *types.WebSocketConfig, msgReport *types.MsgCreateReport) {

	txHash := ""
	if err := msgReport.ValidateBasic(); err != nil {
		fmt.Println(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return
	}

	for try := uint64(1); try <= config.MaxTry; try++ {
		fmt.Println(":e-mail: Try to broadcast report transaction(%d/%d)", try, config.MaxTry)

		txBldr := cliCtx.TxConfig.NewTxBuilder()
		txBldr.SetMemo("websocket")

		bytes, err := cliCtx.JSONMarshaler.MarshalJSON(msgReport)

		res, err := cliCtx.BroadcastTxSync(bytes)
		fmt.Println(":star: response: %v", res)
		if err == nil {
			txHash = res.TxHash
			break
		}
		fmt.Println(":warning: Failed to broadcast tx with error: %s", err.Error())
		time.Sleep(config.RPCPollInterval)
	}
	if txHash == "" {
		fmt.Println(":exploding_head: Cannot try to broadcast more than %d try", config.MaxTry)
		return
	}
	for start := time.Now(); time.Since(start) < config.BroadcastTimeout; {
		time.Sleep(config.RPCPollInterval)
		txRes, err := utils.QueryTx(cliCtx, txHash)
		if err != nil {
			fmt.Println(":warning: Failed to query tx with error: %s", err.Error())
			continue
		}
		if txRes.Code != 0 {
			fmt.Println(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", txRes.Code, txRes.RawLog, txRes.TxHash)
			return
		}
		fmt.Println(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", txHash)
		return
	}
	fmt.Println(":question_mark: Cannot get transaction response from hash: %s transaction might be included in the next few blocks or check your node's health.", txHash)
}
