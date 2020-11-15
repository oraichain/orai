package types

import "github.com/oraichain/orai/x/websocket/exported"

// type check for the implementation of the interface TestCaseResultI
var _ exported.TestCaseResultI = (*TestCaseResult)(nil)

// TestCaseResult stores the test case result
type TestCaseResult struct {
	Name              string                       `json:"test_case"`
	DataSourceResults []exported.DataSourceResultI `json:"data_source_result"`
}

// NewTestCaseResult is the constructor of the test case result struct
func NewTestCaseResult(
	name string,
	dataSourceResults []exported.DataSourceResultI,
) TestCaseResult {
	return TestCaseResult{
		Name:              name,
		DataSourceResults: dataSourceResults,
	}
}

// GetName is the getter method for getting Name of a TestCaseResult
func (tc TestCaseResult) GetName() string {
	return tc.Name
}

// GetDataSourceResults is the getter method for getting DataSourceResults of a TestCaseResult
func (tc TestCaseResult) GetDataSourceResults() []exported.DataSourceResultI {
	return tc.DataSourceResults
}
