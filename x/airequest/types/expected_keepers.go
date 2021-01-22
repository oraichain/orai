package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	provider "github.com/oraichain/orai/x/provider/types"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

/*
When a module wishes to interact with another module, it is good practice to define what it will use
as an interface so the module cannot use things that are not permitted.
TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
*/

// BankKeeper is an interface used by other modules such as Keeper. TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// AuthKeeper defines the expected auth Keeper.
type AuthKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) auth.ModuleAccountI
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) staking.ValidatorI
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator staking.ValidatorI) (stop bool))
	Validator(ctx sdk.Context, address sdk.ValAddress) staking.ValidatorI
	MaxValidators(sdk.Context) uint32
}

// DistrKeeper defines the expected distribution keeper.
type DistrKeeper interface {
	GetCommunityTax(ctx sdk.Context) (percent sdk.Dec)
	GetBaseProposerReward(ctx sdk.Context) (percent sdk.Dec)
	GetBonusProposerReward(ctx sdk.Context) (percent sdk.Dec)
	GetFeePool(ctx sdk.Context) (feePool distr.FeePool)
	SetFeePool(ctx sdk.Context, feePool distr.FeePool)
	AllocateTokensToValidator(ctx sdk.Context, val staking.ValidatorI, tokens sdk.DecCoins)
}

// ProviderKeeper defines the expected provider keeper
type ProviderKeeper interface {
	GetAIDataSource(ctx sdk.Context, name string) (*provider.AIDataSource, error)
	GetOracleScript(ctx sdk.Context, name string) (*provider.OracleScript, error)
	GetTestCase(ctx sdk.Context, name string) (*provider.TestCase, error)
	GetDNamesTcNames(ctx sdk.Context, oScript string) ([]string, []string, error)
	GetMinimumFees(ctx sdk.Context, dNames, tcNames []string, valNum int) (sdk.Coins, error)
	GetKeyOracleScriptRewardPercentage(ctx sdk.Context) int64
}
