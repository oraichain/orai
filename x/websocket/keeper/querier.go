package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/types"
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

// OracleInfo implements the Query/DataSourceInfo gRPC method
func (k *Querier) OracleInfo(goCtx context.Context, req *types.QueryContract) (*types.ResponseContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	contractReq := k.keeper.cdc.MustMarshalJSON(req.Request)
	result, err := k.keeper.QueryContract(ctx, req.Contract, contractReq)
	return &types.ResponseContract{
		Data: result,
	}, err
}
