package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	provider "github.com/oraichain/orai/x/provider/types"
)

// NewAIRequest is the constructor of the ai request struct
func NewAIRequest(
	requestID string,
	oscriptName string,
	creator sdk.AccAddress,
	validators []sdk.ValAddress,
	blockHeight int64,
	aiDataSources []provider.AIDataSource,
	testCases []provider.TestCase,
	fees sdk.Coins,
	input []byte,
	expectedOutput []byte,
) *AIRequest {
	return &AIRequest{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		Validators:       validators,
		BlockHeight:      blockHeight,
		AiDataSources:    aiDataSources,
		TestCases:        testCases,
		Fees:             fees,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}

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

// NewMsgSetAIRequestRes is a constructor function for NewMsgSetAIRequestRes
func NewMsgSetAIRequestRes(requestID string, oscriptName string, creator sdk.AccAddress, fees string, valCount int64, input []byte, expectedOutput []byte) *MsgSetAIRequestRes {
	return &MsgSetAIRequestRes{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		ValidatorCount:   valCount,
		Fees:             fees,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}

// NewQueryAIRequestRes is the constructor for the QueryAIRequestRes
func NewQueryAIRequestRes(requestID string, oscriptName string, creator sdk.AccAddress, fees sdk.Coins, validators []sdk.ValAddress, blockHeight int64, ds []provider.AIDataSource, tc []provider.TestCase, input []byte, expectedOutput []byte) *QueryAIRequestRes {
	return &QueryAIRequestRes{
		RequestId:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		Validators:       validators,
		BlockHeight:      blockHeight,
		AiDataSources:    ds,
		TestCases:        tc,
		Fees:             fees,
		Input:            input,
		Output:           expectedOutput,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(aiResults []AIRequestResult, params Params) GenesisState {
	return GenesisState{
		// TODO: Fill out according to your genesis state
		AIRequestResults: aiResults,
		Params:           params,
	}
}
