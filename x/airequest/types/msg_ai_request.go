package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetKYCRequest defines message for a KYC request
type MsgSetKYCRequest struct {
	ImageHash    string          `json:"image_hash"`
	ImageName    string          `json:"image_name"`
	MsgAIRequest MsgSetAIRequest `json:"msg_set_ai_request"`
}

// NewMsgSetKYCRequest is a constructor function for MsgSetKYCRequest
func NewMsgSetKYCRequest(imageHash string, imageName string, msgSetAIRequest MsgSetAIRequest) MsgSetKYCRequest {
	return MsgSetKYCRequest{
		ImageHash:    imageHash,
		ImageName:    imageName,
		MsgAIRequest: msgSetAIRequest,
	}
}

// Route should return the name of the module
func (msg MsgSetKYCRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetKYCRequest) Type() string { return "set_kyc_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetKYCRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	err := msg.MsgAIRequest.ValidateBasic()
	if err != nil {
		return err
	}
	if len(msg.ImageHash) == 0 || len(msg.ImageName) == 0 {
		return sdkerrors.Wrap(ErrImageFailedToUnzip, "Image name / hash is not valid")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetKYCRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetKYCRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}

// MsgSetAIRequest defines message for an AI request
type MsgSetAIRequest struct {
	RequestID        string         `json:"request_id"`
	OracleScriptName string         `json:"oscript_name"`
	Creator          sdk.AccAddress `json:"creator"`
	ValidatorCount   int            `json:"validator_count"`
	Fees             string         `json:"transaction_fee"`
	Input            string         `json:"request_input"`
	ExpectedOutput   string         `json:"expected_output"`
}

// NewMsgSetAIRequest is a constructor function for NewMsgSetAIRequest
func NewMsgSetAIRequest(requestID string, oscriptName string, creator sdk.AccAddress, fees string, valCount int, input string, expectedOutput string) MsgSetAIRequest {
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
		return sdkerrors.Wrap(ErrNameIsEmpty, "Name or / and validator count cannot be empty")
	}
	_, err := sdk.ParseCoins(msg.Fees)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidFeeType, err.Error())
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

// MsgSetPriceRequest defines message for a price prediction request
type MsgSetPriceRequest struct {
	MsgAIRequest MsgSetAIRequest `json:"msg_set_ai_request"`
}

// NewMsgSetPriceRequest is a constructor function for MsgSetPriceRequest
func NewMsgSetPriceRequest(msgSetAIRequest MsgSetAIRequest) MsgSetPriceRequest {
	return MsgSetPriceRequest{
		MsgAIRequest: msgSetAIRequest,
	}
}

// Route should return the name of the module
func (msg MsgSetPriceRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetPriceRequest) Type() string { return "set_price_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetPriceRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	err := msg.MsgAIRequest.ValidateBasic()
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetPriceRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetPriceRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}
