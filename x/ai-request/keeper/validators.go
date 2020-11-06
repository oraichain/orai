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
	"github.com/oraichain/orai/x/ai-request/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// AllocateTokens allocates the tokens to the validators that participate in the AI request handling
func (k Keeper) AllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := feeCollector.GetCoins()
	// If there are no fees, we do not need to handle anything
	if feesCollectedInt.Empty() {
		return
	}
	requestFees := k.CollectRequestFees(ctx, ctx.BlockHeight()-int64(1))
	// if there are fees from the requests, we remove them from the fee collector
	if !requestFees.IsZero() {
		// 100 - 70 = 30%
		rewardRatio := sdk.NewDecWithPrec(int64(100)-int64(k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)), 2)
		rewardFeesDec := sdk.NewDecCoinsFromCoins(requestFees...)
		rewardFees, _ := rewardFeesDec.MulDecTruncate(rewardRatio).TruncateDecimal()
		// we remove the reward fees from the request fees to reward the proposer afterwards immediately
		feeCollector.SetCoins(feeCollector.GetCoins().Sub(rewardFees))
		// Workaround of the Cosmos sdk bug. When we set coin using feeCollector, the actual coins in bankKeeper does not update. When sending coins to other accounts, it uses coins from the bankKeeper interface => we need to do another step to completely remove coins.
		_, err := k.bankKeeper.SubtractCoins(ctx, feeCollector.GetAddress(), rewardFees)
		if err != nil {
			return
		}
	}
	// get reward from the previous block
	rewardObj, err := k.GetReward(ctx, ctx.BlockHeight()-int64(1))
	// If there's no reward in the previous block, then we do not handle
	if err != nil || rewardObj.BlockHeight == int64(-1) {
		return
	}
	// add all the fees from the report since we only reward those included in the report
	feesCollectedInt = feesCollectedInt.Add(rewardObj.ProviderFees...).Add(rewardObj.ValidatorFees...)
	reward := sdk.NewDecCoinsFromCoins(feesCollectedInt...)
	// append those coins into the fee collector to get ready allocating them to the distr module.
	err = feeCollector.SetCoins(feeCollector.GetCoins().Add(feesCollectedInt...))
	if err != nil {
		fmt.Println("error set coins into fee collector: ", err)
		return
	}
	_, err = k.bankKeeper.AddCoins(ctx, feeCollector.GetAddress(), feesCollectedInt)
	if err != nil {
		fmt.Println("error adding coins using bank keeper: ", err)
		return
	}
	remaining := reward
	//Allocate non-community pool tokens to active validators weighted by voting power.

	// reward for test cases that contribute
	for _, testCase := range rewardObj.TestCases {
		// send coins to test case owner addresses
		k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, testCase.GetOwner(), testCase.GetFees())
		remaining = remaining.Sub(sdk.NewDecCoinsFromCoins(testCase.GetFees()...))
	}

	// reward for test cases that contribute
	for _, dataSource := range rewardObj.DataSources {
		// send coins to data source owner addresses
		k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, dataSource.GetOwner(), dataSource.GetFees())
		remaining = remaining.Sub(sdk.NewDecCoinsFromCoins(dataSource.GetFees()...))
	}
	// reward for the validators that contribute in the ai request test
	// transfer collected fees to the distribution module account to distribute the oracle rewards to the validators. Note that if we transfer all the transaction fees, then other modules won't be able to handle allocation
	for _, val := range rewardObj.Validators {
		powerFraction := sdk.NewDec(val.VotingPower).QuoTruncate(sdk.NewDec(rewardObj.TotalPower))
		valRewardDec := sdk.NewDecCoinsFromCoins(rewardObj.ValidatorFees...).MulDec(powerFraction)
		valRewardInt, _ := valRewardDec.TruncateDecimal()
		err = k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, distr.ModuleName, valRewardInt)
		if err != nil {
			fmt.Println("error in sending coins from fee collector to distrution module: ", err)
			return
		}
		// allocate tokens to validator with a specific commission
		k.distrKeeper.AllocateTokensToValidator(ctx, k.stakingKeeper.Validator(ctx, val.Address), valRewardDec)
		remaining = remaining.Sub(valRewardDec)
	}
	fmt.Println("Finish allocating the tokens")
}

