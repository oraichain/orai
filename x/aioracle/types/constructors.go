package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAIOracle is the constructor of the ai request struct
func NewAIOracle(
	requestID string,
	contract sdk.AccAddress,
	creator sdk.AccAddress,
	validators []sdk.ValAddress,
	blockHeight int64,
	fees sdk.Coins,
	input []byte,
) *AIOracle {
	return &AIOracle{
		RequestID:   requestID,
		Contract:    contract,
		Creator:     creator,
		Validators:  validators,
		BlockHeight: blockHeight,
		Fees:        fees,
		Input:       input,
	}
}

// NewMsgSetAIOracleReq is a constructor function for NewMsgSetAIOracle
func NewMsgSetAIOracleReq(requestID string, contract sdk.AccAddress, creator sdk.AccAddress, fees string, valCount int64, input []byte) *MsgSetAIOracleReq {
	return &MsgSetAIOracleReq{
		RequestID:      requestID,
		Contract:       contract,
		Creator:        creator,
		ValidatorCount: valCount,
		Fees:           fees,
		Input:          input,
	}
}

// NewMsgSetAIOracleRes is a constructor function for NewMsgSetAIOracleRes
func NewMsgSetAIOracleRes(requestID string, contract sdk.AccAddress, creator sdk.AccAddress, fees string, valCount int64, input []byte) *MsgSetAIOracleRes {
	return &MsgSetAIOracleRes{
		RequestID:      requestID,
		Contract:       contract,
		Creator:        creator,
		ValidatorCount: valCount,
		Fees:           fees,
		Input:          input,
	}
}

// NewQueryAIOracleRes is the constructor for the QueryAIOracleRes
func NewQueryAIOracleRes(requestID string, contract sdk.AccAddress, creator sdk.AccAddress, fees sdk.Coins, validators []sdk.ValAddress, blockHeight int64, input []byte) *QueryAIOracleRes {
	return &QueryAIOracleRes{
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
func NewGenesisState(AIOracles []AIOracle, params Params) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		AIOracles: AIOracles,
		Params:    params,
	}
}

// NewAIOracleResult is a constructor for the ai request result struct
func NewAIOracleResult(
	requestID string,
	results []ValResult,
	status string,
) *AIOracleResult {
	return &AIOracleResult{
		RequestID: requestID,
		Results:   results,
		Status:    status,
	}
}

func NewQueryFullRequestRes(
	AIOracle AIOracle,
	reports []Report,
	result AIOracleResult,
) *QueryFullOracleRes {
	return &QueryFullOracleRes{
		AIOracle: AIOracle,
		Reports:  reports,
		Result:   result,
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
	dataSourceResults []*DataSourceResult,
) *Reward {
	return &Reward{
		Validators:        validators,
		BlockHeight:       blockHeight,
		TotalPower:        totalVotingPower,
		ProviderFees:      providerFees,
		ValidatorFees:     validatorFees,
		DatasourceResults: dataSourceResults,
	}
}

// DefaultReward is a default value init for the reward struct
func DefaultReward(blockHeight int64) *Reward {
	return &Reward{
		Validators:        make([]Validator, 0),
		BlockHeight:       blockHeight,
		TotalPower:        int64(0),
		ProviderFees:      sdk.NewCoins(sdk.NewCoin(Denom, sdk.NewInt(int64(0)))),
		ValidatorFees:     sdk.NewCoins(sdk.NewCoin(Denom, sdk.NewInt(int64(0)))),
		DatasourceResults: []*DataSourceResult{},
	}
}

func NewReport(
	requestID string,
	dataSourceResults []*DataSourceResult,
	blockHeight int64,
	aggregatedResult []byte,
	valAddress sdk.ValAddress,
	status string,
) *Report {
	return &Report{
		RequestID:         requestID,
		ValidatorAddress:  valAddress,
		DataSourceResults: dataSourceResults,
		BlockHeight:       blockHeight,
		AggregatedResult:  aggregatedResult,
		ResultStatus:      status,
	}
}

// NewDataSourceResult is the constructor of the data source result struct
func NewDataSourceResult(
	entryPoint *EntryPoint,
	result []byte,
	status string,
) *DataSourceResult {
	return &DataSourceResult{
		EntryPoint: entryPoint,
		Result:     result,
		Status:     status,
	}
}

// NewTestCaseResult is the constructor of the test case result struct
func NewTestCaseResult(
	entryPoint *EntryPoint,
	dataSourceResults []*DataSourceResult,
	status string,
) *TestCaseResult {
	return &TestCaseResult{
		EntryPoint:        entryPoint,
		DataSourceResults: dataSourceResults,
		Status:            status,
	}
}
