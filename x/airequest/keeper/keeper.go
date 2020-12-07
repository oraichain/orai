package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/oraichain/orai/x/airequest/types"
)

// Implements OracleScriptSet interface
var _ types.AIRequestSet = Keeper{}

// Keeper of the provider store
type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	paramSpace     params.Subspace
	stakingKeeper  types.StakingKeeper
	ProviderKeeper types.ProviderKeeper
}

// NewKeeper creates a provider keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, aiRequestSubspace params.Subspace, stakingKeeper types.StakingKeeper, providerKeeper types.ProviderKeeper) Keeper {
	if !aiRequestSubspace.HasKeyTable() {
		// register parameters of the provider module into the param space
		aiRequestSubspace = aiRequestSubspace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		storeKey:       key,
		cdc:            cdc,
		paramSpace:     aiRequestSubspace,
		stakingKeeper:  stakingKeeper,
		ProviderKeeper: providerKeeper,
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
