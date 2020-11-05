package types

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// // NewAIRequest is the constructor of the ai request struct
// func NewAIRequest(
// 	requestID string,
// 	oscriptName string,
// 	creator sdk.AccAddress,
// 	validators []sdk.ValAddress,
// 	blockHeight int64,
// 	aiDataSources []AIDataSource,
// 	testCases []TestCase,
// 	fees sdk.Coins,
// 	input string,
// 	expectedOutput string,
// ) AIRequest {
// 	return AIRequest{
// 		RequestID:        requestID,
// 		OracleScriptName: oscriptName,
// 		Creator:          creator,
// 		Validators:       validators,
// 		BlockHeight:      blockHeight,
// 		AIDataSources:    aiDataSources,
// 		TestCases:        testCases,
// 		Fees:             fees,
// 		Input:            input,
// 		ExpectedOutput:   expectedOutput,
// 	}
// }

// // NewReport is the constructor of the report struct
// func NewReport(
// 	requestID string,
// 	validator Validator,
// 	dataSourceResults []DataSourceResult,
// 	testCaseResults []TestCaseResult,
// 	blockHeight int64,
// 	fees sdk.Coins,
// 	aggregatedResult []byte,
// ) Report {
// 	return Report{
// 		RequestID:         requestID,
// 		Validator:         validator,
// 		DataSourceResults: dataSourceResults,
// 		TestCaseResults:   testCaseResults,
// 		BlockHeight:       blockHeight,
// 		Fees:              fees,
// 		AggregatedResult:  aggregatedResult,
// 	}
// }

// // NewDataSourceResult is the constructor of the data source result struct
// func NewDataSourceResult(
// 	name string,
// 	result []byte,
// 	status string,
// ) DataSourceResult {
// 	return DataSourceResult{
// 		Name:   name,
// 		Result: result,
// 		Status: status,
// 	}
// }

// // NewTestCaseResult is the constructor of the test case result struct
// func NewTestCaseResult(
// 	name string,
// 	dataSourceResults []DataSourceResult,
// ) TestCaseResult {
// 	return TestCaseResult{
// 		Name:              name,
// 		DataSourceResults: dataSourceResults,
// 	}
// }

// // NewValidator is the constructor of the validator struct
// func NewValidator(
// 	address sdk.ValAddress,
// 	votingPower int64,
// 	status string,
// ) Validator {
// 	return Validator{
// 		Address:     address,
// 		VotingPower: votingPower,
// 		Status:      status,
// 	}
// }

// // NewReward is a constructor for the reward struct
// func NewReward(
// 	validators []Validator,
// 	dataSources []AIDataSource,
// 	testCases []TestCase,
// 	blockHeight int64,
// 	totalVotingPower int64,
// 	providerFees sdk.Coins,
// 	validatorFees sdk.Coins,
// ) Reward {
// 	return Reward{
// 		Validators:    validators,
// 		DataSources:   dataSources,
// 		TestCases:     testCases,
// 		BlockHeight:   blockHeight,
// 		TotalPower:    totalVotingPower,
// 		ProviderFees:  providerFees,
// 		ValidatorFees: validatorFees,
// 	}
// }

// // NewAIRequestResult is a constructor for the ai request result struct
// func NewAIRequestResult(
// 	requestID string,
// 	results ValResults,
// 	status string,
// ) AIRequestResult {
// 	return AIRequestResult{
// 		RequestID: requestID,
// 		Results:   results,
// 		Status:    status,
// 	}
// }

// // NewValResult is a constructor for the validator result
// func NewValResult(
// 	val sdk.ValAddress,
// 	result []byte,
// ) ValResult {
// 	return ValResult{
// 		Validator: val,
// 		Result:    result,
// 	}
// }

// // NewStrategy is a constructor for the strategy struct
// func NewStrategy(
// 	stratID uint64,
// 	stratName string,
// 	stratFlow []string,
// 	performanceFee uint64,
// 	performanceMax uint64,
// 	withdrawalFee uint64,
// 	withdrawalMax uint64,
// 	governanceAddr string,
// 	strategistAddr string,
// ) Strategy {
// 	return Strategy{
// 		StratID:        stratID,
// 		StratName:      stratName,
// 		StratFlow:      stratFlow,
// 		PerformanceFee: performanceFee,
// 		PerformanceMax: performanceMax,
// 		WithdrawalFee:  withdrawalFee,
// 		WithdrawalMax:  withdrawalMax,
// 		GovernanceAddr: governanceAddr,
// 		StrategistAddr: strategistAddr,
// 	}
// }
