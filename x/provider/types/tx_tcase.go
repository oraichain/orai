package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgCreateTestCase) Route() string {
	return RouterKey
}

// Type should return the action
func (msg *MsgCreateTestCase) Type() string {
	return EventTypeCreateTestCase
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateTestCase) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))

}

// ValidateBasic runs stateless checks on the message
func (msg *MsgCreateTestCase) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Contract) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "Name or/and contract address cannot be empty")
	}
	if !IsStringAlphabetic(msg.Name) || !IsStringAlphabetic(msg.Contract) || !IsStringAlphabetic(msg.Description) {
		return sdkerrors.Wrap(ErrCannotSetOracleScript, "Input contains invalid characters")
	}

	if len(msg.Contract) > 1*1024*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The size of the source code is too large!\n")
	}
	// verify contract address
	_, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return err
	}
	return checkFees(msg.Fees)
}

// GetSigners defines whose signature is required
func (msg *MsgCreateTestCase) GetSigners() []sdk.AccAddress {
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
	if len(msg.OldName) == 0 || len(msg.Contract) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name and/or Contract cannot be empty")
	}
	if !IsStringAlphabetic(msg.OldName) || !IsStringAlphabetic(msg.NewName) || !IsStringAlphabetic(msg.Contract) || !IsStringAlphabetic(msg.Description) {
		return sdkerrors.Wrap(ErrCannotSetOracleScript, "Input contains invalid characters")
	}

	if len(msg.Contract) > 1*1024*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The size of the source code is too large!\n")
	}
	// verify contract address
	_, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return err
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
