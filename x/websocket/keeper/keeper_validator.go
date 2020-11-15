package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/oraichain/orai/x/websocket/types"
)

// SetValidator saves the validator into the store
func (k Keeper) SetValidator(ctx sdk.Context, id string, rep types.Report) {
	ctx.KVStore(k.storeKey).Set(types.ReportStoreKey(string(rep.Reporter.Validator[:]), id), k.cdc.MustMarshalBinaryBare(rep))
}

// GetValidator return a specific validator given a validator address
func (k Keeper) GetValidator(ctx sdk.Context, valAddress sdk.ValAddress) staking.ValidatorI {
	return k.stakingKeeper.Validator(ctx, valAddress)
}

// AddValidator stores a list of validators to set a test case into the store
func (k Keeper) AddValidator(ctx sdk.Context, validator types.Validator) {

}
