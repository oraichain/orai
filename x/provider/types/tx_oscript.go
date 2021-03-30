package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgCreateOracleScript) Route() string {
	return RouterKey
}

// Type should return the action
func (msg *MsgCreateOracleScript) Type() string {
	return EventTypeSetOracleScript
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))

}

// ValidateBasic runs stateless checks on the message
func (msg *MsgCreateOracleScript) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Name) == 0 || len(msg.Contract) == 0 {
		return sdkerrors.Wrap(ErrCannotSetOracleScript, "Name and/or Contract cannot be empty")
	}
	if !IsStringAlphabetic(msg.Name) || !IsStringAlphabetic(msg.Contract) || !IsStringAlphabetic(msg.Description) {
		return sdkerrors.Wrap(ErrCannotSetOracleScript, "Input contains invalid characters")
	}

	err := validateDSourcesTCases(msg.DataSources, msg.TestCases)
	if err != nil {
		return err
	}

	if len(msg.Contract) > 1*1024*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The size of the source code is too large!\n")
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return err
	}
	return nil
}

func validateDSourcesTCases(dSources, tCases []string) error {
	if len(dSources) == 0 || len(tCases) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data source & test case identifiers cannot be empty")
	}
	for _, dSource := range dSources {
		if len(dSource) == 0 {
			return sdkerrors.Wrap(ErrCannotSetOracleScript, "data source inputs of oracle script cannot be empty")
		}
		if !IsStringAlphabetic(dSource) {
			return sdkerrors.Wrap(ErrCannotSetOracleScript, "Input data source contains invalid characters")
		}
	}

	for _, tCase := range tCases {
		if len(tCase) == 0 {
			return sdkerrors.Wrap(ErrCannotSetOracleScript, "test case inputs of oracle script cannot be empty")
		}
		if !IsStringAlphabetic(tCase) {
			return sdkerrors.Wrap(ErrCannotSetOracleScript, "Input test case contains invalid characters")
		}
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg *MsgCreateOracleScript) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Owner.String())
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

// Route should return the name of the module
func (msg *MsgEditOracleScript) Route() string {
	return RouterKey
}

// Type should return the action
func (msg *MsgEditOracleScript) Type() string {
	return EventTypeEditOracleScript
}

// GetSignBytes encodes the message for signing
func (msg *MsgEditOracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))

}

// ValidateBasic runs stateless checks on the message
func (msg *MsgEditOracleScript) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.OldName) == 0 || len(msg.Contract) == 0 || len(msg.NewName) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name and/or Contract cannot be empty")
	}
	if !IsStringAlphabetic(msg.OldName) || !IsStringAlphabetic(msg.NewName) || !IsStringAlphabetic(msg.Contract) || !IsStringAlphabetic(msg.Description) {
		return sdkerrors.Wrap(ErrCannotSetOracleScript, "Input contains invalid characters")
	}

	err := validateDSourcesTCases(msg.DataSources, msg.TestCases)
	if err != nil {
		return err
	}
	if len(msg.Contract) > 1*1024*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The size of the source code is too large!\n")
	}
	return nil
}

// GetSigners defines whose signature is required
func (msg *MsgEditOracleScript) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Owner.String())
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}
