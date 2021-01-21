package types

// import (
// 	"fmt"
// 	"strings"
// )

// // Strategy stores the information of a strategy for a yAI flow
// type Strategy struct {
// 	StratID        uint64   `json:"strategy_id"`
// 	StratName      string   `json:"strategy_name"`
// 	StratFlow      []string `json:"strategy_flow"`
// 	PerformanceFee uint64   `json:"performance_fee"`
// 	PerformanceMax uint64   `json:"performance_max"`
// 	WithdrawalFee  uint64   `json:"withdrawal_fee"`
// 	WithdrawalMax  uint64   `json:"withdrawal_max"`
// 	GovernanceAddr string   `json:"governance_address"`
// 	StrategistAddr string   `json:"strategist_address"`
// }

// // implement fmt.Stringer
// func (s Strategy) String() string {
// 	stratFlows := fmt.Sprintln(s.StratFlow)
// 	return strings.TrimSpace(fmt.Sprintf(`StratID: %d
// 	StratName: %s StratFlow: %s PerformanceFee: %s PerformanceMax %d WithdrawalFee %d WithdrawalMax %d GovernanceAddr %s StrategistAddr %s`, s.StratID, s.StratName, stratFlows, s.PerformanceFee, s.PerformanceMax, s.WithdrawalFee, s.WithdrawalMax, s.GovernanceAddr, s.StrategistAddr))
// }
