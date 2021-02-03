package keeper

import (
	//"fmt"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AllocateTokens allocates the tokens to the validators that participate in the AI request handling
func (k Keeper) AllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)

	requestFees, creators := k.CollectRequestFees(ctx, ctx.BlockHeight()-int64(1))
	// if there are fees from the requests, we remove them from the fee collector
	if requestFees != nil && creators != nil {
		// since both share the same length, we only need to iterate one time through either array
		for i := range requestFees {
			// we substract fees from the creator based on the fees they provide
			err := k.bankKeeper.SubtractCoins(ctx, creators[i], requestFees[i])
			if err != nil {
				return
			}
		}
	}
	// get reward from the previous block
	rewardObj, err := k.GetReward(ctx, ctx.BlockHeight()-int64(1))
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
		fmt.Println("error adding coins using bank keeper: ", err)
		return
	}
	remaining := reward
	//Allocate non-community pool tokens to active validators weighted by voting power.

	// reward for test cases that contribute
	for _, testCase := range rewardObj.TestCases {
		// send coins to test case owner addresses
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, testCase.GetOwner(), testCase.GetFees())
		remaining = remaining.Sub(sdk.NewDecCoinsFromCoins(testCase.GetFees()...))
	}

	// reward for test cases that contribute
	for _, dataSource := range rewardObj.DataSources {
		// send coins to data source owner addresses
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, dataSource.GetOwner(), dataSource.GetFees())
		remaining = remaining.Sub(sdk.NewDecCoinsFromCoins(dataSource.GetFees()...))
	}
	// reward for the validators that contribute in the ai request test
	// transfer collected fees to the distribution module account to distribute the oracle rewards to the validators. Note that if we transfer all the transaction fees, then other modules won't be able to handle allocation

	// fix check division by zero
	if rewardObj.TotalPower <= int64(0) {
		return
	} else {
		for _, val := range rewardObj.Validators {
			powerFraction := sdk.NewDec(val.GetVotingPower()).QuoTruncate(sdk.NewDec(rewardObj.TotalPower))
			valRewardDec := sdk.NewDecCoinsFromCoins(rewardObj.ValidatorFees...).MulDec(powerFraction)
			valRewardInt, _ := valRewardDec.TruncateDecimal()
			err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, distr.ModuleName, valRewardInt)
			if err != nil {
				fmt.Println("error in sending coins from fee collector to distrution module: ", err)
				return
			}
			// allocate tokens to validator with a specific commission
			k.distrKeeper.AllocateTokensToValidator(ctx, k.stakingKeeper.Validator(ctx, val.GetAddress()), valRewardDec)
			remaining = remaining.Sub(valRewardDec)
		}
	}
	fmt.Println("Finish allocating the tokens")
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
