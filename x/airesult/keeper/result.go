package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	aiRequest "github.com/oraichain/orai/x/airequest/exported"
	"github.com/oraichain/orai/x/airesult/types"
	webSocket "github.com/oraichain/orai/x/websocket/exported"
	webSocketType "github.com/oraichain/orai/x/websocket/types"
)

// ResolveResult aggregates the results from the reports before storing it into the blockchain
func (k Keeper) ResolveResult(ctx sdk.Context, req aiRequest.AIRequestI, rep webSocket.ReportI) {
	// hard code the result first if the request does not have a result
	if !k.HasResult(ctx, req.GetRequestID()) {
		// if the the request only needs a validator to return a result from the report then it's finished
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAaa0")
		resultList := make([]webSocket.ValResultI, 0)
		resultList = append(resultList, webSocketType.NewValResult(rep.GetValidator(), rep.GetAggregatedResult()))
		fmt.Println("report: ", rep)
		fmt.Println("req validator length: ", len(req.GetValidators()))
		if len(req.GetValidators()) == 1 {
			k.SetResult(ctx, req.GetRequestID(), types.NewAIRequestResult(req.GetRequestID(), resultList, types.RequestStatusFinished))
		} else {
			// assume that we have already checked the number of validators required for a request
			k.SetResult(ctx, req.GetRequestID(), types.NewAIRequestResult(req.GetRequestID(), resultList, types.RequestStatusPending))
		}
	} else {
		// if there are more than one validators then we add more results
		if len(req.GetValidators()) > 1 {
			// if already has result then we add more results
			result, _ := k.GetResult(ctx, req.GetRequestID())
			result.Results = append(result.Results, webSocketType.NewValResult(rep.GetValidator(), rep.GetAggregatedResult()))
			// check if there are enough results from the validators or not
			if len(req.GetValidators()) == len(result.Results) {
				result.Status = types.RequestStatusFinished
			}
			k.SetResult(ctx, req.GetRequestID(), result)
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
