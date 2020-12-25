package websocket

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/packages/filecache"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Context is used for the overall websocket information storage
type Context struct {
	client           rpcclient.Client
	validator        sdk.ValAddress
	gasPrices        sdk.DecCoins
	gas              uint64
	gasAdj           float64
	fees             sdk.Coins
	keys             chan keys.Info
	fileCache        filecache.Cache
	broadcastTimeout time.Duration
	maxTry           uint64
	rpcPollInterval  time.Duration
}
