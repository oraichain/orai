package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg MsgCreateTestCase) Route() string {
	return RouterKey
}

// Type should return the action
func (msg MsgCreateTestCase) Type() string {
	return "set_test_case"
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateTestCase) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateTestCase) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Code) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "Name or/and code cannot be empty")
	}
	if len(msg.Code) > 1*1024*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The size of the source code is too large!\n")
	}
	return checkFees(msg.Fees)
}

// GetSigners defines whose signature is required
func (msg MsgCreateTestCase) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Owner.String())
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

// Route should return the name of the module
func (msg MsgEditTestCase) Route() string {
	return RouterKey
}

// Type should return the action
func (msg MsgEditTestCase) Type() string {
	return "edit_test_case"
}

// GetSignBytes encodes the message for signing
func (msg MsgEditTestCase) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))

}

// ValidateBasic runs stateless checks on the message
func (msg MsgEditTestCase) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OldName) == 0 || len(msg.Code) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name and/or Code cannot be empty")
	}
	if len(msg.Code) > 1*1024*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The size of the source code is too large!\n")
	}
	return checkFees(msg.Fees)
}

// GetSigners defines whose signature is required
func (msg MsgEditTestCase) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Owner.String())
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}
