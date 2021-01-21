package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/oraichain/orai/x/provider/types"
)

var _ types.QueryServer = Keeper{}

// DataSourceInfo implements the Query/DataSourceInfo gRPC method
func (k Keeper) DataSourceInfo(ctx context.Context, req *types.DataSourceInfoReq) (*types.DataSourceInfoRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "data source name query cannot be empty")
	}

	// sdkCtx := sdk.UnwrapSDKContext(ctx)

	// balance := k.GetBalance(sdkCtx, address, req.Denom)

	return &types.DataSourceInfoRes{Name: "ABCD"}, nil
}
