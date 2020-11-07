package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// NewDataSourceResult is the constructor of the data source result struct
func NewDataSourceResult(
	name string,
	result []byte,
	status string,
) DataSourceResult {
	return DataSourceResult{
		Name:   name,
		Result: result,
		Status: status,
	}
}

// NewTestCaseResult is the constructor of the test case result struct
func NewTestCaseResult(
	name string,
	dataSourceResults []DataSourceResult,
) TestCaseResult {
	return TestCaseResult{
		Name:              name,
		DataSourceResults: dataSourceResults,
	}
}

// Validator mimics the original validator to store information of those that execute the oScript
type Validator struct {
	Address     sdk.ValAddress `json:"address"`
	VotingPower int64          `json:"voting_power"`
	Status      string         `json:"status"`
}

// Report stores the result of the data source when validator executes it
type Report struct {
	RequestID         string             `json:"request_id"`
	Validator         Validator          `json:"validator"`
	DataSourceResults []DataSourceResult `json:"data_source_results"`
	TestCaseResults   []TestCaseResult   `json:"test_case_results"`
	BlockHeight       int64              `json:"block_height"`
	Fees              sdk.Coins          `json:"report_fee"`
	AggregatedResult  []byte             `json:"aggregated_result"`
}

// Strategy stores the information of a strategy for a yAI flow
type Strategy struct {
	StratID        uint64   `json:"strategy_id"`
	StratName      string   `json:"strategy_name"`
	StratFlow      []string `json:"strategy_flow"`
	PerformanceFee uint64   `json:"performance_fee"`
	PerformanceMax uint64   `json:"performance_max"`
	WithdrawalFee  uint64   `json:"withdrawal_fee"`
	WithdrawalMax  uint64   `json:"withdrawal_max"`
	GovernanceAddr string   `json:"governance_address"`
	StrategistAddr string   `json:"strategist_address"`
}

// implement fmt.Stringer
func (s Strategy) String() string {
	stratFlows := fmt.Sprintln(s.StratFlow)
	return strings.TrimSpace(fmt.Sprintf(`StratID: %d
	StratName: %s StratFlow: %s PerformanceFee: %s PerformanceMax %d WithdrawalFee %d WithdrawalMax %d GovernanceAddr %s StrategistAddr %s`, s.StratID, s.StratName, stratFlows, s.PerformanceFee, s.PerformanceMax, s.WithdrawalFee, s.WithdrawalMax, s.GovernanceAddr, s.StrategistAddr))
}
