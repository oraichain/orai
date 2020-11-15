package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetOCRRequest defines message for a OCR request
type MsgSetOCRRequest struct {
	ImageHash    string          `json:"image_hash"`
	ImageName    string          `json:"image_name"`
	MsgAIRequest MsgSetAIRequest `json:"msg_set_ai_request"`
}

// NewMsgSetOCRRequest is a constructor function for MsgSetOCRRequest
func NewMsgSetOCRRequest(imageHash string, imageName string, msgSetAIRequest MsgSetAIRequest) MsgSetOCRRequest {
	return MsgSetOCRRequest{
		ImageHash:    imageHash,
		ImageName:    imageName,
		MsgAIRequest: msgSetAIRequest,
	}
}

// Route should return the name of the module
func (msg MsgSetOCRRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetOCRRequest) Type() string { return "set_ocr_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetOCRRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	fmt.Println("messgae: ", msg)
	err := msg.MsgAIRequest.ValidateBasic()
	if err != nil {
		fmt.Println("ERROR IN VALIDATING MSG AI REQUEST", err)
		return err
	}
	if len(msg.ImageHash) == 0 || len(msg.ImageName) == 0 {
		return sdkerrors.Wrap(ErrImageFailedToUnzip, "Image name / hash is not valid")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetOCRRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetOCRRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}
