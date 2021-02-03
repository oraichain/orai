package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetSignBytes encodes the message for signing
func (msg *OracleScript) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *OracleScript) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Owner.String())
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

// Route should return the name of the module
func (msg *OracleScript) Route() string {
	return RouterKey
}

// Type should return the action
func (msg *OracleScript) Type() string {
	return AttributeOracleScriptName
}

// ValidateBasic runs stateless checks on the message
func (msg *OracleScript) ValidateBasic() error {

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name cannot be empty")
	}
	return nil
}
