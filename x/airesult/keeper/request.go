package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	aiRequest "github.com/oraichain/orai/x/airequest/exported"
	"github.com/oraichain/orai/x/airesult/types"
	webSocket "github.com/oraichain/orai/x/websocket/exported"
)

// GetRequestsBlockHeight returns all requests for the given block height, or nil if there is none.
func (k Keeper) GetRequestsBlockHeight(ctx sdk.Context, blockHeight int64) (reqs []aiRequest.AIRequestI) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.RequeststoreKeyPrefixAll())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var req aiRequest.AIRequestI
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &req)
		// check if block height is equal or not
		if req.GetBlockHeight() == blockHeight {
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
		fees = fees.Add(request.GetFees()...)
	}
	return fees
}

// ResolveExpiredRequest handles requests that have been expired
func (k Keeper) ResolveExpiredRequest(ctx sdk.Context, reqID string) {
	// hard code the result first if the request does not have a result
	if !k.HasResult(ctx, reqID) {
		valResults := make([]webSocket.ValResultI, 0)
		k.SetResult(ctx, reqID, types.NewAIRequestResult(reqID, valResults, types.RequestStatusExpired))
	} else {
		// if already has result then we change the request status to expired
		result, _ := k.GetResult(ctx, reqID)
		result.Status = types.RequestStatusExpired
		k.SetResult(ctx, reqID, result)
	}
}

// ResolveRequestsFromReports handles the reports received in a block to group all the validators, data source owners and test case owners
func (k Keeper) ResolveRequestsFromReports(ctx sdk.Context, rep webSocket.ReportI, reward *types.Reward, blockHeight int64) {

	req, _ := k.aiRequestKeeper.GetAIRequest(ctx, rep.GetRequestID())
	validation := k.validateBasic(ctx, req, rep, blockHeight)
	// if the report cannot pass the validation basic then we skip the rest
	if !validation {
		return
	}

	// collect data source owners that have their data sources executed to reward
	for _, dataSourceResult := range rep.GetDataSourceResults() {
		if dataSourceResult.GetStatus() == k.webSocketKeeper.GetKeyResultSuccess() {
			dataSource, _ := k.ProviderKeeper.GetAIDataSourceI(ctx, dataSourceResult.GetName())
			reward.DataSources = append(reward.DataSources, dataSource)
			reward.ProviderFees = reward.ProviderFees.Add(dataSource.GetFees()...)
		}
	}

	// collect data source owners that have their data sources executed to reward
	for _, testCaseResult := range rep.GetTestCaseResults() {
		testCase, _ := k.ProviderKeeper.GetTestCaseI(ctx, testCaseResult.GetName())
		reward.TestCases = append(reward.TestCases, testCase)
		reward.ProviderFees = reward.ProviderFees.Add(testCase.GetFees()...)
	}
	// calculate validator fees from the total fees (50% for data sources, test cases, 20% for list of validators)
	valFees, _ := sdk.NewDecCoinsFromCoins(reward.ProviderFees...).MulDec(sdk.NewDecWithPrec(int64(40), 2)).TruncateDecimal()
	// add validator fees into the total fees of all validators
	reward.ValidatorFees = reward.ValidatorFees.Add(valFees...)
	// store information into the reward struct to reward these entities in the next begin block
	valAddress := rep.GetValidator()
	validator := k.webSocketKeeper.NewValidator(valAddress, k.stakingKeeper.Validator(ctx, valAddress).GetConsensusPower(), "active")
	reward.Validators = append(reward.Validators, validator)
	reward.TotalPower += validator.GetVotingPower()
	// Aggregate the result and store it into the blockchain
	k.ResolveResult(ctx, req, rep)
}

func (k Keeper) validateBasic(ctx sdk.Context, req aiRequest.AIRequestI, rep webSocket.ReportI, blockHeight int64) bool {
	// if the request has been expired
	// if req.GetBlockHeight()+int64(k.GetParam(ctx, types.KeyExpirationCount)) < blockHeight {
	// 	//TODO: NEED TO HANDLE THE EXPIRED REQUEST.
	// 	fmt.Println("Request has been expired")
	// 	k.ResolveExpiredRequest(ctx, req.GetRequestID())
	// 	return false
	// }

	// Count the total number of data source results to see if it matches the requested data sources
	if len(rep.GetDataSourceResults()) != len(req.GetAIDataSources()) {
		fmt.Println("data source result length is different")
		return false
	}

	// Count the total number of test case results to see if it matches the requested test cases
	if len(rep.GetTestCaseResults()) != len(req.GetTestCases()) {
		fmt.Println("test case result length is different")
		return false
	}

	// TODO
	err := k.webSocketKeeper.ValidateReport(ctx, rep, req)
	if err != nil {
		fmt.Println("error in validating the report: ", err.Error())
		return false
	}
	return true
}
