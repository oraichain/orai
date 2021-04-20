package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// GetAIOracle returns the information of an AI request
func (k *Keeper) GetAIOracle(ctx sdk.Context, id string) (*types.AIOracle, error) {
	store := ctx.KVStore(k.StoreKey)
	hasAIOracle := store.Has(types.RequestStoreKey(id))
	var err error
	if !hasAIOracle {
		err = fmt.Errorf("")
		return nil, err
	}
	result := &types.AIOracle{}
	err = k.Cdc.UnmarshalBinaryBare(store.Get(types.RequestStoreKey(id)), result)
	return result, err
}

// HasAIOracle checks if there exists an ai request given an id
func (k *Keeper) HasAIOracle(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.StoreKey)
	return store.Has(types.RequestStoreKey(id))
}

// SetAIOracle allows users to set a oScript into the store
func (k *Keeper) SetAIOracle(ctx sdk.Context, id string, request *types.AIOracle) {
	store := ctx.KVStore(k.StoreKey)
	bz, err := k.Cdc.MarshalBinaryBare(request)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error: %v\n", err.Error()))
	}
	store.Set(types.RequestStoreKey(id), bz)
}

// GetAIOracleIDIter get an iterator of all key-value pairs in the store
func (k *Keeper) GetAIOracleIDIter(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIterator(store, types.RequeststoreKeyPrefixAll())
}

// GetPaginatedAIOracles get an iterator of paginated key-value pairs in the store
func (k *Keeper) GetPaginatedAIOracles(ctx sdk.Context, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIteratorPaginated(store, types.RequeststoreKeyPrefixAll(), page, limit)
}

// GetAIOraclesBlockHeight returns all ai oracle requests for the given block height, or nil if there is none for validators to execute.
func (k *Keeper) GetAIOraclesBlockHeight(ctx sdk.Context) (aiOracles []types.AIOracle) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.RequeststoreKeyPrefixAll())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var req types.AIOracle
		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &req)
		// check if block height is equal or not
		if req.GetBlockHeight() == ctx.BlockHeight()-1 {
			aiOracles = append(aiOracles, req)
		}
	}
	return aiOracles
}

// // ResolveExpiredRequest handles requests that have been expired
// func (k *Keeper) ResolveExpiredRequest(ctx sdk.Context, reqID string) {
// 	// hard code the result first if the request does not have a result
// 	if !k.HasResult(ctx, reqID) {
// 		valResults := make([]types.ValResult, 0)
// 		k.SetResult(ctx, reqID, types.NewAIOracleResult(reqID, valResults, types.RequestStatusExpired))
// 	} else {
// 		// if already has result then we change the request status to expired
// 		result, _ := k.GetResult(ctx, reqID)
// 		result.Status = types.RequestStatusExpired
// 		k.SetResult(ctx, reqID, result)
// 	}
// }

// ResolveRequestsFromReports handles the reports received in a block to group all the validators, data source owners and test case owners
func (k *Keeper) ResolveRequestsFromReports(ctx sdk.Context, rep *types.Report, reward *types.Reward) (bool, int) {

	req, _ := k.GetAIOracle(ctx, rep.BaseReport.GetRequestId())
	validation := k.validateReportBasic(ctx, req, rep, ctx.BlockHeight())
	// if the report cannot pass the validation basic then we skip the rest
	if !validation {
		return false, 0
	}

	// this temp var is used to calculate validator fees. Cannot use reward provider fees since it will be stacked by other functions we handle
	var providerFees sdk.Coins
	// collect data source owners that have their data sources executed to reward
	for _, dataSourceResult := range rep.GetDataSourceResults() {
		if dataSourceResult.GetStatus() == k.GetKeyResultSuccess() {
			reward.Results = append(reward.Results, dataSourceResult)
			reward.BaseReward.ProviderFees = reward.BaseReward.ProviderFees.Add(dataSourceResult.GetEntryPoint().GetProviderFees()...)
			providerFees = providerFees.Add(dataSourceResult.GetEntryPoint().GetProviderFees()...)
		}
	}
	// add validator fees into the total fees of all validators
	reward.BaseReward.ValidatorFees = reward.BaseReward.ValidatorFees.Add(k.CalculateValidatorFees(ctx, providerFees)...)
	// collect validator current status
	val := k.StakingKeeper.Validator(ctx, rep.BaseReport.GetValidatorAddress())
	// create a new validator wrapper and append to reward obj
	reward.BaseReward.Validators = append(reward.BaseReward.Validators, *k.NewValidator(rep.BaseReport.GetValidatorAddress(), val.GetConsensusPower(), val.GetStatus().String()))
	reward.BaseReward.TotalPower += val.GetConsensusPower()

	// return boolean and length of validator list to resolve result
	return true, len(req.GetValidators())
}

// ResolveRequestsFromTestCaseReports handles the test case reports received in a block to group all the validators, data source owners and test case owners
func (k *Keeper) ResolveRequestsFromTestCaseReports(ctx sdk.Context, rep *types.TestCaseReport, reward *types.Reward) {

	req, _ := k.GetAIOracle(ctx, rep.BaseReport.GetRequestId())
	validation := k.validateTestCaseReportBasic(ctx, req, rep, ctx.BlockHeight())
	// if the report cannot pass the validation basic then we skip the rest
	if !validation {
		return
	}

	// this temp var is used to calculate validator fees. Cannot use reward provider fees since it will be stacked by other functions we handle
	var providerFees sdk.Coins
	// collect data source owners that have their data sources executed to reward
	for _, result := range rep.GetResultsWithTestCase() {
		for _, tcResult := range result.GetTestCaseResults() {
			if tcResult.Status == types.ResultSuccess {
				reward.BaseReward.ProviderFees = reward.BaseReward.ProviderFees.Add(tcResult.GetEntryPoint().GetProviderFees()...)
				reward.Results = append(reward.Results)
				providerFees = providerFees.Add(tcResult.GetEntryPoint().GetProviderFees()...)
			}
		}
	}
	// add validator fees into the total fees of all validators
	reward.BaseReward.ValidatorFees = reward.BaseReward.ValidatorFees.Add(k.CalculateValidatorFees(ctx, providerFees)...)
	// collect validator current status
	val := k.StakingKeeper.Validator(ctx, rep.BaseReport.GetValidatorAddress())
	// create a new validator wrapper and append to reward obj
	reward.BaseReward.Validators = append(reward.BaseReward.Validators, *k.NewValidator(rep.BaseReport.GetValidatorAddress(), val.GetConsensusPower(), val.GetStatus().String()))
	reward.BaseReward.TotalPower += val.GetConsensusPower()
	return
}
