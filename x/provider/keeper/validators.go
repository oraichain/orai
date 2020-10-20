package keeper

import (
	//"fmt"

	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/oraichain/orai/packages/rng"
	"github.com/oraichain/orai/x/provider/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AllocateTokens allocates the tokens to the validators that participate in the AI request handling
func (k Keeper) AllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo) {

	//logger := k.Logger(ctx)

	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := feeCollector.GetCoins()
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// If there are no fees, we do not need to handle anything
	if feesCollected.Empty() {
		return
	}

	// get reward from the previous block
	rewardObj, err := k.GetReward(ctx, ctx.BlockHeight()-int64(1))

	// If there's no reward in the previous block, then we do not handle
	if err != nil || rewardObj.BlockHeight == int64(-1) {
		return
	}

	fmt.Println("Ready to allocate tokens with amount: ", feesCollected.String())

	// Compute the fee allocated for oracle module to distribute to active validators.
	rewardRatio := sdk.NewDecWithPrec(int64(k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)), 2)
	rewardInt, _ := feesCollected.MulDecTruncate(rewardRatio).TruncateDecimal()

	// transfer collected fees to the distribution module account. Note that if we transfer all the transaction fees, then other modules won't be able to handle allocation
	err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, distr.ModuleName, rewardInt)
	if err != nil {
		panic(err)
	}

	// Convert the transfered tokens back to DecCoins for internal distr allocations.
	reward := sdk.NewDecCoinsFromCoins(rewardInt...)
	remaining := reward
	rewardMultiplier := sdk.OneDec().Sub(k.DistrKeeper.GetCommunityTax(ctx)).Sub(k.DistrKeeper.GetCommunityTax(ctx))
	//Allocate non-community pool tokens to active validators weighted by voting power.

	// reward for the validators that contribute in the ai request test
	for _, val := range rewardObj.Validators {
		powerFraction := sdk.NewDec(val.VotingPower).QuoTruncate(sdk.NewDec(rewardObj.TotalPower))
		finalReward := reward.MulDecTruncate(rewardMultiplier).MulDecTruncate(powerFraction)
		// allocate tokens to validator with a specific commission
		k.DistrKeeper.AllocateTokensToValidator(ctx, k.stakingKeeper.Validator(ctx, val.Address), finalReward)
		remaining = remaining.Sub(finalReward)
	}

	// reward for test cases that contribute
	for _, owner := range rewardObj.TestCaseOwners {
		powerFraction := sdk.NewDec(1).QuoTruncate(sdk.NewDec(100))
		finalReward := reward.MulDecTruncate(rewardMultiplier).MulDecTruncate(powerFraction)
		rewardCoins, _ := finalReward.TruncateDecimal()
		// send coins to test case owner addresses
		k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, owner, rewardCoins)
		remaining = remaining.Sub(finalReward)
	}

	// reward for test cases that contribute
	for _, owner := range rewardObj.DataSourceOwners {
		powerFraction := sdk.NewDec(1).QuoTruncate(sdk.NewDec(100))
		finalReward := reward.MulDecTruncate(rewardMultiplier).MulDecTruncate(powerFraction)
		rewardCoins, _ := finalReward.TruncateDecimal()
		// send coins to data source owner addresses
		k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, owner, rewardCoins)
		remaining = remaining.Sub(finalReward)
	}

	// Allocate the remaining coins to the community pool.
	feePool := k.DistrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.DistrKeeper.SetFeePool(ctx, feePool)

	fmt.Println("Finish allocating the tokens")
}

// DirectAllocateTokens allocates the tokens to the validators, data sources and test cases that participate in the AI request handling directly using coins from the requester account
func (k Keeper) DirectAllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo) {
	reports := k.GetReportsBlockHeight(ctx, ctx.BlockHeight()-int64(1))

	for _, report := range reports {
		request, err := k.GetAIRequest(ctx, report.RequestID)
		if err != nil {
			return
		}
		// at this stage, the validator has run all the test cases and data sources to produce a valid report. This must be checked before => assume we have enough
		for _, dSource := range request.AIDataSources {
			// the creator will directly pays the data source provider
			k.bankKeeper.SendCoins(ctx, request.Creator, dSource.Owner, dSource.Fees)
		}

		for _, testCase := range request.TestCases {
			// the creator will directly pays the data source provider
			k.bankKeeper.SendCoins(ctx, request.Creator, testCase.Owner, testCase.Fees)
		}

		// Allocate tokens directly to the validator that creates a report using the fees given in the report
		k.DistrKeeper.AllocateTokensToValidator(ctx, k.stakingKeeper.Validator(ctx, report.Validator.Address), sdk.NewDecCoinsFromCoins(report.Fees...))
	}
}

// // RandomValidators random a set of validators (currently not based on the voting power) to execute the oracle script
// func (k Keeper) RandomValidators(ctx sdk.Context, size int) ([]sdk.ValAddress, error) {
// 	valOperators := []sdk.ValAddress{}
// 	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
// 		func(idx int64, val exported.ValidatorI) (stop bool) {
// 			valOperators = append(valOperators, val.GetOperator())
// 			return false
// 		})
// 	if len(valOperators) < size {
// 		return nil, sdkerrors.Wrapf(
// 			types.ErrNotEnoughValidators, "%d < %d", len(valOperators), size)
// 	}

