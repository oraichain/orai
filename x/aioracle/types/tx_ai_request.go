package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgSetAIOracleReq) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSetAIOracleReq) Type() string { return EventTypeSetAIOracle }

// ValidateBasic runs stateless checks on the message
func (msg *MsgSetAIOracleReq) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Contract) == 0 || msg.ValidatorCount <= 0 {
		return sdkerrors.Wrap(ErrRequestInvalid, "Name or / and validator count cannot be empty")
	}
	fees, err := sdk.ParseCoinsNormalized(msg.Fees)
	if err != nil {
		return sdkerrors.Wrap(ErrRequestFeesInvalid, err.Error())
	}
	if fees.Len() == 0 {
		return sdkerrors.Wrap(ErrRequestFeesInvalid, "The fee format is not correct or empty fees")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgSetAIOracleReq) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSetAIOracleReq) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
