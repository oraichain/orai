package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/types"
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

// QueryAIRequest implements the Query/QueryAIRequest gRPC method
func (k *Querier) QueryAIRequest(goCtx context.Context, req *types.QueryAIRequestReq) (*types.QueryAIRequestRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetRequestId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ai request id query cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	aiRequest, err := k.keeper.GetAIRequest(ctx, req.GetRequestId())
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrRequestNotFound, err.Error())
	}

	return types.NewQueryAIRequestRes(
		aiRequest.GetRequestID(), aiRequest.GetOracleScriptName(),
		aiRequest.GetCreator(), aiRequest.GetFees(),
		aiRequest.GetValidators(), aiRequest.GetBlockHeight(),
		aiRequest.GetAiDataSources(), aiRequest.GetTestCases(),
		aiRequest.GetInput(), aiRequest.GetExpectedOutput(),
	), nil
}

// QueryAIRequestIDs implements the Query/QueryAIRequestIDs gRPC method
func (k *Querier) QueryAIRequestIDs(goCtx context.Context, req *types.QueryAIRequestIDsReq) (*types.QueryAIRequestIDsRes, error) {

	var requestIDs []string

	ctx := sdk.UnwrapSDKContext(goCtx)
	iterator := k.keeper.GetPaginatedAIRequests(ctx, uint(req.Page), uint(req.Limit))

	for ; iterator.Valid(); iterator.Next() {
		requestIDs = append(requestIDs, string(iterator.Key()))
	}

	return &types.QueryAIRequestIDsRes{
		RequestIds: requestIDs,
	}, nil
}
