package subscribe

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	providerTypes "github.com/oraichain/orai/x/provider/types"
	"github.com/oraichain/orai/x/websocket/types"
)

func (subscriber *Subscriber) handleDataSourceLog(cliCtx *client.Context, queryClient types.QueryClient, attrMap map[string]string) {
	contractAddr, _ := sdk.AccAddressFromBech32(attrMap[providerTypes.AttributeContractAddress])
	query := &types.QueryContract{
		Contract: contractAddr,
		Request: &types.Request{
			Fetch: &types.Fetch{
				Url: "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd",
			},
		},
	}

	response, _ := queryClient.OracleInfo(
		context.Background(),
		query,
	)

	// only get json back, or can process in smart contract
	fmt.Printf("contract address: %s, response: %s\n", contractAddr.String(), string(response.Data))
}
