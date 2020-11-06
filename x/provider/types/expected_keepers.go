package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/oraichain/orai/x/provider/exported"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

// /*
// When a module wishes to interact with another module, it is good practice to define what it will use
// as an interface so the module cannot use things that are not permitted.
// TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
// type BankKeeper interface {
// 	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
// 	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
// }
// */

// // BankKeeper is an interface used by other modules such as Keeper. TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
// type BankKeeper interface {
// 	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
// 	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
// 	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
// 	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
// }

// // SupplyKeeper defines the expected supply Keeper.
// type SupplyKeeper interface {
// 	GetModuleAccount(ctx sdk.Context, name string) supply.ModuleAccountI
// 	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
// 	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
// }

// // StakingKeeper defines the expected staking keeper.
// type StakingKeeper interface {
// 	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) staking.ValidatorI
// 	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator staking.ValidatorI) (stop bool))
// 	Validator(ctx sdk.Context, address sdk.ValAddress) staking.ValidatorI
// }

// // DistrKeeper defines the expected distribution keeper.
// type DistrKeeper interface {
// 	GetCommunityTax(ctx sdk.Context) (percent sdk.Dec)
// 	GetBaseProposerReward(ctx sdk.Context) (percent sdk.Dec)
// 	GetBonusProposerReward(ctx sdk.Context) (percent sdk.Dec)
// 	GetFeePool(ctx sdk.Context) (feePool distr.FeePool)
// 	SetFeePool(ctx sdk.Context, feePool distr.FeePool)
// 	AllocateTokensToValidator(ctx sdk.Context, val staking.ValidatorI, tokens sdk.DecCoins)
// }

// OracleScriptSet is an interface for all the related properties for interacting with the Oracle Script struct for Keeper
type OracleScriptSet interface {
	// GetOracleScript(ctx sdk.Context, name string) (exported.OracleScriptI, error)
	// GetOracleScripts(ctx sdk.Context, page, limit uint) ([]exported.OracleScriptI, error)
	// GetAllOracleScriptNames(ctx sdk.Context) sdk.Iterator
	// GetPaginatedOracleScriptNames(ctx sdk.Context, page, limit uint) sdk.Iterator
	// EditOracleScript(ctx sdk.Context, oldName string, newName string, oScript exported.OracleScriptI)
	// SetOracleScript(ctx sdk.Context, name string, oScript exported.OracleScriptI)
	// AddOracleScriptFile(file []byte, name string)
	// EraseOracleScriptFile(name string)
	// EditOracleScriptFile(file []byte, name string)
	// GetOracleScriptFile(name string) []byte
}

// AIDataSourceSet is an interface for all the related properties for interacting with the Data Source struct for Keeper that are exported for other modules to use
type AIDataSourceSet interface {
	GetAIDataSourceI(ctx sdk.Context, name string) (exported.AIDataSourceI, error)
	DefaultAIDataSourceI() exported.AIDataSourceI
	// GetAIDataSources(ctx sdk.Context, page, limit uint) ([]exported.AIDataSourceI, error)
	// GetAllAIDataSourceNames(ctx sdk.Context) sdk.Iterator
	// GetPaginatedAIDataSourceNames(ctx sdk.Context, page, limit uint) sdk.Iterator
	// EditAIDataSource(ctx sdk.Context, oldName string, newName string, oScript exported.AIDataSourceI)
	// SetAIDataSource(ctx sdk.Context, name string, oScript exported.AIDataSourceI)
	// AddAIDataSourceFile(file []byte, name string)
	// EraseAIDataSourceFile(name string)
	// EditAIDataSourceFile(file []byte, name string)
	// GetAIDataSourceFile(name string) []byte
}

// TestCaseSet is an interface for all the related properties for interacting with the Test Case struct for Keeper that are exported for other modules to use
type TestCaseSet interface {
	GetTestCaseI(ctx sdk.Context, name string) (exported.TestCaseI, error)
	DefaultTestCaseI() exported.TestCaseI
	// GetTestCases(ctx sdk.Context, page, limit uint) ([]exported.TestCaseI, error)
	// GetAllTestCaseNames(ctx sdk.Context) sdk.Iterator
	// GetPaginatedTestCaseNames(ctx sdk.Context, page, limit uint) sdk.Iterator
	// EditTestCase(ctx sdk.Context, oldName string, newName string, oScript exported.TestCaseI)
	// SetTestCase(ctx sdk.Context, name string, oScript exported.TestCaseI)
	// AddTestCaseFile(file []byte, name string)
	// EraseTestCaseFile(name string)
	// EditTestCaseFile(file []byte, name string)
	// GetTestCaseFile(name string) []byte
}
