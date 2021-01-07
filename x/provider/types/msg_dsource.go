package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	msgSetDSource  = "set_datasource"
	msgEditDSource = "edit_datasource"
)

// MsgCreateAIDataSource defines a MsgCreateAIDataSource message
type MsgCreateAIDataSource struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        string         `json:"transaction_fee"`
}

// NewMsgCreateAIDataSource is a constructor function for MsgCreateAIDataSource
func NewMsgCreateAIDataSource(name string, code []byte, owner sdk.AccAddress, fees string, des string) MsgCreateAIDataSource {
	return MsgCreateAIDataSource{
		Name:        name,
		Description: des,
		Code:        code,
		Owner:       owner,
		Fees:        fees,
	}
}

// Route should return the name of the module
func (msg MsgCreateAIDataSource) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateAIDataSource) Type() string { return msgSetDSource }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateAIDataSource) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Name) == 0 || len(msg.Code) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "Name and/or Code cannot be empty")
	}
	if len(msg.Code) > 1*1024*1024 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The size of the source code is too large!\n")
	}
	return checkFees(msg.Fees)
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateAIDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateAIDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgEditAIDataSource defines a message for editing a data source in the store
type MsgEditAIDataSource struct {
	OldName     string         `json:"old_name"`
	NewName     string         `json:"new_name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        string         `json:"new_transaction_fee"`
}

// NewMsgEditAIDataSource is a constructor function for MsgEditAIDataSource
func NewMsgEditAIDataSource(oldName string, newName string, code []byte, owner sdk.AccAddress, fees string, des string) MsgEditAIDataSource {
	return MsgEditAIDataSource{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Code:        code,
		Owner:       owner,
		Fees:        fees,
	}
}

// Route should return the name of the module
func (msg MsgEditAIDataSource) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEditAIDataSource) Type() string { return msgEditDSource }

// ValidateBasic runs stateless checks on the message
func (msg MsgEditAIDataSource) ValidateBasic() error {
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

// GetSignBytes encodes the message for signing
func (msg MsgEditAIDataSource) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditAIDataSource) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
