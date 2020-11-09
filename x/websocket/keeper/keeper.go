package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/types"
)

// Implements ReportSet interface
var _ types.ReportSet = Keeper{}

// Implements ReporterSet interface
var _ types.ReporterSet = Keeper{}

// Implements DataSourceResultSet interface
var _ types.DataSourceResultSet = Keeper{}

// Implements TestCaseResultSet interface
var _ types.TestCaseResultSet = Keeper{}

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements ValidatorSet interface
var _ types.ValResultSet = Keeper{}

// Keeper of the provider store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	stakingKeeper types.StakingKeeper
	//paramSpace       params.Subspace
}

// NewKeeper creates a provider keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, stakingKeeper types.StakingKeeper) Keeper {
	// if !aiRequestSubspace.HasKeyTable() {
	// 	// register parameters of the provider module into the param space
	// 	aiRequestSubspace = aiRequestSubspace.WithKeyTable(types.ParamKeyTable())
	// }
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		stakingKeeper: stakingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//IsNamePresent checks if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// // CreateStrategy allows users to create a new strategy into the store
// func (k Keeper) CreateStrategy(ctx sdk.Context, name string, strategy types.Strategy) {
// 	store := ctx.KVStore(k.storeKey)

// 	bz := k.cdc.MustMarshalBinaryLengthPrefixed(strategy)
// 	store.Set(types.StrategyStoreKey(strategy.StratID, name), bz)
// }

// // GetParam returns the parameter as specified by key as an uint64.
// func (k Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
// 	k.paramSpace.Get(ctx, key, &res)
// 	return res
// }

// // SetParam saves the given key-value parameter to the store.
// func (k Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
// 	k.paramSpace.Set(ctx, key, value)
// }

// // GetParams returns all current parameters as a types.Params instance.
// func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
// 	k.paramSpace.GetParamSet(ctx, &params)
// 	return params
// }
