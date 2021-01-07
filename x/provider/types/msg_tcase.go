package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	msgSetTCase  = "set_test_case"
	msgEditTCase = "edit_test_case"
)

// MsgCreateTestCase defines message for an AI request test case
type MsgCreateTestCase struct {
	Name        string         `json:"test_case_name"`
	Owner       sdk.AccAddress `json:"owner"`
	Code        []byte         `json:"code"`
	Fees        string         `json:"transaction_fee"`
	Description string         `json:"description"`
}

// NewMsgCreateTestCase is a constructor function for MsgCreateTestCase
func NewMsgCreateTestCase(name string, code []byte, owner sdk.AccAddress, fees string, des string) MsgCreateTestCase {
	return MsgCreateTestCase{
		Name:        name,
		Description: des,
		Owner:       owner,
		Code:        code,
		Fees:        fees,
	}
}

// Route should return the name of the module
func (msg MsgCreateTestCase) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateTestCase) Type() string { return msgSetTCase }

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

// GetSignBytes encodes the message for signing
func (msg MsgCreateTestCase) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateTestCase) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgEditTestCase defines a message for editing a test case in the store
type MsgEditTestCase struct {
	OldName     string         `json:"old_name"`
	NewName     string         `json:"new_name"`
	Description string         `json:"description"`
	Code        []byte         `json:"code"`
	Owner       sdk.AccAddress `json:"owner"`
	Fees        string         `json:"new_transaction_fee"`
}

// NewMsgEditTestCase is a constructor function for MsgEditTestCase
func NewMsgEditTestCase(oldName string, newName string, code []byte, owner sdk.AccAddress, fees string, des string) MsgEditTestCase {
	return MsgEditTestCase{
		OldName:     oldName,
		NewName:     newName,
		Description: des,
		Code:        code,
		Owner:       owner,
		Fees:        fees,
	}
}

// Route should return the name of the module
func (msg MsgEditTestCase) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEditTestCase) Type() string { return msgEditTCase }

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

// GetSignBytes encodes the message for signing
func (msg MsgEditTestCase) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditTestCase) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
