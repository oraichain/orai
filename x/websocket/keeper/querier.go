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

// DataSourceContract implements the Query/DataSourceInfo gRPC method
func (k *Querier) DataSourceContract(goCtx context.Context, req *types.QueryDataSourceContract) (*types.ResponseContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	contractReq := k.keeper.cdc.MustMarshalJSON(&types.QueryDataSourceSmartContract{
		Get: req.Request,
	})
	result, err := k.keeper.QueryContract(ctx, req.Contract, contractReq)
	return &types.ResponseContract{
		Data: result,
	}, err
}

// TestCaseContract implements the Query/DataSourceInfo gRPC method
func (k *Querier) TestCaseContract(goCtx context.Context, req *types.QueryTestCaseContract) (*types.ResponseContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	contractReq := k.keeper.cdc.MustMarshalJSON(&types.QueryTestCaseSmartContract{
		Test: req.Request,
	})

	result, err := k.keeper.QueryContract(ctx, req.Contract, contractReq)
	return &types.ResponseContract{
		Data: result,
	}, err
}

// OracleScriptContract implements the Query/DataSourceInfo gRPC method
func (k *Querier) OracleScriptContract(goCtx context.Context, req *types.QueryOracleScriptContract) (*types.ResponseContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	contractReq := k.keeper.cdc.MustMarshalJSON(&types.QueryOracleScriptSmartContract{
		Aggregate: req.Request,
	})

	result, err := k.keeper.QueryContract(ctx, req.Contract, contractReq)
	return &types.ResponseContract{
		Data: result,
	}, err
}

// DataSourceEntries implements the Query/DataSourceInfo gRPC method
func (k *Querier) DataSourceEntries(goCtx context.Context, req *types.QueryDataSourceEntriesContract) (*types.ResponseEntryPointContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	contractReq := k.keeper.cdc.MustMarshalJSON(&types.QueryDataSourceEntriesSmartContract{
		GetDataSources: req.Request,
	})

	result, err := k.keeper.QueryContract(ctx, req.Contract, contractReq)
	if err != nil {
		return nil, err
	}
	var entries []*types.EntryPoint
	err = types.ModuleCdc.Amino.UnmarshalJSON(result, &entries)

	return &types.ResponseEntryPointContract{
		Data: entries,
	}, err
}

// TestCaseEntries implements the Query/DataSourceInfo gRPC method
func (k *Querier) TestCaseEntries(goCtx context.Context, req *types.QueryTestCaseEntriesContract) (*types.ResponseEntryPointContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	contractReq := k.keeper.cdc.MustMarshalJSON(&types.QueryTestCaseEntriesSmartContract{
		GetTestCases: req.Request,
	})

	result, err := k.keeper.QueryContract(ctx, req.Contract, contractReq)
	if err != nil {
		return nil, err
	}
	var entries []*types.EntryPoint
	err = types.ModuleCdc.Amino.UnmarshalJSON(result, &entries)

	return &types.ResponseEntryPointContract{
		Data: entries,
	}, err
}
