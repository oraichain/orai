package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgSetAIRequestReq) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSetAIRequestReq) Type() string { return EventTypeSetAIRequest }

// ValidateBasic runs stateless checks on the message
func (msg *MsgSetAIRequestReq) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Contract) == 0 || msg.ValidatorCount <= 0 {
		return sdkerrors.Wrap(ErrRequestInvalid, "Name or / and validator count cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrap(ErrRequestInvalid, "Contract address is invalid")
	}
	_, err = sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrap(ErrRequestInvalid, "creator address is invalid")
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
func (msg *MsgSetAIRequestReq) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSetAIRequestReq) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{creator}
}
