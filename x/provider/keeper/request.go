package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

// GetAIRequest returns the information of an AI request
func (k Keeper) GetAIRequest(ctx sdk.Context, id string) (types.AIRequest, error) {
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

	bz := k.cdc.MustMarshalBinaryBare(request)
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

// ResolveExpiredRequest handles requests that have been expired
func (k Keeper) ResolveExpiredRequest(ctx sdk.Context, reqID string) {
	// hard code the result first if the request does not have a result
	if !k.HasResult(ctx, reqID) {
		k.SetResult(ctx, reqID, types.NewAIRequestResult(reqID, types.ValResults{}, types.RequestStatusExpired))
	} else {
		// if already has result then we change the request status to expired
		result, _ := k.GetResult(ctx, reqID)
		result.Status = types.RequestStatusExpired
		k.SetResult(ctx, reqID, result)
	}
}

// ResolveRequestsFromReports handles the reports received in a block to group all the validators, data source owners and test case owners
func (k Keeper) ResolveRequestsFromReports(ctx sdk.Context, rep types.Report, reward *types.Reward, blockHeight int64) {
	// fmt.Println("Param of the provider module: ", k.GetParam(ctx, types.KeyOracleScriptRewardPercentage))

	req, _ := k.GetAIRequest(ctx, rep.RequestID)
	validation := k.validateBasic(ctx, req, rep, blockHeight)
	// if the report cannot pass the validation basic then we skip the rest
	if !validation {
		return
	}

	// collect data source owners that have their data sources executed to reward
	for _, dataSourceResult := range rep.DataSourceResults {
		if dataSourceResult.Status == types.ResultSuccess {
			dataSource, _ := k.GetAIDataSource(ctx, dataSourceResult.Name)
			reward.DataSources = append(reward.DataSources, dataSource)
			reward.ProviderFees = reward.ProviderFees.Add(dataSource.Fees...)
		}
	}

	// collect data source owners that have their data sources executed to reward
	for _, testCaseResult := range rep.TestCaseResults {
		testCase, _ := k.GetTestCase(ctx, testCaseResult.Name)
		reward.TestCases = append(reward.TestCases, testCase)
		reward.ProviderFees = reward.ProviderFees.Add(testCase.Fees...)
	}
	// calculate validator fees from the total fees (50% for data sources, test cases, 20% for list of validators)
	valFees, _ := sdk.NewDecCoinsFromCoins(reward.ProviderFees...).MulDec(sdk.NewDecWithPrec(int64(40), 2)).TruncateDecimal()
	// add validator fees into the total fees of all validators
	reward.ValidatorFees = reward.ValidatorFees.Add(valFees...)
	// store information into the reward struct to reward these entities in the next begin block
	reward.Validators = append(reward.Validators, rep.Validator)
	reward.TotalPower += rep.Validator.VotingPower

	// Aggregate the result and store it into the blockchain
	k.ResolveResult(ctx, req, rep)
}

func (k Keeper) validateBasic(ctx sdk.Context, req types.AIRequest, rep types.Report, blockHeight int64) bool {
	// if the request has been expired
	if req.BlockHeight+int64(k.GetParam(ctx, types.KeyExpirationCount)) < blockHeight {
		//TODO: NEED TO HANDLE THE EXPIRED REQUEST.
		fmt.Println("Request has been expired")
		k.ResolveExpiredRequest(ctx, req.RequestID)
		return false
	}

	// Count the total number of data source results to see if it matches the requested data sources
	if len(rep.DataSourceResults) != len(req.AIDataSources) {
		fmt.Println("data source result length is different")
		return false
	}

	// Count the total number of test case results to see if it matches the requested test cases
	if len(rep.TestCaseResults) != len(req.TestCases) {
		fmt.Println("test case result length is different")
		return false
	}
	return true
}
