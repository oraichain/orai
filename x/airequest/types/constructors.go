package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAIRequest is the constructor of the ai request struct
func NewAIRequest(
	requestID string,
	contract sdk.AccAddress,
	creator sdk.AccAddress,
	validators []sdk.ValAddress,
	blockHeight int64,
	fees sdk.Coins,
	input []byte,
	expectedOutput []byte,
) *AIRequest {
	return &AIRequest{
		RequestID:      requestID,
		Contract:       contract,
		Creator:        creator,
		Validators:     validators,
		BlockHeight:    blockHeight,
		Fees:           fees,
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}

// NewMsgSetAIRequest is a constructor function for NewMsgSetAIRequest
func NewMsgSetAIRequest(requestID string, contract sdk.AccAddress, creator sdk.AccAddress, fees string, valCount int64, input []byte, expectedOutput []byte) *MsgSetAIRequest {
	return &MsgSetAIRequest{
		RequestID:      requestID,
		Contract:       contract,
		Creator:        creator,
		ValidatorCount: valCount,
		Fees:           fees,
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}

// NewMsgSetAIRequestRes is a constructor function for NewMsgSetAIRequestRes
func NewMsgSetAIRequestRes(requestID string, contract sdk.AccAddress, creator sdk.AccAddress, fees string, valCount int64, input []byte, expectedOutput []byte) *MsgSetAIRequestRes {
	return &MsgSetAIRequestRes{
		RequestID:      requestID,
		Contract:       contract,
		Creator:        creator,
		ValidatorCount: valCount,
		Fees:           fees,
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}

// NewQueryAIRequestRes is the constructor for the QueryAIRequestRes
func NewQueryAIRequestRes(requestID string, contract sdk.AccAddress, creator sdk.AccAddress, fees sdk.Coins, validators []sdk.ValAddress, blockHeight int64, input []byte, expectedOutput []byte) *QueryAIRequestRes {
	return &QueryAIRequestRes{
		RequestId:   requestID,
		Contract:    contract,
		Creator:     creator,
		Validators:  validators,
		BlockHeight: blockHeight,
		Fees:        fees,
		Input:       input,
		Output:      expectedOutput,
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
