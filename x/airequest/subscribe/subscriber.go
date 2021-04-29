package subscribe

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/oraichain/orai/x/airequest/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"
)

const (
	// TxQuery ...
	TxQuery        = "tm.event = 'NewBlock'"
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

func GetAttributeMap(attrs []abci.EventAttribute) map[string][]string {
	ret := make(map[string][]string)
	for _, attr := range attrs {
		key := string(attr.Key)
		if _, ok := ret[key]; ok {
			ret[key] = append(ret[key], string(attr.Value))
		} else {
			ret[key] = []string{string(attr.Value)}
		}
	}
	return ret
}

// Implementation
func (subscriber *Subscriber) handleReport(queryClient types.QueryClient, ev abci.Event) error {
	subscriber.log.Info(":eyes: Inspecting incoming report: %X", ev)
	for _, attr := range ev.Attributes {
		if string(attr.Key) == types.AttributeReport {
			subscriber.submitReport(attr.Value)
		}
		if string(attr.Key) == types.AttributeTestCaseReport {
			subscriber.submitTestCaseReport(attr.Value)
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
				switch dataType := ev.Data.(type) {
				case tmtypes.EventDataNewBlock:
					events := dataType.ResultBeginBlock.Events
					for _, ev := range events {
						if ev.Type == types.EventTypeReportWithData {
							err = subscriber.handleReport(queryClient, ev)
							if err != nil {
								return nil
							}
							break
						}
					}
				}
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
