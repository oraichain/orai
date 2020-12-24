package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataSourceResultI is the generic DataSourceResult interface that is used for other modules
type DataSourceResultI interface {
	GetName() string
	GetResult() []byte
	GetStatus() string
}

// TestCaseResultI is the generic TestCaseResult interface that is used for other modules
type TestCaseResultI interface {
	GetName() string
	GetDataSourceResults() []DataSourceResultI
}

// ReportI is the generic Report interface that is used for other modules
type ReportI interface {
	GetRequestID() string
	GetDataSourceResults() []DataSourceResultI
	GetTestCaseResults() []TestCaseResultI
	GetValidator() sdk.ValAddress
	GetReporter() ReporterI
	GetAggregatedResult() []byte
	GetBlockHeight() int64
	GetFees() sdk.Coins
	GetResultStatus() string
}

// ReporterI is the generic Reporter interface that is used for other modules
type ReporterI interface {
	GetAddress() sdk.AccAddress
	GetName() string
	GetValidator() sdk.ValAddress
}

// ValidatorI is the generic Validator interface that is used for other modules
type ValidatorI interface {
	GetAddress() sdk.ValAddress
	GetVotingPower() int64
	GetStatus() string
}

// ValResultI is the generic ValidatorResult interface that is used for other modules
type ValResultI interface {
	GetValidator() ValidatorI
	GetResult() []byte
	GetResultStatus() string
}
