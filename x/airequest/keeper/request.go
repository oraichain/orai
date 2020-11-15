package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airequest/exported"
	"github.com/oraichain/orai/x/airequest/types"
)

// GetAIRequest returns the information of an AI request
func (k Keeper) GetAIRequest(ctx sdk.Context, id string) (exported.AIRequestI, error) {
	store := ctx.KVStore(k.storeKey)
	var result types.AIRequest
	err := k.cdc.UnmarshalBinaryBare(store.Get(types.RequestStoreKey(id)), &result)
	if err != nil {
		return types.AIRequest{}, err
	}
	return result, nil
}

// SetAIRequest allows users to set a oScript into the store
func (k Keeper) SetAIRequest(ctx sdk.Context, id string, request types.AIRequest) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryBare(request)
	if err != nil {
		fmt.Println("error: ", err)
	}
	store.Set(types.RequestStoreKey(id), bz)
}

// GetAllAIRequestIDs get an iterator of all key-value pairs in the store
func (k Keeper) GetAllAIRequestIDs(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte("req"))
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
		return sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0))))
	}
	for _, request := range requests {
		fees = fees.Add(request.Fees...)
	}
	return fees
}
