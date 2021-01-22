package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgAddReporter) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgAddReporter) Type() string { return "add_reporter" }

// ValidateBasic runs stateless checks on the message
func (msg *MsgAddReporter) ValidateBasic() error {
	if msg.Validator.Empty() || msg.Adder.Empty() || msg.Reporter.Empty() {
		return sdkerrors.Wrap(ErrReporterMsgInvalid, "The message attibutes cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgAddReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgAddReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Adder}
}

// Route should return the name of the module
func (msg *MsgRemoveReporter) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgRemoveReporter) Type() string { return "remove_reporter" }

// ValidateBasic runs stateless checks on the message
func (msg *MsgRemoveReporter) ValidateBasic() error {
	if msg.Validator.Empty() || msg.Remover.Empty() || msg.Reporter.Empty() {
		return sdkerrors.Wrap(ErrReporterMsgInvalid, "The message attibutes cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgRemoveReporter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgRemoveReporter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Remover}
}