// 	validators := make([]sdk.ValAddress, size)
// 	for i := 0; i < size; i++ {
// 		validators[i] = valOperators[rand.Intn(size)]
// 	}
// 	return validators, nil
// }

// RandomValidators random a set of validators (currently not based on the voting power) to execute the oracle script
func (k Keeper) RandomValidators(ctx sdk.Context, size int, nonce []byte) ([]sdk.ValAddress, error) {
	var curVotingP int64
	var prevVotingP int64
	specialIndex := 0 // this index stores the first validator that has equal index to the next val
	valOperators := []sdk.ValAddress{}
	maxValidatorSize := 0
	totalPowers := int64(0)
	validators := make([]sdk.ValAddress, size)
	// store a mapping of validators that have already been chosen
	chosenVal := make(map[string]string)
	// count the total current validator
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
		func(idx int64, val exported.ValidatorI) (stop bool) {
			// the highest staked validator has the highest freq appearance in the list. When random => higher chance of getting picked
			maxValidatorSize++
			totalPowers += val.GetConsensusPower()
			return false
		})
	// if there is no voting power, we return error to prevent x % 0 sampling
	if totalPowers == int64(0) {
		return nil, sdkerrors.Wrapf(
			types.ErrValidatorsHaveNoVotes, "%d < %d", len(valOperators), size)
	} else if maxValidatorSize < size {
		fmt.Println("not enough validators")
		return nil, sdkerrors.Wrapf(
			types.ErrNotEnoughValidators, "%d < %d", maxValidatorSize, size)
	} else {
		fmt.Println("enough validators")
		k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
			func(idx int64, val exported.ValidatorI) (stop bool) {
				// store the prev voting power validator
				prevVotingP = curVotingP
				// collect the new voting power
				curVotingP = val.GetConsensusPower()

				// if we meet the equal sistuation the first time then we note down the index
				if prevVotingP == curVotingP {
					// increment the index by one to make up for the index loss of the current validator
					specialIndex++
				} else {
					// reset the index to 0 since the sequence has ended
					specialIndex = 0
				}

				// the highest staked validator has the highest freq appearance in the list. When random => higher chance of getting picked
				for i := 0; i < maxValidatorSize+specialIndex; i++ {
					valOperators = append(valOperators, val.GetOperator())
				}
				maxValidatorSize--
				return false
			})

		fmt.Println("valOperators list: ", valOperators)

		randomGenerator, err := rng.NewRng(k.GetRngSeed(ctx), nonce, []byte(ctx.ChainID()))
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrSeedinitiation, err.Error())
		}
		for i := 0; i < size; i++ {
			// the dividend is randomed to make sure no one can predict the next validator
			dividend := randomGenerator.RandUint64()
			divisor := uint64(totalPowers)
			// this value init makes sure that we at least calculate the modulo once
			quotient := uint64(len(valOperators))
			// we keep calculating the new modulo until we get in the range
			for quotient >= uint64(len(valOperators)) {
				quotient, err = calucateMol(dividend, divisor)
				if err != nil {
					return nil, err
				}
				dividend = divisor
				divisor = quotient
			}
			// if the quotient is in the sampling list, and it is not in the chosen validator map range then we pick it
			valStr := valOperators[quotient].String()
			fmt.Println("chosen val: ", chosenVal[valStr])
			fmt.Println("chosen val has existed ? ", chosenVal[valStr] != valStr)
			if chosenVal[valStr] != valStr {
				// add to the chosen validator list
				chosenVal[valStr] = valStr
				validators[i] = valOperators[quotient]
			} else {
				// if it has been chosen already, we decrement the loop index to continue choosing a new one
				i--
			}
			fmt.Println("validator: ", validators[i].String())
			fmt.Println("All validators: ", validators)
		}
	}
	return validators, nil
}

// // GetValidatorFees calculates the total fees needed for a set of provided validators
// func (k Keeper) GetValidatorFees(ctx sdk.Context, providedCoins sdk.DecCoins, validators []sdk.ValAddress) (sdk.Coins, error) {
// 	for _, validator := range validators {
// 		power := k.GetValidator(ctx, validator).GetConsensusPower()
// 	}
// }

func calucateMol(dividend, divisor uint64) (uint64, error) {
	dividendBig := new(big.Int)
	dividendBig.SetUint64(dividend)
	divisorBig := new(big.Int)
	divisorBig.SetUint64(divisor)

	tenmodfour := new(big.Int)

	quotient := tenmodfour.Mod(dividendBig, divisorBig)
	return quotient.Uint64(), nil
}

// SetValidator saves the validator into the store
func (k Keeper) SetValidator(ctx sdk.Context, id string, rep types.Report) {
	ctx.KVStore(k.storeKey).Set(types.ReportStoreKey(string(rep.Validator.Address[:]), id), k.cdc.MustMarshalBinaryBare(rep))
}

// GetValidator return a specific validator given a validator address
func (k Keeper) GetValidator(ctx sdk.Context, valAddress sdk.ValAddress) staking.ValidatorI {
	return k.stakingKeeper.Validator(ctx, valAddress)
}

// AddValidator stores a list of validators to set a test case into the store
func (k Keeper) AddValidator(ctx sdk.Context, validator types.Validator) {

}
