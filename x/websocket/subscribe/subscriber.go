package subscribe

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	artypes "github.com/oraichain/orai/x/airequest/types"
	providerTypes "github.com/oraichain/orai/x/provider/types"
	"github.com/oraichain/orai/x/websocket/keeper"
	"github.com/oraichain/orai/x/websocket/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// TxQuery ...
	TxQuery = "tm.event = 'Tx'"
	// EventChannelCapacity is a buffer size of channel between node and this program
	EventChannelCapacity = 5000
)

type Subscriber struct {
	keeper *keeper.Keeper
	log    log.Logger
	config *types.WebSocketConfig
}

// NewQuerier return querier implementation
func NewSubscriber(keeper *keeper.Keeper, log log.Logger, config *types.WebSocketConfig) *Subscriber {
	return &Subscriber{
		keeper: keeper,
		log:    log,
		config: config,
	}
}

// RegisterSubscribes register all sbuscribe
func RegisterSubscribes(cliCtx client.Context, subscriber *Subscriber) {
	go subscriber.registerWebSocketSubscribe(&cliCtx)
	fmt.Printf("Websocket Subscribing config: %v ...\n", subscriber.config)
}

func getAttributeMap(attrs []sdk.Attribute) map[string]string {
	ret := make(map[string]string)
	for _, attr := range attrs {
		ret[attr.Key] = attr.Value
	}
	return ret
}

// Implementation
func (subscriber *Subscriber) handleTransaction(cliCtx *client.Context, tx *abci.TxResult) {
	fmt.Printf(":eyes: Inspecting incoming transaction: %X\n", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		subscriber.log.Debug(":alien: Skipping transaction with non-zero code: %d", tx.Result.Code)
		return
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		subscriber.log.Error(":cold_sweat: Failed to parse transaction logs with error: %s", err.Error())
		return
	}

	queryClient := types.NewQueryClient(cliCtx)

	for _, log := range logs {
		for _, ev := range log.Events {
			// process with each event type
			attrMap := getAttributeMap(ev.GetAttributes())
			switch ev.Type {
			case providerTypes.EventTypeSetDataSource:
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
				fmt.Printf("contract address: %v, response: %s\n", attrMap[providerTypes.AttributeContractAddress], string(response.Data))

			case artypes.EventTypeSetAIRequest:
				subscriber.handleAIRequestLog(cliCtx, attrMap)
			default:
				subscriber.log.Debug(":ghost: Skipping non-{request/packet} type: %s", ev.Type)
			}
		}
	}
}

func (subscriber *Subscriber) registerWebSocketSubscribe(cliCtx *client.Context) {

	// Instantiate and start tendermint RPC client
	client, err := cliCtx.GetNode()
	if err != nil {
		panic(err)
	}

	// Initialize a new error group
	goCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventChan, err := client.Subscribe(goCtx, "", TxQuery, EventChannelCapacity)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case ev := <-eventChan:
			txResult := ev.Data.(tmtypes.EventDataTx).TxResult
			subscriber.handleTransaction(cliCtx, &txResult)
		}
	}
}
