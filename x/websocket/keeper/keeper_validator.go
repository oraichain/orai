package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/oraichain/orai/x/websocket/types"
)

// GetValidator return a specific validator given a validator address
func (k Keeper) GetValidator(ctx sdk.Context, valAddress sdk.ValAddress) staking.ValidatorI {
	return k.stakingKeeper.Validator(ctx, valAddress)
}

// NewValidator is a wrapper function of the websocket module that allow others to initiate a new validator entity through the keeper
func (k Keeper) NewValidator(address sdk.ValAddress, votingPower int64, status string) *types.Validator {
	return &types.Validator{
		Address:     address,
		VotingPower: votingPower,
		Status:      status,
	}
}
