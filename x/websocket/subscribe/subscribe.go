package subscribe

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/oraichain/orai/x/websocket/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// TxQuery ...
	TxQuery = "tm.event = 'Tx'"
	// EventChannelCapacity is a buffer size of channel between node and this program
	EventChannelCapacity = 5000
)

// RegisterSubscribes register all sbuscribe
func RegisterSubscribes(cliCtx client.Context, log log.Logger, config *types.WebSocketConfig) {
	go registerWebSocketSubscribe(cliCtx, log, config)
	fmt.Printf("Websocket Subscribing config: %v ...\n", config)
}

func registerWebSocketSubscribe(cliCtx client.Context, log log.Logger, config *types.WebSocketConfig) {

	// Instantiate and start tendermint RPC client
	client, err := cliCtx.GetNode()
	if err != nil {
		panic(err)
	}

	// Initialize a new error group
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventChan, err := client.Subscribe(ctx, "", TxQuery, EventChannelCapacity)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case ev := <-eventChan:
			txResult := ev.Data.(tmtypes.EventDataTx).TxResult
			handleTransaction(ctx, log, config, &txResult)
		}
	}
}
