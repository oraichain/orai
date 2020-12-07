package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	provider "github.com/oraichain/orai/x/provider/exported"
)

// AIRequestI is an interface of all the functions related to the AIRequest that are used by other modules
type AIRequestI interface {
	GetRequestID() string
	GetOScriptName() string
	GetCreator() sdk.AccAddress
	GetValidators() []sdk.ValAddress
	GetBlockHeight() int64
	GetAIDataSources() []provider.AIDataSourceI
	GetTestCases() []provider.TestCaseI
	GetFees() sdk.Coins
	GetInput() []byte
	GetExpectedOutput() []byte
}
