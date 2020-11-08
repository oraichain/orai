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
	Input            string                   `json:"request_input"`
	ExpectedOutput   string                   `json:"expected_output"`
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
	input string,
	expectedOutput string,
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
