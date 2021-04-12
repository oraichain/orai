package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
	aioracle "github.com/oraichain/orai/x/aioracle/types"
)

// ResolveResult aggregates the results from the reports before storing it into the blockchain
func (k Keeper) ResolveResult(ctx sdk.Context, rep *aioracle.Report, valCount int, totalReportPercentage uint64) {

	id := rep.GetRequestID()
	// get a new validator object for the result object.
	valAddress := rep.GetReporter().GetValidator()
	validator := k.NewValidator(valAddress, k.StakingKeeper.Validator(ctx, valAddress).GetConsensusPower(), "active")

	if !k.HasResult(ctx, id) {
		// if the the request only needs a validator to return a result from the report then it's finished
		resultList := make([]aioracle.ValResult, 0)
		resultList = append(resultList, *k.NewValResult(validator, rep.GetAggregatedResult(), rep.GetResultStatus()))
		if valCount == 1 {
			k.SetResult(ctx, id, types.NewAIOracleResult(id, resultList, types.RequestStatusFinished))
		} else {
			// assume that we have already checked the number of validators required for a request
			k.SetResult(ctx, id, types.NewAIOracleResult(id, resultList, types.RequestStatusPending))
		}
	} else {
		// if there are more than one validators then we add more results
		if valCount > 1 {
			// if already has result then we add more results
			result, _ := k.GetResult(ctx, id)
			result.Results = append(result.Results, *k.NewValResult(validator, rep.GetAggregatedResult(), rep.GetResultStatus()))

			// check if there are enough results from the validators or not
			ratio := sdk.NewDecWithPrec(int64(totalReportPercentage), 2)

			// the number of reports that the user requires
			reportLengths := sdk.NewDec(int64(valCount))

			// the threshold that the length of the result must pass
			threshold := reportLengths.Mul(ratio)

			// the actual result length
			resultLengths := sdk.NewDec(int64(len(result.Results)))

			// if the result length is GTE the threshold then the result is valid, and considered finished
			if resultLengths.GTE(threshold) {
				result.Status = types.RequestStatusFinished
			}
			// store the result
			k.SetResult(ctx, id, result)
		}
	}
}

// HasResult checks if a given request has result or not
func (k Keeper) HasResult(ctx sdk.Context, reqID string) bool {
	store := ctx.KVStore(k.StoreKey)
	return store.Has(types.ResultStoreKey(reqID))
}

// GetResult returns the result of a given request
func (k Keeper) GetResult(ctx sdk.Context, reqID string) (*types.AIOracleResult, error) {
	store := ctx.KVStore(k.StoreKey)
	var result types.AIOracleResult
	err := k.Cdc.UnmarshalBinaryBare(store.Get(types.ResultStoreKey(reqID)), &result)
	if err != nil {
		return &types.AIOracleResult{}, err
	}
	return &result, nil
}

// SetResult allows users to set a result into the store
func (k Keeper) SetResult(ctx sdk.Context, reqID string, result *types.AIOracleResult) error {
	store := ctx.KVStore(k.StoreKey)

	bz, err := k.Cdc.MarshalBinaryBare(result)
	store.Set(types.ResultStoreKey(reqID), bz)
	return err
}
