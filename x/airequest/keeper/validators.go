package keeper

import (
	//"fmt"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/oraichain/orai/packages/rng"
	"github.com/oraichain/orai/x/airequest/types"
)

// RandomValidators random a set of validators (currently not based on the voting power) to execute the oracle script
func (k Keeper) RandomValidators(ctx sdk.Context, size int, nonce []byte) ([]sdk.ValAddress, error) {
	maxValidatorSize := 0
	totalPowers := int64(0)
	// count the total current validator
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(idx int64, val staking.ValidatorI) (stop bool) {
		// the highest staked validator has the highest freq appearance in the list. When random => higher chance of getting picked
		maxValidatorSize++
		totalPowers += val.GetConsensusPower()
		return false
	})
	// if there is no voting power, we return error to prevent x % 0 sampling
	if totalPowers == int64(0) {
		return nil, sdkerrors.Wrapf(types.ErrValidatorsHaveNoVotes, "%d < %d", maxValidatorSize, size)
	} else if maxValidatorSize < size {
		k.Logger(ctx).Error("not enough validators")
		return nil, sdkerrors.Wrapf(types.ErrNotEnoughValidators, "%d < %d", maxValidatorSize, size)
	} else {
		k.Logger(ctx).Info("enough validators")
		valOperators := k.createValSamplingList(ctx, maxValidatorSize)
		k.Logger(ctx).Error(fmt.Sprintf("AI request validator operators: %v\n", valOperators))
		randomGenerator, err := rng.NewRng(k.GetRngSeed(ctx), nonce, []byte(ctx.ChainID()))
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrSeedinitiation, err.Error())
		}
		validators := k.SampleIndexes(valOperators, size, randomGenerator, totalPowers)
		k.Logger(ctx).Error(fmt.Sprintf("AI request validator list final: %v\n", validators))
		return validators, nil
	}
}

// func calucateMol(dividend, divisor uint64) uint64 {
// 	dividendBig := new(big.Int)
// 	dividendBig.SetUint64(dividend)
// 	divisorBig := new(big.Int)

// 	// check division by zero or negative
// 	if divisor <= uint64(0) {
// 		// fix divisor to 1 to prevent division by zero
// 		divisorBig.SetInt64(1)
// 	} else {
// 		divisorBig.SetUint64(divisor)
// 	}

// 	tenmodfour := new(big.Int)

// 	quotient := tenmodfour.Mod(dividendBig, divisorBig)
// 	return quotient.Uint64()
// }

func (k Keeper) createValSamplingList(ctx sdk.Context, maxValidatorSize int) (valOperators []sdk.ValAddress) {
	var curVotingP int64
	var prevVotingP int64
	specialIndex := 0 // this index stores the first validator that has equal index to the next val
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(idx int64, val staking.ValidatorI) (stop bool) {
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

// SampleIndexes return random of indexes of chosen validators
func (k *Keeper) SampleIndexes(valOperators []sdk.ValAddress, size int, randomGenerator *rng.Rng, totalPowers int64) []sdk.ValAddress {

	valOperatorLen := uint64(len(valOperators))
	validators := make([]sdk.ValAddress, size)

	// store a mapping of validators that have already been chosen
	chosenVal := make([]bool, valOperatorLen)
	// time := 0
	for i := 0; i < size; {
		// the dividend is randomed to make sure no one can predict the next validator
		quotient := randomGenerator.RandUint64() % valOperatorLen
		// time++
		// fmt.Printf("%d) quotient :%v\n", time, quotient)
		// if the quotient is in the sampling list, and it is not in the chosen validator map range then we pick it
		if quotient == valOperatorLen {
			quotient--
		}
		if !chosenVal[quotient] {
			// add to the chosen validator list
			chosenVal[quotient] = true
			validators[i] = valOperators[quotient]
			i++
		}
	}
	return validators
}
