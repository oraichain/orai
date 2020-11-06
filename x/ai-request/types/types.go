package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
