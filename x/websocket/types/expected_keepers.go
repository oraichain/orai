package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/exported"
)

// ParamSubspace defines the expected Subspace interfacace
// type ParamSubspace interface {
// 	WithKeyTable(table params.KeyTable) params.Subspace
// 	Get(ctx sdk.Context, key []byte, ptr interface{})
// 	GetParamSet(ctx sdk.Context, ps params.ParamSet)
// 	SetParamSet(ctx sdk.Context, ps params.ParamSet)
// }

/*
When a module wishes to interact with another module, it is good practice to define what it will use
as an interface so the module cannot use things that are not permitted.
TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
*/

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) staking.ValidatorI
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator staking.ValidatorI) (stop bool))
	Validator(ctx sdk.Context, address sdk.ValAddress) staking.ValidatorI
	MaxValidators(sdk.Context) uint16
}
