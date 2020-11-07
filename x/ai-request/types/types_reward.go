package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/ai-request/exported"
	provider "github.com/oraichain/orai/x/provider/exported"
)

// Reward implements the RewardI interface
var _ exported.RewardI = Reward{}

// Reward stores a list of validators, data source owners and test case owners that receive rewards for a specific block height
type Reward struct {
	Validators    []Validator              `json:"validators"`
	DataSources   []provider.AIDataSourceI `json:"data_sources"`
	TestCases     []provider.TestCaseI     `json:"test_cases"`
	BlockHeight   int64                    `json:"block_height"`
	TotalPower    int64                    `json:"total_voting_power"`
	ProviderFees  sdk.Coins                `json:"provider_fees"`
	ValidatorFees sdk.Coins                `json:"validator_fees"`
}

// NewReward is a constructor for the reward struct
func NewReward(
	validators []Validator,
	dataSources []provider.AIDataSourceI,
	testCases []provider.TestCaseI,
	blockHeight int64,
	totalVotingPower int64,
	providerFees sdk.Coins,
	validatorFees sdk.Coins,
) Reward {
	return Reward{
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
func DefaultReward(blockHeight int64) Reward {
	return Reward{
		Validators:    []Validator{},
		DataSources:   make([]provider.AIDataSourceI, 0),
		TestCases:     make([]provider.TestCaseI, 0),
		BlockHeight:   blockHeight,
		TotalPower:    int64(0),
		ProviderFees:  sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0)))),
		ValidatorFees: sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0)))),
	}
}