// DirectAllocateTokens allocates the tokens to the validators, data sources and test cases that participate in the AI request handling directly using coins from the requester account
func (k Keeper) DirectAllocateTokens(ctx sdk.Context, prevVotes []abci.VoteInfo) {
	reports := k.GetReportsBlockHeight(ctx, ctx.BlockHeight()-int64(1))

	// TODO: instead of directly allocating tokens like this which is insecure, we get the total fees and then put it back to the fee collector. By doing this, we make sure that the total fee is enough to allocate for all using the fee collector.
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
		k.distrKeeper.AllocateTokensToValidator(ctx, k.stakingKeeper.Validator(ctx, report.Validator.Address), sdk.NewDecCoinsFromCoins(report.Fees...))
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
	maxValidatorSize := 0
	totalPowers := int64(0)
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
			types.ErrValidatorsHaveNoVotes, "%d < %d", maxValidatorSize, size)
	} else if maxValidatorSize < size {
		fmt.Println("not enough validators")
		return nil, sdkerrors.Wrapf(
			types.ErrNotEnoughValidators, "%d < %d", maxValidatorSize, size)
	} else {
		fmt.Println("enough validators")
		valOperators := k.createValSamplingList(ctx, maxValidatorSize)
		fmt.Println("valOperators list: ", valOperators)

		randomGenerator, err := rng.NewRng(k.GetRngSeed(ctx), nonce, []byte(ctx.ChainID()))
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrSeedinitiation, err.Error())
		}
		validators := k.sampleIndexes(valOperators, size, randomGenerator, totalPowers)
		return validators, nil
	}
}

func calucateMol(dividend, divisor uint64) uint64 {
	dividendBig := new(big.Int)
	dividendBig.SetUint64(dividend)
	divisorBig := new(big.Int)
	divisorBig.SetUint64(divisor)

	tenmodfour := new(big.Int)

	quotient := tenmodfour.Mod(dividendBig, divisorBig)
	return quotient.Uint64()
}

func (k Keeper) createValSamplingList(ctx sdk.Context, maxValidatorSize int) (valOperators []sdk.ValAddress) {
	var curVotingP int64
	var prevVotingP int64
	specialIndex := 0 // this index stores the first validator that has equal index to the next val
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
	return valOperators
}

func (k Keeper) sampleIndexes(valOperators []sdk.ValAddress, size int, randomGenerator *rng.Rng, totalPowers int64) []sdk.ValAddress {
	// store a mapping of validators that have already been chosen
	chosenVal := make(map[string]string)
	validators := make([]sdk.ValAddress, size)
	for i := 0; i < size; i++ {
		// the dividend is randomed to make sure no one can predict the next validator
		dividend := randomGenerator.RandUint64()
		divisor := uint64(totalPowers)
		// this value init makes sure that we at least calculate the modulo once
		quotient := uint64(len(valOperators))
		// we keep calculating the new modulo until we get in the range
		for quotient >= uint64(len(valOperators)) {
			quotient = calucateMol(dividend, divisor)
			dividend = divisor
			divisor = quotient
		}
		// if the quotient is in the sampling list, and it is not in the chosen validator map range then we pick it
		valStr := valOperators[quotient].String()
		if chosenVal[valStr] != valStr {
			// add to the chosen validator list
			chosenVal[valStr] = valStr
			validators[i] = valOperators[quotient]
		} else {
			// if it has been chosen already, we decrement the loop index to continue choosing a new one
			i--
		}
		fmt.Println("All validators: ", validators)
	}
	return validators
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
