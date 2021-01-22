package rest

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// TxQuery ...
	TxQuery = "tm.event = 'Tx'"
	// EventChannelCapacity is a buffer size of channel between node and this program
	EventChannelCapacity = 5000
)

func registerSubscribe(cliCtx client.Context) {

	// Instantiate and start tendermint RPC client
	client, err := cliCtx.GetNode()
	if err != nil {
		panic(err)
	}

	if err = client.Start(); err != nil {
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
			fmt.Printf("%v\n", ev.Data.(tmtypes.EventDataTx).TxResult)
			// go handleTransaction(c, l, ev.Data.(tmtypes.EventDataTx).TxResult)
		}
	}
}
