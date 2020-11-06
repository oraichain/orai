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
