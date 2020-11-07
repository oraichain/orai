package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewReport is the constructor of the report struct
func NewReport(
	requestID string,
	validator Validator,
	dataSourceResults []DataSourceResult,
	testCaseResults []TestCaseResult,
	blockHeight int64,
	fees sdk.Coins,
	aggregatedResult []byte,
) Report {
	return Report{
		RequestID:         requestID,
		Validator:         validator,
		DataSourceResults: dataSourceResults,
		TestCaseResults:   testCaseResults,
		BlockHeight:       blockHeight,
		Fees:              fees,
		AggregatedResult:  aggregatedResult,
	}
}

// NewValidator is the constructor of the validator struct
func NewValidator(
	address sdk.ValAddress,
	votingPower int64,
	status string,
) Validator {
	return Validator{
		Address:     address,
		VotingPower: votingPower,
		Status:      status,
	}
}

// NewStrategy is a constructor for the strategy struct
func NewStrategy(
	stratID uint64,
	stratName string,
	stratFlow []string,
	performanceFee uint64,
	performanceMax uint64,
	withdrawalFee uint64,
	withdrawalMax uint64,
	governanceAddr string,
	strategistAddr string,
) Strategy {
	return Strategy{
		StratID:        stratID,
		StratName:      stratName,
		StratFlow:      stratFlow,
		PerformanceFee: performanceFee,
		PerformanceMax: performanceMax,
		WithdrawalFee:  withdrawalFee,
		WithdrawalMax:  withdrawalMax,
		GovernanceAddr: governanceAddr,
		StrategistAddr: strategistAddr,
	}
}
