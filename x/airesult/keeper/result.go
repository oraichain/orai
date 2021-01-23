package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	airequest "github.com/oraichain/orai/x/airequest/types"
	"github.com/oraichain/orai/x/airesult/types"
	websocket "github.com/oraichain/orai/x/websocket/types"
)

// ResolveResult aggregates the results from the reports before storing it into the blockchain
func (k Keeper) ResolveResult(ctx sdk.Context, req *airequest.AIRequest, rep *websocket.Report) {
	// hard code the result first if the request does not have a result

	// get a new validator object for the result object.
	valAddress := rep.GetReporter().GetValidator()
	validator := k.webSocketKeeper.NewValidator(valAddress, k.stakingKeeper.Validator(ctx, valAddress).GetConsensusPower(), "active")

	if !k.HasResult(ctx, req.GetRequestID()) {
		// if the the request only needs a validator to return a result from the report then it's finished
		resultList := make([]websocket.ValResult, 0)
		resultList = append(resultList, *k.webSocketKeeper.NewValResult(validator, rep.GetAggregatedResult(), rep.GetResultStatus()))
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
			result.Results = append(result.Results, *k.webSocketKeeper.NewValResult(validator, rep.GetAggregatedResult(), rep.GetResultStatus()))

			// validate param
			totalReportsParam := k.GetTotalReportsParam(ctx)
			if k.GetTotalReportsParam(ctx) <= int64(0) {
				totalReportsParam = int64(70)
			}
			// check if there are enough results from the validators or not
			ratio := sdk.NewDecWithPrec(totalReportsParam, 2)

			// the number of reports that the user requires
			reportLengths := sdk.NewDec(int64(len(req.GetValidators())))

			// the threshold that the length of the result must pass
			threshold := reportLengths.Mul(ratio)

			// the actual result length
			resultLengths := sdk.NewDec(int64(len(result.Results)))

			// if the result length is GTE the threshold then the result is valid, and considered finished
			if resultLengths.GTE(threshold) {
				result.Status = types.RequestStatusFinished
			}
			// store the result
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
func (k Keeper) GetResult(ctx sdk.Context, reqID string) (*types.AIRequestResult, error) {
	store := ctx.KVStore(k.storeKey)
	var result types.AIRequestResult
	err := k.cdc.UnmarshalBinaryBare(store.Get(types.ResultStoreKey(reqID)), &result)
	if err != nil {
		return &types.AIRequestResult{}, err
	}
	return &result, nil
}

// SetResult allows users to set a result into the store
func (k Keeper) SetResult(ctx sdk.Context, reqID string, result *types.AIRequestResult) error {
	store := ctx.KVStore(k.storeKey)

	bz, err := k.cdc.MarshalBinaryBare(result)
	store.Set(types.ResultStoreKey(reqID), bz)
	return err
}
