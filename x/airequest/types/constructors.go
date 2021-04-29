package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAIRequest is the constructor of the ai request struct
func NewAIRequest(
	requestID string,
	contract sdk.AccAddress,
	creator sdk.AccAddress,
	validators []sdk.ValAddress,
	blockHeight int64,
	fees sdk.Coins,
	input []byte,
	testOnly bool,
) *AIRequest {
	return &AIRequest{
		RequestID:   requestID,
		Contract:    contract,
		Creator:     creator,
		Validators:  validators,
		BlockHeight: blockHeight,
		Fees:        fees,
		Input:       input,
		TestOnly:    testOnly,
	}
}

// NewMsgSetAIRequestReq is a constructor function for NewMsgSetAIRequest
func NewMsgSetAIRequestReq(requestID string, contract, creator, fees string, valCount int64, input []byte, testOnly bool) *MsgSetAIRequestReq {
	return &MsgSetAIRequestReq{
		RequestID:      requestID,
		Contract:       contract,
		Creator:        creator,
		ValidatorCount: valCount,
		Fees:           fees,
		Input:          input,
		TestOnly:       testOnly,
	}
}

// NewMsgSetAIRequestRes is a constructor function for NewMsgSetAIRequestRes
func NewMsgSetAIRequestRes(requestID string, contract, creator, fees string, valCount int64, input []byte, testOnly bool) *MsgSetAIRequestRes {
	return &MsgSetAIRequestRes{
		RequestID:      requestID,
		Contract:       contract,
		Creator:        creator,
		ValidatorCount: valCount,
		Fees:           fees,
		Input:          input,
		TestOnly:       testOnly,
	}
}

// NewQueryAIRequestRes is the constructor for the QueryAIRequestRes
func NewQueryAIRequestRes(requestID string, contract sdk.AccAddress, creator sdk.AccAddress, fees sdk.Coins, validators []sdk.ValAddress, blockHeight int64, input []byte) *QueryAIRequestRes {
	return &QueryAIRequestRes{
		RequestId:   requestID,
		Contract:    contract,
		Creator:     creator,
		Validators:  validators,
		BlockHeight: blockHeight,
		Fees:        fees,
		Input:       input,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(AIRequests []AIRequest, params Params) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		AIRequests: AIRequests,
		Params:     params,
	}
}

// NewAIRequestResult is a constructor for the ai request result struct
func NewAIRequestResult(
	requestID string,
	results []ValResult,
	status string,
) *AIRequestResult {
	return &AIRequestResult{
		RequestID: requestID,
		Results:   results,
		Status:    status,
	}
}

func NewQueryFullRequestRes(
	AIRequest AIRequest,
	reports []Report,
	tcReports []TestCaseReport,
	result AIRequestResult,
) *QueryFullOracleRes {
	return &QueryFullOracleRes{
		AIRequest:       AIRequest,
		Reports:         reports,
		TestcaseReports: tcReports,
		Result:          result,
	}
}

func NewQueryRewardRes(
	reward Reward,
) *QueryRewardRes {
	return &QueryRewardRes{
		Reward: reward,
	}
}

// NewReward is a constructor for the reward struct
func NewReward(
	validators []Validator,
	blockHeight int64,
	totalVotingPower int64,
	providerFees sdk.Coins,
	validatorFees sdk.Coins,
	results []*Result,
) *Reward {
	return &Reward{
		BaseReward: &BaseReward{
			Validators:    validators,
			BlockHeight:   blockHeight,
			TotalPower:    totalVotingPower,
			ProviderFees:  providerFees,
			ValidatorFees: validatorFees,
		},
		Results: results,
	}
}

// DefaultReward is a default value init for the reward struct
func DefaultReward(blockHeight int64) *Reward {
	return &Reward{
		BaseReward: &BaseReward{
			Validators:    make([]Validator, 0),
			BlockHeight:   blockHeight,
			TotalPower:    int64(0),
			ProviderFees:  sdk.NewCoins(sdk.NewCoin(Denom, sdk.NewInt(int64(0)))),
			ValidatorFees: sdk.NewCoins(sdk.NewCoin(Denom, sdk.NewInt(int64(0)))),
		},
		Results: []*Result{},
	}
}

func NewReport(
	requestID string,
	dataSourceResults []*Result,
	blockHeight int64,
	aggregatedResult []byte,
	valAddress sdk.ValAddress,
	status string,
	fees sdk.Coins,
) *Report {
	return &Report{
		BaseReport: &BaseReport{
			RequestId:        requestID,
			ValidatorAddress: valAddress,
			BlockHeight:      blockHeight,
			Fees:             fees,
		},
		DataSourceResults: dataSourceResults,
		AggregatedResult:  aggregatedResult,
		ResultStatus:      status,
	}
}

func NewTestCaseReport(
	requestID string,
	results []*ResultWithTestCase,
	blockHeight int64,
	valAddress sdk.ValAddress,
	fees sdk.Coins,
) *TestCaseReport {
	return &TestCaseReport{
		BaseReport: &BaseReport{
			RequestId:        requestID,
			ValidatorAddress: valAddress,
			BlockHeight:      blockHeight,
			Fees:             fees,
		},
		ResultsWithTestCase: results,
	}
}

// NewResult is the constructor of the result struct
func NewResult(
	entryPoint *EntryPoint,
	result []byte,
	status string,
) *Result {
	return &Result{
		EntryPoint: entryPoint,
		Result:     result,
		Status:     status,
	}
}

// NewResultWithTestCase is the constructor of the result with test case struct
func NewResultWithTestCase(
	entryPoint *EntryPoint,
	result []*Result,
	status string,
) *ResultWithTestCase {
	return &ResultWithTestCase{
		EntryPoint:      entryPoint,
		TestCaseResults: result,
		Status:          status,
	}
}

// // NewTestCaseResult is the constructor of the test case result struct
// func NewTestCaseResult(
// 	entryPoint *EntryPoint,
// 	dataSourceResults []*Result,
// 	status string,
// ) *TestCaseResult {
// 	return &TestCaseResult{
// 		EntryPoint:        entryPoint,
// 		DataSourceResults: dataSourceResults,
// 		Status:            status,
// 	}
// }
