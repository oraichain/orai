package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// GetAIOracle returns the information of an AI request
func (k Keeper) GetAIOracle(ctx sdk.Context, id string) (*types.AIOracle, error) {
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
func (k Keeper) HasAIOracle(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.StoreKey)
	return store.Has(types.RequestStoreKey(id))
}

// SetAIOracle allows users to set a oScript into the store
func (k Keeper) SetAIOracle(ctx sdk.Context, id string, request *types.AIOracle) {
	store := ctx.KVStore(k.StoreKey)
	bz, err := k.Cdc.MarshalBinaryBare(request)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error: %v\n", err.Error()))
	}
	store.Set(types.RequestStoreKey(id), bz)
}

// GetAIOracleIDIter get an iterator of all key-value pairs in the store
func (k Keeper) GetAIOracleIDIter(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIterator(store, types.RequeststoreKeyPrefixAll())
}

// GetPaginatedAIOracles get an iterator of paginated key-value pairs in the store
func (k *Keeper) GetPaginatedAIOracles(ctx sdk.Context, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIteratorPaginated(store, types.RequeststoreKeyPrefixAll(), page, limit)
}

// ResolveExpiredRequest handles requests that have been expired
func (k Keeper) ResolveExpiredRequest(ctx sdk.Context, reqID string) {
	// hard code the result first if the request does not have a result
	if !k.HasResult(ctx, reqID) {
		valResults := make([]types.ValResult, 0)
		k.SetResult(ctx, reqID, types.NewAIOracleResult(reqID, valResults, types.RequestStatusExpired))
	} else {
		// if already has result then we change the request status to expired
		result, _ := k.GetResult(ctx, reqID)
		result.Status = types.RequestStatusExpired
		k.SetResult(ctx, reqID, result)
	}
}

// ResolveRequestsFromReports handles the reports received in a block to group all the validators, data source owners and test case owners
func (k Keeper) ResolveRequestsFromReports(ctx sdk.Context, rep *types.Report, reward *types.Reward, rewardPercentage int64) (bool, int) {

	req, _ := k.GetAIOracle(ctx, rep.GetRequestID())
	validation := k.validateBasic(ctx, req, rep, ctx.BlockHeight())
	// if the report cannot pass the validation basic then we skip the rest
	if !validation {
		return false, 0
	}

	// // collect data source owners that have their data sources executed to reward
	// for _, dataSourceResult := range rep.GetDataSourceResults() {
	// 	if dataSourceResult.GetStatus() == k.webSocketKeeper.GetKeyResultSuccess() {
	// 		dataSource, _ := k.providerKeeper.GetAIDataSource(ctx, dataSourceResult.GetName())
	// 		reward.DataSources = append(reward.DataSources, *dataSource)
	// 		reward.ProviderFees = reward.ProviderFees.Add(dataSource.GetFees()...)
	// 	}
	// }

	// // collect data source owners that have their data sources executed to reward
	// for _, testCaseResult := range rep.GetTestCaseResults() {
	// 	testCase, _ := k.providerKeeper.GetTestCase(ctx, testCaseResult.GetName())
	// 	reward.TestCases = append(reward.TestCases, *testCase)
	// 	reward.ProviderFees = reward.ProviderFees.Add(testCase.GetFees()...)
	// }
	// change reward ratio to the ratio of validator
	rewardRatio := k.GetParam(ctx, types.KeyAIOracleRewardPercentages)

	// reward = 1 - oracle reward percentage × (data source fees + test case fees)
	valFees, _ := sdk.NewDecCoinsFromCoins(reward.ProviderFees...).MulDec(sdk.NewDec(int64(rewardRatio))).TruncateDecimal()
	// add validator fees into the total fees of all validators
	reward.ValidatorFees = reward.ValidatorFees.Add(valFees...)
	// store information into the reward struct to reward these entities in the next begin block
	valAddress := rep.GetReporter().GetValidator()

	// collect validator current status
	val := k.StakingKeeper.Validator(ctx, valAddress)
	// create a new validator wrapper and append to reward obj
	validator := k.NewValidator(valAddress, val.GetConsensusPower(), val.GetStatus().String())
	reward.Validators = append(reward.Validators, *validator)
	reward.TotalPower += validator.GetVotingPower()

	// return boolean and length of validator list to resolve result
	return true, len(req.GetValidators())
}

func (k Keeper) validateBasic(ctx sdk.Context, req *types.AIOracle, rep *types.Report, blockHeight int64) bool {
	// if the request has been expired
	// if req.GetBlockHeight()+int64(k.GetParam(ctx, types.KeyExpirationCount)) < blockHeight {
	// 	//TODO: NEED TO HANDLE THE EXPIRED REQUEST.
	// 	fmt.Println("Request has been expired")
	// 	k.ResolveExpiredRequest(ctx, req.GetRequestID())
	// 	return false
	// }

	if rep.ResultStatus == types.ResultFailure {
		k.Logger(ctx).Error("result status is fail")
		return false
	}

	// Check if validator exists and active
	_, isExist := k.StakingKeeper.GetValidator(ctx, rep.GetReporter().GetValidator())
	if !isExist {
		k.Logger(ctx).Error(fmt.Sprintf("error in validating the report: validator does not exist"))
		return false
	}

	// // check if validator is bonded or not
	// if isBonded := validator.IsBonded(); !isBonded {
	// 	k.Logger(ctx).Error(fmt.Sprintf("error in validating the report: validator is not bonded, cannot send reports to receive rewards"))
	// 	return false
	// }

	// TODO
	err := k.ValidateReport(ctx, rep.GetReporter(), req)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error in validating the report: %v\n", err.Error()))
		return false
	}
	return true
}
