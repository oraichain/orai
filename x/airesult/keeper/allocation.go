package keeper

import (
	//"fmt"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/oraichain/orai/x/provider"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AllocateTokens allocates the tokens to the validators that participate in the AI request handling
func (k Keeper) AllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo, blockHeight int64) {
	// get reward from the previous block
	rewardObj, err := k.GetReward(ctx, blockHeight-1)
	// If there's no reward in the previous block, then we do not handle
	if err != nil || rewardObj.BlockHeight == int64(-1) {
		return
	}

	// retrieve fee collector module account to prepare token allocation1
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	// add all the fees from the report since we only reward those included in the report
	feesCollected := rewardObj.ProviderFees.Add(rewardObj.ValidatorFees...)
	reward := sdk.NewDecCoinsFromCoins(feesCollected...)
	// append those coins into the fee collector to get ready allocating them to the distr module.
	err = k.bankKeeper.AddCoins(ctx, feeCollector.GetAddress(), feesCollected)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error adding coins using bank keeper: %v\n", err.Error()))
		return
	}
	remaining := reward
	hasNeg := false

	//Allocate non-community pool tokens to active validators weighted by voting power.
	// reward for test cases that contribute
	for _, testCase := range rewardObj.TestCases {

		// safesub to prevent panic
		remaining, hasNeg = remaining.SafeSub(sdk.NewDecCoinsFromCoins(testCase.GetFees()...))
		if hasNeg {
			k.Logger(ctx).Error(fmt.Sprintf("not enough balance to reward test case :%v, \n", testCase.Name))
			return
		}

		// send coins to test case owner addresses
		temp := k.bankKeeper.GetBalance(ctx, testCase.GetOwner(), provider.Denom)
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, testCase.GetOwner(), testCase.GetFees())
		rewardCollected := k.bankKeeper.GetBalance(ctx, testCase.GetOwner(), provider.Denom).Sub(temp)
		k.Logger(ctx).Info(fmt.Sprintf("Reward collected for the following address %v - %v\n", testCase.GetOwner().String(), rewardCollected))
	}

	// reward for test cases that contribute
	for _, dataSource := range rewardObj.DataSources {

		// safesub to prevent panic
		remaining, hasNeg = remaining.SafeSub(sdk.NewDecCoinsFromCoins(dataSource.GetFees()...))
		if hasNeg {
			k.Logger(ctx).Error(fmt.Sprintf("not enough balance to reward data source :%v, \n", dataSource.Name))
			return
		}

		// send coins to data source owner addresses
		temp := k.bankKeeper.GetBalance(ctx, dataSource.GetOwner(), provider.Denom)
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, dataSource.GetOwner(), dataSource.GetFees())
		rewardCollected := k.bankKeeper.GetBalance(ctx, dataSource.GetOwner(), provider.Denom).Sub(temp)
		k.Logger(ctx).Info(fmt.Sprintf("Reward collected for the following address %v - %v\n", dataSource.GetOwner().String(), rewardCollected))

	}
	// reward for the validators that contribute in the ai request test
	// transfer collected fees to the distribution module account to distribute the oracle rewards to the validators. Note that if we transfer all the transaction fees, then other modules won't be able to handle allocation

	decValLen := sdk.NewDec(int64(len(rewardObj.Validators)))
	decTotalPower := sdk.NewDec(rewardObj.TotalPower)

	// fix check division by zero, no validator or zero total power
	if decValLen.IsZero() || decValLen.IsZero() {
		k.Logger(ctx).Error(fmt.Sprintf("total power zero\n"))
		return
	}

	for _, val := range rewardObj.Validators {
		powerFraction := sdk.NewDec(val.GetVotingPower()).QuoTruncate(decTotalPower)
		// since validator fees here is the sum of all validator fees, so we need to divide with total number of validators to get fees for one validator.
		valRewardDec := sdk.NewDecCoinsFromCoins(rewardObj.ValidatorFees...).QuoDec(decValLen).MulDec(powerFraction)

		// safesub to prevent panic
		remaining, hasNeg = remaining.SafeSub(valRewardDec)
		if hasNeg {
			k.Logger(ctx).Error(fmt.Sprintf("not enough balance to reward validator :%v, \n", val.GetAddress()))
			return
		}

		valRewardInt, _ := valRewardDec.TruncateDecimal()
		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, distr.ModuleName, valRewardInt)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("error in sending coins from fee collector to distrution module: %v\n", err.Error()))
			return
		}
		// allocate tokens to validator with a specific commission
		k.distrKeeper.AllocateTokensToValidator(ctx, k.stakingKeeper.Validator(ctx, val.GetAddress()), valRewardDec)
		k.Logger(ctx).Info(fmt.Sprintf("outstanding reward of validator %v - %v\n", val.GetAddress().String(), k.distrKeeper.GetValidatorAccumulatedCommission(ctx, val.GetAddress())))

	}

	// allocate community funding
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.distrKeeper.SetFeePool(ctx, feePool)
	k.Logger(ctx).Info("finish allocating tokens")
}

// // DirectAllocateTokens allocates the tokens to the validators, data sources and test cases that participate in the AI request handling directly using coins from the requester account
// func (k Keeper) DirectAllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo) {
// 	reports := k.webSocketKeeper.GetReportsBlockHeight(ctx, ctx.BlockHeight()-int64(1))

// 	// TODO: instead of directly allocating tokens like this which is insecure, we get the total fees and then put it back to the fee collector. By doing this, we make sure that the total fee is enough to allocate for all using the fee collector.
// 	for _, report := range reports {
// 		request, err := k.GetAIRequest(ctx, report.GetRequestID())
// 		if err != nil {
// 			return
// 		}
// 		// at this stage, the validator has run all the test cases and data sources to produce a valid report. This must be checked before => assume we have enough
// 		for _, dSource := range request.AIDataSources {
// 			// the creator will directly pays the data source provider
// 			k.bankKeeper.SendCoins(ctx, request.Creator, dSource.GetOwner(), dSource.GetFees())
// 		}

// 		for _, testCase := range request.TestCases {
// 			// the creator will directly pays the data source provider
// 			k.bankKeeper.SendCoins(ctx, request.Creator, testCase.GetOwner(), testCase.GetFees())
// 		}

// 		// Allocate tokens directly to the validator that creates a report using the fees given in the report
// 		k.distrKeeper.AllocateTokensToValidator(ctx, k.stakingKeeper.Validator(ctx, report.GetValidator()), sdk.NewDecCoinsFromCoins(report.GetFees()...))
// 	}
// }
