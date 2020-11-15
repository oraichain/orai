package airequest

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
)

// NewHandler creates an sdk.Handler for all the airequest type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetKYCRequest:
			return handleMsgSetKYCRequest(ctx, k, msg)
		case types.MsgSetClassificationRequest:
			return handleMsgSetClassificationRequest(ctx, k, msg)
		case types.MsgSetOCRRequest:
			return handleMsgSetOCRRequest(ctx, k, msg)
		case types.MsgSetPriceRequest:
			return handleMsgSetPriceRequest(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
