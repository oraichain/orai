package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	airequest "github.com/oraichain/orai/x/airequest"
	"github.com/oraichain/orai/x/airesult/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	keeper *Keeper
}

// NewQuerier return querier implementation
func NewQuerier(keeper *Keeper) *Querier {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = &Querier{}

// QueryFullRequest implements the Query/QueryFullRequest gRPC method
func (k *Querier) QueryFullRequest(goCtx context.Context, req *types.QueryFullRequestReq) (*types.QueryFullRequestRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// id of the request
	id := req.GetRequestId()

	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "request id query cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	request, err := k.keeper.aiRequestKeeper.GetAIRequest(ctx, id)
	if err != nil {
		return nil, sdkerrors.Wrapf(airequest.ErrRequestNotFound, err.Error())
	}
	// collect all the reports of a given request id
	reports := k.keeper.webSocketKeeper.GetReports(ctx, id)

	// collect the result of a given request id

	result, err := k.keeper.GetResult(ctx, id)
	if err != nil {
		// init a nil result showing that the result does not have a result yet
		result = types.NewAIRequestResult(id, nil, types.RequestStatusPending)
	}

	return types.NewQueryFullRequestRes(*request, reports, *result), nil
}

// QueryReward implements the Query/QueryReward gRPC method
func (k *Querier) QueryReward(goCtx context.Context, req *types.QueryRewardReq) (*types.QueryRewardRes, error) {

	// id of the request
	blockHeight := req.GetBlockHeight()
	blockHeightInt, err := strconv.Atoi(blockHeight)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrBlockHeightInvalid, err.Error())
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	// get reward
	reward, err := k.keeper.GetReward(ctx, int64(blockHeightInt))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRewardNotfound, err.Error())
	}

	return types.NewQueryRewardRes(*reward), nil
}
