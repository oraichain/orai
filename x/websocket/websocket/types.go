package websocket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AIRequest mocks the real AIRequest types
type AIRequest struct {
	RequestID        string         `json:"request_id"`
	OracleScriptName string         `json:"oscript_name"`
	Creator          sdk.AccAddress `json:"creator"`
	ValidatorCount   int64          `json:"validator_count"`
	Input            string         `json:"request_input"`
	ExpectedOutput   string         `json:"expected_output"`
}

// NewAIRequest is the constructor for the mock struct
func NewAIRequest(
	requestID string,
	oscriptName string,
	creator sdk.AccAddress,
	validatorCount int64,
	input string,
	expectedOutput string,
) AIRequest {
	return AIRequest{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		ValidatorCount:   validatorCount,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}
