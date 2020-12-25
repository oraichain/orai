package websocket

import (
	"fmt"
	"time"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/oraichain/orai/app"
	"github.com/oraichain/orai/x/websocket/types"
)

var (
	cdc = app.MakeCodec()
)

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func SubmitReport(c *Context, l *Logger, key keys.Info, msgReport types.MsgCreateReport) {
	cliCtx := sdkCtx.CLIContext{Client: c.client, TrustNode: true, Codec: cdc}
	txHash := ""
	if err := msgReport.ValidateBasic(); err != nil {
		l.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return
	}

	for try := uint64(1); try <= c.maxTry; try++ {
		l.Info(":e-mail: Try to broadcast report transaction(%d/%d)", try, c.maxTry)
		acc, err := auth.NewAccountRetriever(cliCtx).GetAccount(key.GetAddress())
		if err != nil {
			l.Info(":warning: Failed to retreive account with error: %s", err.Error())
			time.Sleep(c.rpcPollInterval)
			continue
		}

		txBldr := auth.NewTxBuilder(
			auth.DefaultTxEncoder(cdc), acc.GetAccountNumber(), acc.GetSequence(),
			c.gas, c.gasAdj, false, cfg.ChainID, fmt.Sprintf("websocket"), sdk.NewCoins(), c.gasPrices,
		)
		// txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, []sdk.Msg{msg})
		// if err != nil {
		// 	l.Error(":exploding_head: Failed to enrich with gas with error: %s", err.Error())
		// 	return
		// }
		out, err := txBldr.WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, []sdk.Msg{msgReport})
		if err != nil {
			l.Info(":warning: Failed to build tx with error: %s", err.Error())
			time.Sleep(c.rpcPollInterval)
			continue
		}

		res, err := cliCtx.BroadcastTxSync(out)
		l.Info(":star: response: %v", res)
		if err == nil {
			txHash = res.TxHash
			break
		}
		l.Info(":warning: Failed to broadcast tx with error: %s", err.Error())
		time.Sleep(c.rpcPollInterval)
	}
	if txHash == "" {
		l.Error(":exploding_head: Cannot try to broadcast more than %d try", c.maxTry)
		return
	}
	for start := time.Now(); time.Since(start) < c.broadcastTimeout; {
		time.Sleep(c.rpcPollInterval)
		txRes, err := utils.QueryTx(cliCtx, txHash)
		if err != nil {
			l.Debug(":warning: Failed to query tx with error: %s", err.Error())
			continue
		}
		if txRes.Code != 0 {
			l.Error(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", txRes.Code, txRes.RawLog, txRes.TxHash)
			return
		}
		l.Info(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", txHash)
		return
	}
	l.Info(":question_mark: Cannot get transaction response from hash: %s transaction might be included in the next few blocks or check your node's health.", txHash)
}
