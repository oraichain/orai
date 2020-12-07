package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airesult/exported"
	provider "github.com/oraichain/orai/x/provider/exported"
	webSocket "github.com/oraichain/orai/x/websocket/exported"
)

// Reward implements the RewardI interface
var _ exported.RewardI = Reward{}

// Reward stores a list of validators, data source owners and test case owners that receive rewards for a specific block height
type Reward struct {
	Validators    []webSocket.ValidatorI   `json:"validators"`
	DataSources   []provider.AIDataSourceI `json:"data_sources"`
	TestCases     []provider.TestCaseI     `json:"test_cases"`
	BlockHeight   int64                    `json:"block_height"`
	TotalPower    int64                    `json:"total_voting_power"`
	ProviderFees  sdk.Coins                `json:"provider_fees"`
	ValidatorFees sdk.Coins                `json:"validator_fees"`
}

// NewReward is a constructor for the reward struct
func NewReward(
	validators []webSocket.ValidatorI,
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
		Validators:    make([]webSocket.ValidatorI, 0),
		DataSources:   make([]provider.AIDataSourceI, 0),
		TestCases:     make([]provider.TestCaseI, 0),
		BlockHeight:   blockHeight,
		TotalPower:    int64(0),
		ProviderFees:  sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0)))),
		ValidatorFees: sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(int64(0)))),
	}
}

// SetValidators setter
func (re Reward) SetValidators(vals []webSocket.ValidatorI) error {
	if len(vals) == 0 {
		return errors.New("Cannot set reward validators when empty")
	}
	re.Validators = vals
	return nil
}

// GetValidators getter
func (re Reward) GetValidators() []webSocket.ValidatorI {
	return re.Validators
}

// SetDataSources setter
func (re Reward) SetDataSources(dSources []provider.AIDataSourceI) error {
	if len(dSources) == 0 {
		return errors.New("Cannot set AI data sources when empty")
	}
	re.DataSources = dSources
	return nil
}

// GetDataSources getter
func (re Reward) GetDataSources() []provider.AIDataSourceI {
	return re.DataSources
}

// SetTestCases setter
func (re Reward) SetTestCases(tCases []provider.TestCaseI) error {
	return nil
}

// GetTestCases getter
func (re Reward) GetTestCases() []provider.TestCaseI {
	return re.TestCases
}

// SetBlockHeight setter
func (re Reward) SetBlockHeight(blockHeight int64) error {
	return nil
}

// GetBlockHeight getter
func (re Reward) GetBlockHeight() int64 {
	return re.BlockHeight
}

// SetTotalPower setter
func (re Reward) SetTotalPower(power int64) error {
	return nil
}

// GetTotalPower getter
func (re Reward) GetTotalPower() int64 {
	return re.TotalPower
}

// SetProviderFees setter
func (re Reward) SetProviderFees(fees sdk.Coins) error {
	return nil
}

// GetProviderFees getter
func (re Reward) GetProviderFees() sdk.Coins {
	return re.ProviderFees
}

// SetValidatorFees setter
func (re Reward) SetValidatorFees(fees sdk.Coins) error {
	return nil
}

// GetValidatorFees getter
func (re Reward) GetValidatorFees() sdk.Coins {
	return re.ValidatorFees
}
