package subscribe

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	artypes "github.com/oraichain/orai/x/airequest/types"
	"github.com/oraichain/orai/x/websocket/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"
)

const (
	// TxQuery ...
	TxQuery        = "tm.event = 'Tx'"
	SubscriberName = "txevents"
	// EventChannelCapacity is a buffer size of channel between node and this program
	EventChannelCapacity = 5000
)

type Subscriber struct {
	cliCtx *client.Context
	config *WebSocketConfig
	log    *Logger
}

// NewQuerier return querier implementation
func NewSubscriber(cliCtx *client.Context, config *WebSocketConfig) *Subscriber {
	log := NewLogger(config.AllowLogLevel)
	return &Subscriber{
		cliCtx: cliCtx,
		config: config,
		log:    log,
	}
}

func GetAttributeMap(attrs []sdk.Attribute) map[string][]string {
	ret := make(map[string][]string)
	for _, attr := range attrs {
		if _, ok := ret[attr.Key]; ok {
			ret[attr.Key] = append(ret[attr.Key], attr.Value)
		} else {
			ret[attr.Key] = []string{attr.Value}
		}
	}
	return ret
}

// Implementation
func (subscriber *Subscriber) handleTransaction(queryClient types.QueryClient, tx *abci.TxResult) error {
	subscriber.log.Info(":eyes: Inspecting incoming transaction: %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		subscriber.log.Debug(":alien: Skipping transaction with non-zero code: %d", tx.Result.Code)
		return nil
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		subscriber.log.Error(":cold_sweat: Failed to parse transaction logs with error: %s", err.Error())
		return err
	}

	for _, log := range logs {
		for _, ev := range log.Events {
			// process with each event type
			switch ev.Type {
			case artypes.EventTypeRequestWithData:
				msgReport, err := subscriber.handleAIRequestLog(queryClient, &ev)
				if err != nil {
					return err
				}
				// split two functions to do unit test easier an decouple the business logic
				err = subscriber.submitReport(msgReport)
			default:
				subscriber.log.Debug(":ghost: Skipping non-{request/packet} type: %s", ev.Type)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Subscribe subscribe to event log
func (subscriber *Subscriber) Subscribe() error {

	subscriber.log.Info(":beer: Websocket Subscribing with validator: %s, exit on error: %b ...",
		subscriber.cliCtx.GetFromName(), subscriber.config.ErrExit)

	// Instantiate and start tendermint RPC client
	client, err := subscriber.cliCtx.GetNode()
	if err != nil {
		return err
	}

	if err = client.Start(); err != nil {
		return err
	}

	// Initialize a new error group
	goCtx, cancel := context.WithCancel(context.Background())
	eventChan, err := client.Subscribe(goCtx, SubscriberName, TxQuery, EventChannelCapacity)
	if err != nil {
		cancel()
		return err
	}

	defer func() {
		cancel()
		err = client.UnsubscribeAll(goCtx, SubscriberName)
	}()

	queryClient := types.NewQueryClient(subscriber.cliCtx)

	g, ctx := errgroup.WithContext(goCtx)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				break
			case ev := <-eventChan:
				txResult := ev.Data.(tmtypes.EventDataTx).TxResult
				err = subscriber.handleTransaction(queryClient, &txResult)
			}

			// check exit on error
			if err != nil && subscriber.config.ErrExit {
				return err
			}
		}
	})

	// Exit on error or context cancelation
	return g.Wait()
}

func (subscriber *Subscriber) newTxFactory(memo string) tx.Factory {
	// set sequence = 0 to retrieve later
	return subscriber.config.Txf.
		WithSequence(0).
		WithMemo(memo)
}
