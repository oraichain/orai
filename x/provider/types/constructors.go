package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewOracleScript is the constructor of the oScript struct
func NewOracleScript(
	name string,
	owner sdk.AccAddress,
	des string,
	minimumFees sdk.Coins,
) OracleScript {
	return OracleScript{
		Name:        name,
		Owner:       owner,
		Description: des,
		MinimumFees: minimumFees,
	}
}

// NewAIDataSource is the constructor of the data source struct
func NewAIDataSource(
	name string,
	owner sdk.AccAddress,
	fees sdk.Coins,
	des string,
) AIDataSource {
	return AIDataSource{
		Name:        name,
		Owner:       owner,
		Fees:        fees,
		Description: des,
	}
}

// NewTestCase is the constructor of the test case struct
func NewTestCase(
	name string,
	owner sdk.AccAddress,
	fees sdk.Coins,
	des string,
) TestCase {
	return TestCase{
		Name:        name,
		Owner:       owner,
		Fees:        fees,
		Description: des,
	}
}

// NewAIRequest is the constructor of the ai request struct
func NewAIRequest(
	requestID string,
	oscriptName string,
	creator sdk.AccAddress,
	validators []sdk.ValAddress,
	blockHeight int64,
	aiDataSources []AIDataSource,
	testCases []TestCase,
	fees sdk.Coins,
	input string,
	expectedOutput string,
) AIRequest {
	return AIRequest{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		Validators:       validators,
		BlockHeight:      blockHeight,
		AIDataSources:    aiDataSources,
		TestCases:        testCases,
		Fees:             fees,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}

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

// NewReward is a constructor for the reward struct
func NewReward(
	validators []Validator,
	dataSourceOwners []sdk.AccAddress,
	testCaseOwners []sdk.AccAddress,
	blockHeight int64,
	totalVotingPower int64,
) Reward {
	return Reward{
		Validators:       validators,
		DataSourceOwners: dataSourceOwners,
		TestCaseOwners:   testCaseOwners,
		BlockHeight:      blockHeight,
		TotalPower:       totalVotingPower,
	}
}

// NewAIRequestResult is a constructor for the ai request result struct
func NewAIRequestResult(
	requestID string,
	results [][]byte,
	status string,
) AIRequestResult {
	return AIRequestResult{
		RequestID: requestID,
		Results:   results,
		Status:    status,
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
