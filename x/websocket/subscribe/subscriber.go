package subscribe

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	artypes "github.com/oraichain/orai/x/airequest/types"
	providerTypes "github.com/oraichain/orai/x/provider/types"
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
	log    log.Logger
	config *types.WebSocketConfig
}

// NewQuerier return querier implementation
func NewSubscriber(log log.Logger, config *types.WebSocketConfig) *Subscriber {
	return &Subscriber{
		log:    log,
		config: config,
	}
}

func getAttributeMap(attrs []sdk.Attribute) map[string][]string {
	ret := make(map[string][]string)
	for _, attr := range attrs {
		if values, ok := ret[attr.Key]; ok {
			values = append(values, attr.Value)
		} else {
			ret[attr.Key] = []string{attr.Value}
		}
	}
	return ret
}

// RegisterSubscribes register all sbuscribe
func RegisterSubscribes(cliCtx client.Context, subscriber *Subscriber) {
	go subscriber.subscribe(&cliCtx)
	fmt.Printf("Websocket Subscribing config: %v ...\n", subscriber.config)
}

// Implementation
func (subscriber *Subscriber) handleTransaction(cliCtx *client.Context, queryClient types.QueryClient, tx *abci.TxResult) {
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

	for _, log := range logs {
		for _, ev := range log.Events {
			// process with each event type
			attrMap := getAttributeMap(ev.GetAttributes())
			switch ev.Type {
			case providerTypes.EventTypeSetDataSource:
				subscriber.handleDataSourceLog(cliCtx, queryClient, attrMap)
			case artypes.EventTypeSetAIRequest:
				subscriber.handleAIRequestLog(cliCtx, queryClient, attrMap)
			default:
				subscriber.log.Debug(":ghost: Skipping non-{request/packet} type: %s", ev.Type)
			}
		}
	}
}

func (subscriber *Subscriber) subscribe(cliCtx *client.Context) {

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

	queryClient := types.NewQueryClient(cliCtx)

	for {
		select {
		case ev := <-eventChan:
			txResult := ev.Data.(tmtypes.EventDataTx).TxResult
			subscriber.handleTransaction(cliCtx, queryClient, &txResult)
		}
	}
}

func (subscriber *Subscriber) newFactory(cliCtx *client.Context, memo string) tx.Factory {

	signModeStr := cliCtx.SignModeStr
	signMode := signing.SignMode_SIGN_MODE_UNSPECIFIED
	switch signModeStr {
	case flags.SignModeDirect:
		signMode = signing.SignMode_SIGN_MODE_DIRECT
	case flags.SignModeLegacyAminoJSON:
		signMode = signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	}
	return subscriber.config.Txf.WithTxConfig(cliCtx.TxConfig).
		WithAccountRetriever(cliCtx.AccountRetriever).
		WithKeybase(cliCtx.Keyring).
		WithChainID(cliCtx.ChainID).
		WithSignMode(signMode).
		WithMemo(memo)

}
