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
<<<<<<< HEAD
	if len(msg.Name) == 0 || len(msg.Contract) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name and/or Contract cannot be empty")
	}
	if len(msg.DataSources) == 0 || len(msg.TestCases) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data source & test case identifiers cannot be empty")
	}
	if len(msg.Contract) > MaximumContractLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The length of the smart contract address is too large!\n")
	}
=======
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name and/or Contract cannot be empty")
	}
	// verify contract address
	_, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return err
	}

	if len(msg.DataSources) == 0 || len(msg.TestCases) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data source & test case identifiers cannot be empty")
	}

>>>>>>> 59e687a529c4f597ee6c64b6ec8370f0da512597
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
	if len(msg.DataSources) == 0 || len(msg.TestCases) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data source & test case identifiers cannot be empty")
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
