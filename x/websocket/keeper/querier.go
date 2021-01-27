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

// OracleContract implements the Query/DataSourceInfo gRPC method
func (k *Querier) OracleContract(goCtx context.Context, req *types.QueryOracleContract) (*types.ResponseContract, error) {
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

// DataSourceContract implements the Query/DataSourceInfo gRPC method
func (k *Querier) DataSourceContract(goCtx context.Context, req *types.QueryDataSourceContract) (*types.ResponseContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	datasource, err := k.keeper.providerKeeper.GetAIDataSource(ctx, req.Name)
	if err != nil {
		// already wrapped at providerKeeper
		return nil, err
	}

	result, err := k.keeper.QueryContract(ctx, sdk.AccAddress(datasource.Contract), []byte(`{}`))
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
	testcase, err := k.keeper.providerKeeper.GetTestCase(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	datasource, err := k.keeper.providerKeeper.GetAIDataSource(ctx, req.DataSourcName)
	if err != nil {
		// already wrapped at providerKeeper
		return nil, err
	}

	// re-bind
	req.Request.DataSourceContract = sdk.AccAddress(datasource.Contract)
	contractReq := k.keeper.cdc.MustMarshalJSON(req.Request)
	result, err := k.keeper.QueryContract(ctx, sdk.AccAddress(testcase.Contract), contractReq)
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
	oscript, err := k.keeper.providerKeeper.GetTestCase(ctx, req.Name)
	if err != nil {
		// already wrapped at providerKeeper
		return nil, err
	}

	contractReq := k.keeper.cdc.MustMarshalJSON(req.Request)
	result, err := k.keeper.QueryContract(ctx, sdk.AccAddress(oscript.Contract), contractReq)
	return &types.ResponseContract{
		Data: result,
	}, err
}
