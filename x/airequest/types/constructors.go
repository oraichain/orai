package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMsgSetAIRequest is a constructor function for NewMsgSetAIRequest
func NewMsgSetAIRequest(requestID string, oscriptName string, creator sdk.AccAddress, fees string, valCount int64, input []byte, expectedOutput []byte) *MsgSetAIRequest {
	return &MsgSetAIRequest{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		ValidatorCount:   valCount,
		Fees:             fees,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(aiRequests []AIRequest, params Params) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		AIRequests: aiRequests,
		Params:     params,
	}
}
