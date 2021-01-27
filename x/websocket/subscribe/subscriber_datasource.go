package subscribe

import (
	"context"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	providerTypes "github.com/oraichain/orai/x/provider/types"
	"github.com/oraichain/orai/x/websocket/types"
)

func (subscriber *Subscriber) handleDataSourceLog(queryClient types.QueryClient, attrMap map[string][]string) {
	contractAddr, _ := sdk.AccAddressFromBech32(attrMap[providerTypes.AttributeContractAddress][0])
	query := &types.QueryOracleContract{
		Contract: contractAddr,
		Request: &types.Request{
			Fetch: &types.Fetch{
				Url: "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd",
			},
		},
	}

	ret, pub, err := subscriber.cliCtx.Keyring.Sign("duc", []byte("hello"))
	subscriber.log.Info("ret :%v %v %v", ret, pub, err)
	validator := subscriber.cliCtx.GetFromAddress().String()
	subscriber.log.Info("validator :%v", validator)

	response, _ := queryClient.OracleContract(
		context.Background(),
		query,
	)

	// only get json back, or can process in smart contract
	subscriber.log.Info("contract address: %s, response: %s", contractAddr.String(), string(response.Data))

	// test sign transaction to update smart contract
	contractMsg := &wasm.MsgExecuteContract{
		Sender:   validator,
		Contract: contractAddr.String(),
		Msg:      []byte(`{"increment":{}}`),
	}
	txf := subscriber.newTxFactory("websocket")
	tx.BroadcastTx(*subscriber.cliCtx, txf, contractMsg)
}
