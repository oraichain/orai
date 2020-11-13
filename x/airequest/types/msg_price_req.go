package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MsgSetPriceRequest defines message for a price prediction request
type MsgSetPriceRequest struct {
	MsgAIRequest MsgSetAIRequest `json:"msg_set_ai_request"`
}

// NewMsgSetPriceRequest is a constructor function for MsgSetPriceRequest
func NewMsgSetPriceRequest(msgSetAIRequest MsgSetAIRequest) MsgSetPriceRequest {
	return MsgSetPriceRequest{
		MsgAIRequest: msgSetAIRequest,
	}
}

// Route should return the name of the module
func (msg MsgSetPriceRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetPriceRequest) Type() string { return "set_price_request" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetPriceRequest) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	err := msg.MsgAIRequest.ValidateBasic()
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetPriceRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetPriceRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MsgAIRequest.Creator}
}
