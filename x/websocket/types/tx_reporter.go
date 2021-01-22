package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


// NewMsgAddReporter is a constructor function for MsgAddReporter
func NewMsgAddReporter(validator sdk.ValAddress, reporter sdk.AccAddress, adder sdk.AccAddress) *MsgAddReporter {
	return &MsgAddReporter{
		Adder:     adder,
		Validator: validator,
		Reporter:  reporter,
	}
}

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



// NewMsgRemoveReporter is a constructor function for MsgRemoveReporter
func NewMsgRemoveReporter(validator sdk.ValAddress, reporter sdk.AccAddress, remover sdk.AccAddress) MsgRemoveReporter {
	return MsgRemoveReporter{
		Remover:   remover,
		Validator: validator,
		Reporter:  reporter,
	}
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
