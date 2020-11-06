package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/ai-request/types"
)

// ResolveResult aggregates the results from the reports before storing it into the blockchain
func (k Keeper) ResolveResult(ctx sdk.Context, req types.AIRequest, rep types.Report) {
	// hard code the result first if the request does not have a result
	if !k.HasResult(ctx, req.RequestID) {
		// if the the request only needs a validator to return a result from the report then it's finished
		var resultList types.ValResults
		resultList = append(resultList, types.NewValResult(rep.Validator.Address, rep.AggregatedResult))
		if len(req.Validators) == 1 {
			k.SetResult(ctx, req.RequestID, types.NewAIRequestResult(req.RequestID, resultList, types.RequestStatusFinished))
		} else {
			// assume that we have already checked the number of validators required for a request
			k.SetResult(ctx, req.RequestID, types.NewAIRequestResult(req.RequestID, resultList, types.RequestStatusPending))
		}
	} else {
		// if there are more than one validators then we add more results
		if len(req.Validators) > 1 {
			// if already has result then we add more results
			result, _ := k.GetResult(ctx, req.RequestID)
			result.Results = append(result.Results, types.NewValResult(rep.Validator.Address, rep.AggregatedResult))
			// check if there are enough results from the validators or not
			if len(req.Validators) == len(result.Results) {
				result.Status = types.RequestStatusFinished
			}
			k.SetResult(ctx, req.RequestID, result)
		}
	}
}

// HasResult checks if a given request has result or not
func (k Keeper) HasResult(ctx sdk.Context, reqID string) bool {
	_, err := k.GetResult(ctx, reqID)
	if err != nil {
		return false
	}
	return true
}

// GetResult returns the result of a given request
func (k Keeper) GetResult(ctx sdk.Context, reqID string) (types.AIRequestResult, error) {
	store := ctx.KVStore(k.storeKey)
	var result types.AIRequestResult
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.ResultStoreKey(reqID)), &result)
	if err != nil {
		return types.AIRequestResult{}, err
	}
	return result, nil
}

// SetResult allows users to set a result into the store
func (k Keeper) SetResult(ctx sdk.Context, reqID string, result types.AIRequestResult) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(result)
	store.Set(types.ResultStoreKey(reqID), bz)
}
