package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/exported"
)

// type check for the implementation of the interface ReportI
var _ exported.ReportI = (*Report)(nil)

// Report stores the result of the data source when validator executes it
type Report struct {
	RequestID         string                       `json:"request_id"`
	DataSourceResults []exported.DataSourceResultI `json:"data_source_results"`
	TestCaseResults   []exported.TestCaseResultI   `json:"test_case_results"`
	BlockHeight       int64                        `json:"block_height"`
	Fees              sdk.Coins                    `json:"report_fee"`
	AggregatedResult  []byte                       `json:"aggregated_result"`
	ResultStatus      string                       `json:"result_status"`
	Reporter          Reporter                     `json:"reporter"`
}

// NewReport is the constructor of the report struct
func NewReport(
	requestID string,
	dataSourceResults []exported.DataSourceResultI,
	testCaseResults []exported.TestCaseResultI,
	blockHeight int64,
	fees sdk.Coins,
	aggregatedResult []byte,
	reporter Reporter,
	status string,
) Report {
	return Report{
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

// GetRequestID is the getter method for getting request ID of a report
func (r Report) GetRequestID() string {
	return r.RequestID
}

// GetDataSourceResults is the getter method for getting DataSourceResults of a report
func (r Report) GetDataSourceResults() []exported.DataSourceResultI {
	return r.DataSourceResults
}

// GetTestCaseResults is the getter method for getting TestCaseResults of a report
func (r Report) GetTestCaseResults() []exported.TestCaseResultI {
	return r.TestCaseResults
}

// GetValidator is the getter method for getting Validator of a report
func (r Report) GetValidator() sdk.ValAddress {
	return r.Reporter.Validator
}

// GetReporter is the getter method for getting Reporter of a report
func (r Report) GetReporter() exported.ReporterI {
	return r.Reporter
}

// GetAggregatedResult is the getter method for getting AggregatedResult of a report
func (r Report) GetAggregatedResult() []byte {
	return r.AggregatedResult
}

// GetBlockHeight is the getter method for getting BlockHeight of a report
func (r Report) GetBlockHeight() int64 {
	return r.BlockHeight
}

// GetFees is the getter method for getting Fees of a report
func (r Report) GetFees() sdk.Coins {
	return r.Fees
}

// GetResultStatus is the getter method for getting aggregated result status of a report
func (r Report) GetResultStatus() string {
	return r.ResultStatus
}
