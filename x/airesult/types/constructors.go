package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	airequest "github.com/oraichain/orai/x/airequest/types"
	provider "github.com/oraichain/orai/x/provider/types"
	websocket "github.com/oraichain/orai/x/websocket/types"
)

// NewAIRequestResult is a constructor for the ai request result struct
func NewAIRequestResult(
	requestID string,
	results []websocket.ValResult,
	status string,
) *AIRequestResult {
	return &AIRequestResult{
		RequestID: requestID,
		Results:   results,
		Status:    status,
	}
}

func NewQueryFullRequestRes(
	aiRequest airequest.AIRequest,
	reports []websocket.Report,
	result AIRequestResult,
) *QueryFullRequestRes {
	return &QueryFullRequestRes{
		AIRequest: aiRequest,
		Reports:   reports,
		Result:    result,
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
	validators []websocket.Validator,
	dataSources []provider.AIDataSource,
	testCases []provider.TestCase,
	blockHeight int64,
	totalVotingPower int64,
	providerFees sdk.Coins,
	validatorFees sdk.Coins,
) *Reward {
	return &Reward{
		Validators:    validators,
		DataSources:   dataSources,
		TestCases:     testCases,
		BlockHeight:   blockHeight,
		TotalPower:    totalVotingPower,
		ProviderFees:  providerFees,
		ValidatorFees: validatorFees,
	}
}

// DefaultReward is a default value init for the reward struct
func DefaultReward(blockHeight int64) *Reward {
	return &Reward{
		Validators:    make([]websocket.Validator, 0),
		DataSources:   make([]provider.AIDataSource, 0),
		TestCases:     make([]provider.TestCase, 0),
		BlockHeight:   blockHeight,
		TotalPower:    int64(0),
		ProviderFees:  sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0)))),
		ValidatorFees: sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0)))),
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(aiResults []AIRequestResult, params Params) *GenesisState {
	return &GenesisState{
		// TODO: Fill out according to your genesis state
		AIRequestResults: aiResults,
		Params:           params,
	}
}
