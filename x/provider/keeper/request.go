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
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.RequestStoreKey(id)), &result)
	if err != nil {
		return types.AIRequest{}, err
	}
	return result, nil
}

// SetAIRequest allows users to set a oScript into the store
func (k Keeper) SetAIRequest(ctx sdk.Context, id string, request types.AIRequest) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(types.RequestStoreKey(id), bz)
}

// GetAllAIRequestIDs get an iterator of all key-value pairs in the store
func (k Keeper) GetAllAIRequestIDs(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte("req"))
}

// ResolveExpiredRequest handles requests that have been expired
func (k Keeper) ResolveExpiredRequest(ctx sdk.Context, reqID string) {
	// hard code the result first if the request does not have a result
	if !k.HasResult(ctx, reqID) {
		k.SetResult(ctx, reqID, types.NewAIRequestResult(reqID, [][]byte{}, types.RequestStatusExpired))
	} else {
		// if already has result then we change the request status to expired
		result, _ := k.GetResult(ctx, reqID)
		result.Status = types.RequestStatusExpired
		k.SetResult(ctx, reqID, result)
	}
}

// ResolveRequest handles the reports received in a block to group all the validators, data source owners and test case owners
func (k Keeper) ResolveRequest(ctx sdk.Context, rep types.Report, blockHeight int64) ([]types.Validator, []sdk.AccAddress, []sdk.AccAddress, int64) {
	// fmt.Println("Param of the provider module: ", k.GetParam(ctx, types.KeyOracleScriptRewardPercentage))

	var validators []types.Validator
	var dataSourceOwners []sdk.AccAddress
	var testCaseOwners []sdk.AccAddress
	votingPower := int64(0)

	req, _ := k.GetAIRequest(ctx, rep.RequestID)

	// if the request has been expired
	if req.BlockHeight+int64(k.GetParam(ctx, types.KeyExpirationCount)) < blockHeight {
		//TODO: NEED TO HANDLE THE EXPIRED REQUEST.
		fmt.Println("Request has been expired")
		k.ResolveExpiredRequest(ctx, req.RequestID)
		return nil, nil, nil, votingPower
	}

	// Count the total number of data source results to see if it matches the requested data sources
	if len(rep.DataSourceResults) != len(req.AIDataSources) {
		fmt.Println("data source result length is different")
		return nil, nil, nil, votingPower
	}

	// Count the total number of test case results to see if it matches the requested test cases
	if len(rep.TestCaseResults) != len(req.TestCases) {
		fmt.Println("test case result length is different")
		return nil, nil, nil, votingPower
	}

	validators = append(validators, rep.Validator)
	votingPower = votingPower + rep.Validator.VotingPower

	// collect data source owners that have their data sources executed to reward
	for _, dataSourceResult := range rep.DataSourceResults {
		//TODO: one problem: what happens if an data source is called many times ? => reward one time or reward for each time it is called ?
		dataSource, _ := k.GetAIDataSource(ctx, dataSourceResult.Name)
		dataSourceOwners = append(dataSourceOwners, dataSource.Owner)
	}

	// collect data source owners that have their data sources executed to reward
	for _, testCaseResult := range rep.TestCaseResults {
		testCase, _ := k.GetTestCase(ctx, testCaseResult.Name)
		testCaseOwners = append(testCaseOwners, testCase.Owner)
	}

	// Aggregate the result and store it into the blockchain
	k.ResolveResult(ctx, req, rep)

	return validators, dataSourceOwners, testCaseOwners, votingPower
}
