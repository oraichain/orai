package types

import (
	"fmt"
	"strings"

	"github.com/oraichain/orai/x/airequest/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"

	provider "github.com/oraichain/orai/x/provider/exported"
)

// implements AIRequestI interface
var _ exported.AIRequestI = AIRequest{}

// AIRequest stores the request metadata to an AI model
type AIRequest struct {
	RequestID        string                   `json:"request_id"`
	OracleScriptName string                   `json:"oscript_name"`
	Creator          sdk.AccAddress           `json:"request_creator"`
	Validators       []sdk.ValAddress         `json:"validator_addr"`
	BlockHeight      int64                    `json:"block_height"`
	AIDataSources    []provider.AIDataSourceI `json:"data_sources"`
	TestCases        []provider.TestCaseI     `json:"testcases"`
	Fees             sdk.Coins                `json:"transaction_fee"`
	Input            []byte                   `json:"request_input"`
	ExpectedOutput   []byte                   `json:"expected_output"`
	//Time         string `json:"time"`
}

// implement fmt.Stringer
func (ai AIRequest) String() string {
	valString := fmt.Sprintln(ai.Validators)
	dataSourceString := fmt.Sprintln(ai.AIDataSources)
	testCaseString := fmt.Sprintln(ai.TestCases)
	return strings.TrimSpace(fmt.Sprintf(`RequestID: %s
	OracleScript name: %s Creator: %s Validators: %s BlockHeight: %d AIDataSources: %s TestCases: %s Fees: %s`, ai.RequestID, ai.OracleScriptName, string(ai.Creator[:]), valString, ai.BlockHeight, dataSourceString, testCaseString, ai.Fees.String()))
}

// NewAIRequest is the constructor of the ai request struct
func NewAIRequest(
	requestID string,
	oscriptName string,
	creator sdk.AccAddress,
	validators []sdk.ValAddress,
	blockHeight int64,
	aiDataSources []provider.AIDataSourceI,
	testCases []provider.TestCaseI,
	fees sdk.Coins,
	input []byte,
	expectedOutput []byte,
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

// GetRequestID is getter method for AIRequest struct
func (ai AIRequest) GetRequestID() string {
	return ai.RequestID
}

// GetOScriptName is getter method for AIRequest struct
func (ai AIRequest) GetOScriptName() string {
	return ai.OracleScriptName
}

// GetCreator is getter method for AIRequest struct
func (ai AIRequest) GetCreator() sdk.AccAddress {
	return ai.Creator
}

// GetValidators is getter method for AIRequest struct
func (ai AIRequest) GetValidators() []sdk.ValAddress {
	return ai.Validators
}

// GetBlockHeight is getter method for AIRequest struct
func (ai AIRequest) GetBlockHeight() int64 {
	return ai.BlockHeight
}

// GetAIDataSources is getter method for AIRequest struct
func (ai AIRequest) GetAIDataSources() []provider.AIDataSourceI {
	return ai.AIDataSources
}

// GetTestCases is getter method for AIRequest struct
func (ai AIRequest) GetTestCases() []provider.TestCaseI {
	return ai.TestCases
}

// GetFees is getter method for AIRequest struct
func (ai AIRequest) GetFees() sdk.Coins {
	return ai.Fees
}

// GetInput is getter method for AIRequest struct
func (ai AIRequest) GetInput() []byte {
	return ai.Input
}

// GetExpectedOutput is getter method for AIRequest struct
func (ai AIRequest) GetExpectedOutput() []byte {
	return ai.ExpectedOutput
}
