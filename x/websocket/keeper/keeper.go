package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/oraichain/orai/x/airequest"
	"github.com/oraichain/orai/x/provider"
	"github.com/oraichain/orai/x/websocket/types"
)

// Keeper of the provider store
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             codec.Marshaler
	stakingKeeper   stakingkeeper.Keeper
	wasmKeeper      *wasm.Keeper
	providerKeeper  *provider.Keeper
	aiRequestKeeper *airequest.Keeper
	//paramSpace       params.Subspace
}

// NewKeeper creates a provider keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, wasmKeeper *wasm.Keeper, provider *provider.Keeper, s stakingkeeper.Keeper, r *airequest.Keeper) *Keeper {
	// if !aiRequestSubspace.HasKeyTable() {
	// 	// register parameters of the provider module into the param space
	// 	aiRequestSubspace = aiRequestSubspace.WithKeyTable(types.ParamKeyTable())
	// }
	return &Keeper{
		storeKey:        key,
		cdc:             cdc,
		wasmKeeper:      wasmKeeper,
		stakingKeeper:   s,
		providerKeeper:  provider,
		aiRequestKeeper: r,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
