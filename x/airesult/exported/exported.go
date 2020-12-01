package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	provider "github.com/oraichain/orai/x/provider/exported"
	webSocket "github.com/oraichain/orai/x/websocket/exported"
)

type RewardI interface {
	SetValidators(vals []webSocket.ValidatorI) error
	GetValidators() []webSocket.ValidatorI
	SetDataSources(dSources []provider.AIDataSourceI) error
	GetDataSources() []provider.AIDataSourceI
	SetTestCases(tCases []provider.TestCaseI) error
	GetTestCases() []provider.TestCaseI
	SetBlockHeight(blockHeight int64) error
	GetBlockHeight() int64
	SetTotalPower(power int64) error
	GetTotalPower() int64
	SetProviderFees(fees sdk.Coins) error
	GetProviderFees() sdk.Coins
	SetValidatorFees(fees sdk.Coins) error
	GetValidatorFees() sdk.Coins
}

type AIRequestResultI interface {
	GetRequestID() string
	GetValResults() []webSocket.ValResultI
	GetStatus() string
}
