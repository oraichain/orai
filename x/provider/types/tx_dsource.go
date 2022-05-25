package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgCreateAIDataSource) Route() string {
	return RouterKey
}

// Type should return the action
func (msg *MsgCreateAIDataSource) Type() string {
	return EventTypeSetDataSource
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateAIDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg *MsgCreateAIDataSource) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Name) == 0 || len(msg.Contract) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "Name and/or Contract cannot be empty")
	}

	if !IsStringAlphabetic(msg.Name) || !IsStringAlphabetic(msg.Contract) || !IsStringAlphabetic(msg.Description) {
		return sdkerrors.Wrap(ErrCannotSetDataSource, "Input contains invalid characters")
	}
	// verify contract address
	_, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return err
	}
	return checkFees(msg.Fees)
}

// GetSigners defines whose signature is required
func (msg *MsgCreateAIDataSource) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Owner.String())
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

// Route should return the name of the module
func (msg *MsgEditAIDataSource) Route() string {
	return RouterKey
}

// Type should return the action
func (msg *MsgEditAIDataSource) Type() string {
	return EventTypeEditDataSource
}

// GetSignBytes encodes the message for signing
func (msg *MsgEditAIDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))

}

// ValidateBasic runs stateless checks on the message
func (msg *MsgEditAIDataSource) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OldName) == 0 || len(msg.Contract) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name and/or Contract cannot be empty")
	}
	if !IsStringAlphabetic(msg.OldName) || !IsStringAlphabetic(msg.Contract) || !IsStringAlphabetic(msg.NewName) || !IsStringAlphabetic(msg.Description) {
		return sdkerrors.Wrap(ErrCannotSetDataSource, "input contains invalid characters")
	}

	// verify contract address
	_, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return err
	}

	return checkFees(msg.Fees)
}

// GetSigners defines whose signature is required
func (msg *MsgEditAIDataSource) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Owner.String())
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}
