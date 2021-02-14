package websocket

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/oraichain/orai/x/websocket/keeper"
)

// NewHandler creates an sdk.Handler for all the provider type messages
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
		case *MsgCreateReport:
			res, err = msgServer.AddReport(goCtx, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}

		return sdk.WrapServiceResult(ctx, res, err)
	}
}
