package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/ai-request/exported"
)

// implements AIRequestResultI interface
var _ exported.AIRequestResultI = AIRequestResult{}

// implements ValResultI interface
var _ exported.ValResultI = ValResult{}

// AIRequestResult stores the final result after aggregating the results from the reports of an AI request
type AIRequestResult struct {
	RequestID string     `json:"request_id"`
	Results   ValResults `json:"results"`
	Status    string     `json:"request_status"`
}

// ValResult stores the result information from a validator that has executed the oracle script
type ValResult struct {
	Validator sdk.ValAddress `json:"validator_address"`
	Result    []byte         `json:"result"`
}

// ValResults is the list of results struct
type ValResults []ValResult

// DataSourceResult stores the data source result
type DataSourceResult struct {
	Name   string `json:"data_source"`
	Result []byte `json:"result"`
	Status string `json:"result_status"`
}

// TestCaseResult stores the test case result
type TestCaseResult struct {
	Name              string             `json:"test_case"`
	DataSourceResults []DataSourceResult `json:"data_source_result"`
}
