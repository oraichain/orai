package airequest

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/oraichain/orai/x/airequest/keeper"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k *Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		goCtx := sdk.WrapSDKContext(ctx)
		var (
			res proto.Message
			err error
		)
		// this is for server to broadcast, query is for client to query
		switch msg := msg.(type) {
		case *MsgSetAIRequest:
			res, err = msgServer.CreateAIRequest(goCtx, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		return sdk.WrapServiceResult(ctx, res, err)
	}
}
