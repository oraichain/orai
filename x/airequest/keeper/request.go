package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/oraichain/orai/x/provider"
)

// GetAIRequest returns the information of an AI request
func (k Keeper) GetAIRequest(ctx sdk.Context, id string) (*types.AIRequest, error) {
	store := ctx.KVStore(k.storeKey)
	result := &types.AIRequest{}
	err := k.cdc.UnmarshalBinaryBare(store.Get(types.RequestStoreKey(id)), result)
	return result, err
}

// SetAIRequest allows users to set a oScript into the store
func (k Keeper) SetAIRequest(ctx sdk.Context, id string, request *types.AIRequest) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryBare(request)
	if err != nil {
		fmt.Println("error: ", err)
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

// GetRequestsBlockHeight returns all requests for the given block height, or nil if there is none.
func (k Keeper) GetRequestsBlockHeight(ctx sdk.Context, blockHeight int64) (reqs []types.AIRequest) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.RequeststoreKeyPrefixAll())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var req types.AIRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &req)
		// check if block height is equal or not
		if req.BlockHeight == blockHeight {
			reqs = append(reqs, req)
		}
	}
	return reqs
}

// CollectRequestFees collects total fees of the requests from the previous block to remove them from the fee collector
func (k Keeper) CollectRequestFees(ctx sdk.Context, blockHeight int64) (fees sdk.Coins) {
	// collect requests from the previous block
	requests := k.GetRequestsBlockHeight(ctx, blockHeight)
	if len(requests) == 0 {
		return sdk.NewCoins(sdk.NewCoin(provider.Denom, sdk.NewInt(int64(0))))
	}
	for _, request := range requests {
		fees = fees.Add(request.Fees...)
	}
	return fees
}
