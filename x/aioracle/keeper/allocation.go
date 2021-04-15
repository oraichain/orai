package keeper

import (
	//"fmt"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/oraichain/orai/x/aioracle/types"
	aioracle "github.com/oraichain/orai/x/aioracle/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AllocateTokens allocates the tokens to the validators that participate in the AI request handling
func (k *Keeper) AllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo, blockHeight int64) {
	// get reward from the previous block
	rewardObj, err := k.GetReward(ctx, blockHeight-1)
	// If there's no reward in the previous block, then we do not handle
	if err != nil || rewardObj.BaseReward.BlockHeight == int64(-1) {
		return
	}

	// retrieve fee collector module account to prepare token allocation1
	feeCollector := k.AuthKeeper.GetModuleAccount(ctx, k.FeeCollectorName)
	// add all the fees from the report since we only reward those included in the report
	feesCollected := rewardObj.BaseReward.ProviderFees.Add(rewardObj.BaseReward.ValidatorFees...)
	reward := sdk.NewDecCoinsFromCoins(feesCollected...)
	// append those coins into the fee collector to get ready allocating them to the distr module.
	err = k.BankKeeper.AddCoins(ctx, feeCollector.GetAddress(), feesCollected)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("error adding coins using bank *Keeper: %v\n", err.Error()))
		return
	}
	remaining := reward
	hasNeg := false

	// reward for data sources that contribute
	for _, result := range rewardObj.Results {
		providerFees := result.GetEntryPoint().GetProviderFees()
		owner := result.GetEntryPoint().GetOwner()
		// safesub to prevent panic
		remaining, hasNeg = remaining.SafeSub(sdk.NewDecCoinsFromCoins(providerFees...))
		if hasNeg {
			k.Logger(ctx).Error(fmt.Sprintf("not enough balance to reward provider with URL:%v, \n", result.GetEntryPoint().GetUrl()))
			return
		}

		// send coins to data source owner addresses
		temp := k.BankKeeper.GetBalance(ctx, owner, types.Denom)
		k.BankKeeper.SendCoinsFromModuleToAccount(ctx, k.FeeCollectorName, owner, providerFees)
		rewardCollected := k.BankKeeper.GetBalance(ctx, owner, aioracle.Denom).Sub(temp)
		k.Logger(ctx).Info(fmt.Sprintf("Reward collected for the following address %v - %v\n", owner.String(), rewardCollected))

	}
	//reward for the validators that contribute in the ai request test
	// transfer collected fees to the distribution module account to distribute the oracle rewards to the validators. Note that if we transfer all the transaction fees, then other modules won't be able to handle allocation

	decValLen := sdk.NewDec(int64(len(rewardObj.BaseReward.Validators)))
	decTotalPower := sdk.NewDec(rewardObj.BaseReward.TotalPower)

	// fix check division by zero, no validator or zero total power
	if decValLen.IsZero() || decValLen.IsZero() {
		k.Logger(ctx).Error(fmt.Sprintf("total power zero\n"))
		return
	}

	for _, val := range rewardObj.BaseReward.Validators {
		powerFraction := sdk.NewDec(val.GetVotingPower()).QuoTruncate(decTotalPower)
		// since validator fees here is the sum of all validator fees, so we need to divide with total number of validators to get fees for one validator.
		valRewardDec := sdk.NewDecCoinsFromCoins(rewardObj.BaseReward.ValidatorFees...).QuoDec(decValLen).MulDec(powerFraction)

		// safesub to prevent panic
		remaining, hasNeg = remaining.SafeSub(valRewardDec)
		if hasNeg {
			k.Logger(ctx).Error(fmt.Sprintf("not enough balance to reward validator :%v, \n", val.GetAddress()))
			return
		}

		valRewardInt, _ := valRewardDec.TruncateDecimal()
		err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, k.FeeCollectorName, distr.ModuleName, valRewardInt)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("error in sending coins from fee collector to distrution module: %v\n", err.Error()))
			return
		}
		// allocate tokens to validator with a specific commission
		k.DistrKeeper.AllocateTokensToValidator(ctx, k.StakingKeeper.Validator(ctx, val.GetAddress()), valRewardDec)
		if valRewardInt.Empty() {
			valRewardInt = sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0))))
		}
		k.Logger(ctx).Info(fmt.Sprintf("outstanding reward of validator %v - %v\n", val.GetAddress().String(), valRewardInt))

	}

	// allocate community funding
	feePool := k.DistrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.DistrKeeper.SetFeePool(ctx, feePool)
	k.Logger(ctx).Info("finish allocating tokens")
}
