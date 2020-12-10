package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/exported"
)

// OracleScriptSet is an interface for all the related properties for interacting with the Oracle Script struct for Keeper
type OracleScriptSet interface {
	GetDNamesTcNames(ctx sdk.Context, oScript string) ([]string, []string, error)
	GetOScriptPath(oScriptName string) string
	GetOracleScriptI(ctx sdk.Context, name string) (exported.OracleScriptI, error)
	GetKeyOracleScriptRewardPercentage(ctx sdk.Context) int64
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
