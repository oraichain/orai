package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	provider "github.com/oraichain/orai/x/provider/types"
)

// Route should return the name of the module
func (msg *MsgSetAIRequest) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSetAIRequest) Type() string { return EventTypeSetAIRequest }

// ValidateBasic runs stateless checks on the message
func (msg *MsgSetAIRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OracleScriptName) == 0 || msg.ValidatorCount <= 0 {
		return sdkerrors.Wrap(ErrRequestInvalid, "Name or / and validator count cannot be empty")
	}
	// regex allow only alphabet, numeric and underscore characters
	if !provider.IsStringAlphabetic(msg.OracleScriptName) {
		return sdkerrors.Wrap(ErrRequestInvalid, "The oracle script name contains invalid characters")
	}
	_, err := sdk.ParseCoinsNormalized(msg.Fees)
	if err != nil {
		return sdkerrors.Wrap(ErrRequestFeesInvalid, err.Error())
	}
	if len(msg.Fees) == 0 {
		return sdkerrors.Wrap(ErrRequestFeesInvalid, "The fee format is not correct")
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
