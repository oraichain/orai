package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airequest/types"
)

// GetAIRequest returns the information of an AI request
func (k Keeper) GetAIRequest(ctx sdk.Context, id string) (*types.AIRequest, error) {
	store := ctx.KVStore(k.storeKey)
	hasAIRequest := store.Has(types.RequestStoreKey(id))
	var err error
	if !hasAIRequest {
		err = fmt.Errorf("")
		return nil, err
	}
	result := &types.AIRequest{}
	err = k.cdc.UnmarshalBinaryBare(store.Get(types.RequestStoreKey(id)), result)
	return result, err
}

// HasAIRequest checks if there exists an ai request given an id
func (k Keeper) HasAIRequest(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.RequestStoreKey(id))
}

// SetAIRequest allows users to set a oScript into the store
func (k Keeper) SetAIRequest(ctx sdk.Context, id string, request *types.AIRequest) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryBare(request)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error: %v\n", err.Error()))
	}
	store.Set(types.RequestStoreKey(id), bz)
}

// GetAIRequestIDIter get an iterator of all key-value pairs in the store
func (k Keeper) GetAIRequestIDIter(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.RequeststoreKeyPrefixAll())
}

// GetPaginatedAIRequests get an iterator of paginated key-value pairs in the store
func (k *Keeper) GetPaginatedAIRequests(ctx sdk.Context, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIteratorPaginated(store, types.RequeststoreKeyPrefixAll(), page, limit)
}
