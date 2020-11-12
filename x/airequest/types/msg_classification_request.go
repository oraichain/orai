package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetClassificationRequest defines message for a KYC request
type MsgSetClassificationRequest struct {
	ImageHash    string          `json:"image_hash"`
	ImageName    string          `json:"image_name"`
	MsgAIRequest MsgSetAIRequest `json:"msg_set_ai_request"`
}

// NewMsgSetClassificationRequest is a constructor function for MsgSetKYCRequest
func NewMsgSetClassificationRequest(imageHash string, imageName string, msgSetAIRequest MsgSetAIRequest) MsgSetClassificationRequest {
	return MsgSetClassificationRequest{
		ImageHash:    imageHash,
		ImageName:    imageName,
		MsgAIRequest: msgSetAIRequest,
	}
}

// Route should return the name of the module
func (msg MsgSetClassificationRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetClassificationRequest) Type() string { return "set_classification_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetClassificationRequest) ValidateBasic() error {
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
func (msg MsgSetClassificationRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetClassificationRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}
