package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/oraichain/orai/x/airequest"
	"github.com/oraichain/orai/x/airesult/types"
	"github.com/oraichain/orai/x/provider"
	"github.com/oraichain/orai/x/wasm"
	"github.com/oraichain/orai/x/websocket"
	"github.com/tendermint/tendermint/libs/log"
)

// always clone keeper to make it immutable
type (
	Keeper struct {
		cdc              codec.Marshaler
		storeKey         sdk.StoreKey
		wasmKeeper       *wasm.Keeper
		paramSpace       params.Subspace
		stakingKeeper    staking.Keeper
		providerKeeper   *provider.Keeper
		bankKeeper       bank.Keeper
		distrKeeper      distr.Keeper
		authKeeper       auth.AccountKeeper
		webSocketKeeper  websocket.Keeper
		aiRequestKeeper  airequest.Keeper
		feeCollectorName string
	}
)

// NewKeeper creates a airequest keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, wasmKeeper *wasm.Keeper, subspace params.Subspace, stakingKeeper staking.Keeper, providerKeeper *provider.Keeper, bankKeeper bank.Keeper, distrKeeper distr.Keeper, authKeeper auth.AccountKeeper, webSocketKeeper websocket.Keeper, aiRequestKeeper airequest.Keeper, feeCollectorName string) *Keeper {
	if !subspace.HasKeyTable() {
		// register parameters of the airequest module into the param space
		subspace = subspace.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		storeKey:         key,
		cdc:              cdc,
		wasmKeeper:       wasmKeeper,
		paramSpace:       subspace,
		stakingKeeper:    stakingKeeper,
		providerKeeper:   providerKeeper,
		bankKeeper:       bankKeeper,
		distrKeeper:      distrKeeper,
		authKeeper:       authKeeper,
		webSocketKeeper:  webSocketKeeper,
		aiRequestKeeper:  aiRequestKeeper,
		feeCollectorName: feeCollectorName,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
