package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// verify interface at compile time
var _ sdk.Msg = &MsgCreateReport{}

func NewReport(
	requestID string,
	dataSourceResults []*DataSourceResult,
	testCaseResults []*TestCaseResult,
	blockHeight int64,
	fees sdk.Coins,
	aggregatedResult []byte,
	reporter *Reporter,
	status string,
) *Report {
	return &Report{
		RequestID:         requestID,
		Reporter:          reporter,
		DataSourceResults: dataSourceResults,
		TestCaseResults:   testCaseResults,
		BlockHeight:       blockHeight,
		Fees:              fees,
		AggregatedResult:  aggregatedResult,
		ResultStatus:      status,
	}
}

// NewReporter is the constructor of the Reporter struct
func NewReporter(addr sdk.AccAddress, name string, valAddr sdk.ValAddress) *Reporter {
	return &Reporter{
		Address:   addr,
		Name:      name,
		Validator: valAddr,
	}
}

// NewDataSourceResult is the constructor of the data source result struct
func NewDataSourceResult(
	name string,
	result []byte,
	status string,
) *DataSourceResult {
	return &DataSourceResult{
		Name:   name,
		Result: result,
		Status: status,
	}
}

// NewTestCaseResult is the constructor of the test case result struct
func NewTestCaseResult(
	name string,
	dataSourceResults []*DataSourceResult,
) *TestCaseResult {
	return &TestCaseResult{
		Name:              name,
		DataSourceResults: dataSourceResults,
	}
}

// NewMsgCreateReport is a constructor function for MsgCreateReport
func NewMsgCreateReport(
	requestID string,
	dataSourceResults []*DataSourceResult,
	testCaseResults []*TestCaseResult,
	reporter *Reporter,
	fees sdk.Coins,
	aggregatedResult []byte,
	status string,
) *MsgCreateReport {
	return &MsgCreateReport{
		RequestID:         requestID,
		DataSourceResults: dataSourceResults,
		TestCaseResults:   testCaseResults,
		Reporter:          reporter,
		Fees:              fees,
		AggregatedResult:  aggregatedResult,
		ResultStatus:      status,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(reports []Report, reporters []Reporter) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		Reports:   reports,
		Reporters: reporters,
	}
}
