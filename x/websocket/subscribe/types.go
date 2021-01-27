package subscribe

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/tendermint/tendermint/libs/log"
)

// WebSocketConfig is the extra config required for wasm
type WebSocketConfig struct {
	FromValidator    string
	BroadcastTimeout time.Duration
	RPCPollInterval  time.Duration
	MaxTry           uint64
	Txf              tx.Factory
	AllowLogLevel    log.Option
}
