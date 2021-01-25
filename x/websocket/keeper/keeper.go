package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/oraichain/orai/x/wasm"
	"github.com/oraichain/orai/x/websocket/types"
)

// Keeper of the provider store
type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             codec.Marshaler
	stakingKeeper   stakingkeeper.Keeper
	wasmKeeper      *wasm.Keeper
	websocketConfig types.WebSocketConfig
	//paramSpace       params.Subspace
}

// NewKeeper creates a provider keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, wasmKeeper *wasm.Keeper, s stakingkeeper.Keeper, config types.WebSocketConfig) *Keeper {
	// if !aiRequestSubspace.HasKeyTable() {
	// 	// register parameters of the provider module into the param space
	// 	aiRequestSubspace = aiRequestSubspace.WithKeyTable(types.ParamKeyTable())
	// }
	return &Keeper{
		storeKey:        key,
		cdc:             cdc,
		wasmKeeper:      wasmKeeper,
		stakingKeeper:   s,
		websocketConfig: config,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) GetConfig() *types.WebSocketConfig {
	return &k.websocketConfig
}

//IsNamePresent checks if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}
