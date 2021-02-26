package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	airequest "github.com/oraichain/orai/x/airequest/types"
	"github.com/oraichain/orai/x/airesult/types"
	websocket "github.com/oraichain/orai/x/websocket/types"
)

// ResolveExpiredRequest handles requests that have been expired
func (k Keeper) ResolveExpiredRequest(ctx sdk.Context, reqID string) {
	// hard code the result first if the request does not have a result
	if !k.HasResult(ctx, reqID) {
		valResults := make([]websocket.ValResult, 0)
		k.SetResult(ctx, reqID, types.NewAIRequestResult(reqID, valResults, types.RequestStatusExpired))
	} else {
		// if already has result then we change the request status to expired
		result, _ := k.GetResult(ctx, reqID)
		result.Status = types.RequestStatusExpired
		k.SetResult(ctx, reqID, result)
	}
}

// ResolveRequestsFromReports handles the reports received in a block to group all the validators, data source owners and test case owners
func (k Keeper) ResolveRequestsFromReports(ctx sdk.Context, rep *websocket.Report, reward *types.Reward, blockHeight int64, rewardPercentage int64) (bool, int) {

	req, _ := k.aiRequestKeeper.GetAIRequest(ctx, rep.GetRequestID())
	validation := k.validateBasic(ctx, req, rep, blockHeight)
	// if the report cannot pass the validation basic then we skip the rest
	if !validation {
		return false, 0
	}

	// collect data source owners that have their data sources executed to reward
	for _, dataSourceResult := range rep.GetDataSourceResults() {
		if dataSourceResult.GetStatus() == k.webSocketKeeper.GetKeyResultSuccess() {
			dataSource, _ := k.providerKeeper.GetAIDataSource(ctx, dataSourceResult.GetName())
			reward.DataSources = append(reward.DataSources, *dataSource)
			reward.ProviderFees = reward.ProviderFees.Add(dataSource.GetFees()...)
		}
	}

	// collect data source owners that have their data sources executed to reward
	for _, testCaseResult := range rep.GetTestCaseResults() {
		testCase, _ := k.providerKeeper.GetTestCase(ctx, testCaseResult.GetName())
		reward.TestCases = append(reward.TestCases, *testCase)
		reward.ProviderFees = reward.ProviderFees.Add(testCase.GetFees()...)
	}
	// change reward ratio to the ratio of validator
	// 0.6 by default, 2 decimals for percentage
	percentageDec := k.providerKeeper.GetPercentageDec(rewardPercentage)
	rewardRatio := sdk.NewDecWithPrec(int64(100), 2).Sub(percentageDec)

	// reward = 1 - oracle reward percentage × (data source fees + test case fees)
	valFees, _ := sdk.NewDecCoinsFromCoins(reward.ProviderFees...).MulDec(rewardRatio).TruncateDecimal()
	// add validator fees into the total fees of all validators
	reward.ValidatorFees = reward.ValidatorFees.Add(valFees...)
	// store information into the reward struct to reward these entities in the next begin block
	valAddress := rep.GetReporter().GetValidator()
	validator := k.webSocketKeeper.NewValidator(valAddress, k.stakingKeeper.Validator(ctx, valAddress).GetConsensusPower(), "active")
	reward.Validators = append(reward.Validators, *validator)
	reward.TotalPower += validator.GetVotingPower()

	// return boolean and length of validator list to resolve result
	return true, len(req.GetValidators())
}

func (k Keeper) validateBasic(ctx sdk.Context, req *airequest.AIRequest, rep *websocket.Report, blockHeight int64) bool {
	// if the request has been expired
	// if req.GetBlockHeight()+int64(k.GetParam(ctx, types.KeyExpirationCount)) < blockHeight {
	// 	//TODO: NEED TO HANDLE THE EXPIRED REQUEST.
	// 	fmt.Println("Request has been expired")
	// 	k.ResolveExpiredRequest(ctx, req.GetRequestID())
	// 	return false
	// }

	if rep.ResultStatus == websocket.ResultFailure {
		k.Logger(ctx).Error("result status is fail")
		return false
	}

	// Count the total number of data source results to see if it matches the requested data sources
	if len(rep.GetDataSourceResults()) != len(req.GetAiDataSources()) {
		k.Logger(ctx).Error("data source result length is different")
		return false
	}

	// Count the total number of test case results to see if it matches the requested test cases
	if len(rep.GetTestCaseResults()) != len(req.GetTestCases()) {
		k.Logger(ctx).Error("test case result length is different")
		return false
	}

	// TODO
	err := k.webSocketKeeper.ValidateReport(ctx, rep.GetReporter(), req)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error in validating the report: %v\n", err.Error()))
		return false
	}
	return true
}
