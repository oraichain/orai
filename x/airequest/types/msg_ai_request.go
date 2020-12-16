package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetAIRequest defines message for an AI request
type MsgSetAIRequest struct {
	RequestID        string         `json:"request_id"`
	OracleScriptName string         `json:"oscript_name"`
	Creator          sdk.AccAddress `json:"creator"`
	ValidatorCount   int            `json:"validator_count"`
	Fees             string         `json:"transaction_fee"`
	Input            []byte         `json:"request_input"`
	ExpectedOutput   []byte         `json:"expected_output"`
}

// NewMsgSetAIRequest is a constructor function for NewMsgSetAIRequest
func NewMsgSetAIRequest(requestID string, oscriptName string, creator sdk.AccAddress, fees string, valCount int, input []byte, expectedOutput []byte) MsgSetAIRequest {
	return MsgSetAIRequest{
		RequestID:        requestID,
		OracleScriptName: oscriptName,
		Creator:          creator,
		ValidatorCount:   valCount,
		Fees:             fees,
		Input:            input,
		ExpectedOutput:   expectedOutput,
	}
}

// Route should return the name of the module
func (msg MsgSetAIRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetAIRequest) Type() string { return "set_ai_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetAIRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OracleScriptName) == 0 || msg.ValidatorCount <= 0 {
		return sdkerrors.Wrap(ErrRequestInvalid, "Name or / and validator count cannot be empty")
	}
	_, err := sdk.ParseCoins(msg.Fees)
	if err != nil {
		return sdkerrors.Wrap(ErrRequestFeesInvalid, err.Error())
	}

	// threshold for the size of the request
	if len(msg.ExpectedOutput)+len(msg.Input) > MaximumRequestBytesThreshold {
		return sdkerrors.Wrap(ErrRequestInvalid, "The request is too large")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetAIRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetAIRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
