package websocket

import (
	"time"

	keys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	broadcastTimeout time.Duration
	maxTry           uint64
	rpcPollInterval  time.Duration
}
