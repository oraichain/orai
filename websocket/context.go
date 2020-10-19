package websocket

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ducphamle2/dexai/packages/filecache"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Context is used for the overall websocket information storage
type Context struct {
	client    rpcclient.Client
	validator sdk.ValAddress
	gasPrices sdk.DecCoins
	keys      chan keys.Info
	//executor         executor.Executor
	fileCache        filecache.Cache
	broadcastTimeout time.Duration
	maxTry           uint64
	rpcPollInterval  time.Duration
}
