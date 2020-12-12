package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	aiRequest "github.com/oraichain/orai/x/airequest/exported"
	"github.com/oraichain/orai/x/websocket/exported"
)

// ReportSet is an interface for all the related properties for interacting with the Report struct for Keeper that are exported for other modules to use
type ReportSet interface {
	GetReports(ctx sdk.Context, rid string) (reports []exported.ReportI)
	HasReport(ctx sdk.Context, id string, val sdk.ValAddress) bool
	GetReportsBlockHeight(ctx sdk.Context, blockHeight int64) (reports []exported.ReportI)
	ValidateReport(ctx sdk.Context, rep exported.ReportI, req aiRequest.AIRequestI) error
}

// ReporterSet is an interface for all the related properties for interacting with the Reporter struct for Keeper that are exported for other modules to use
type ReporterSet interface {
}

// DataSourceResultSet is an interface for all the related properties for interacting with the DataSourceResult struct for Keeper that are exported for other modules to use
type DataSourceResultSet interface {
}

// TestCaseResultSet is an interface for all the related properties for interacting with the TestCase struct for Keeper that are exported for other modules to use
type TestCaseResultSet interface {
}

// ValidatorSet is an interface for all the related properties for interacting with the Validator struct for Keeper that are exported for other modules to use
type ValidatorSet interface {
	NewValidator(address sdk.ValAddress, votingPower int64, status string) exported.ValidatorI
}

// ValResultSet is an interface for all the related properties for interacting with the ValResult struct for Keeper that are exported for other modules to use
type ValResultSet interface {
	DefaultValResultI() exported.ValResultI
	GetKeyResultSuccess() string
	NewValResult(val exported.ValidatorI, result []byte, status string) exported.ValResultI
}
