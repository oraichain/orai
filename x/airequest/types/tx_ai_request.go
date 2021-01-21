package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgSetAIRequest) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSetAIRequest) Type() string { return "set_ai_request" }

// ValidateBasic runs stateless checks on the message
func (msg *MsgSetAIRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OracleScriptName) == 0 || msg.ValidatorCount <= 0 {
		return sdkerrors.Wrap(ErrRequestInvalid, "Name or / and validator count cannot be empty")
	}
	_, err := sdk.ParseCoinsNormalized(msg.Fees)
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
func (msg *MsgSetAIRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSetAIRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
