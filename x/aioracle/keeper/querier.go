package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/aioracle/types"
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

// QueryAIOracle implements the Query/QueryAIOracle gRPC method
func (k *Querier) QueryAIOracle(goCtx context.Context, req *types.QueryAIOracleReq) (*types.QueryAIOracleRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GetRequestId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ai request id query cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	AIOracle, err := k.keeper.GetAIOracle(ctx, req.GetRequestId())
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrRequestNotFound, err.Error())
	}

	return types.NewQueryAIOracleRes(
		AIOracle.GetRequestID(), AIOracle.GetContract(),
		AIOracle.GetCreator(), AIOracle.GetFees(),
		AIOracle.GetValidators(), AIOracle.GetBlockHeight(),
		AIOracle.GetInput(),
	), nil
}

// QueryAIOracleIDs implements the Query/QueryAIOracleIDs gRPC method
func (k *Querier) QueryAIOracleIDs(goCtx context.Context, req *types.QueryAIOracleIDsReq) (*types.QueryAIOracleIDsRes, error) {

	var requestIDs []string

	ctx := sdk.UnwrapSDKContext(goCtx)
	iterator := k.keeper.GetPaginatedAIOracles(ctx, uint(req.Page), uint(req.Limit))

	for ; iterator.Valid(); iterator.Next() {
		requestIDs = append(requestIDs, string(iterator.Key()))
	}

	return &types.QueryAIOracleIDsRes{
		RequestIds: requestIDs,
	}, nil
}

// DataSourceContract implements the Query/DataSourceInfo gRPC method
func (k *Querier) DataSourceContract(goCtx context.Context, req *types.QueryDataSourceContract) (*types.ResponseContract, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	contractReq := k.keeper.Cdc.MustMarshalJSON(&types.QueryDataSourceSmartContract{
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

	contractReq := k.keeper.Cdc.MustMarshalJSON(&types.QueryTestCaseSmartContract{
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

	contractReq := k.keeper.Cdc.MustMarshalJSON(&types.QueryOracleScriptSmartContract{
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

	contractReq := k.keeper.Cdc.MustMarshalJSON(&types.QueryDataSourceEntriesSmartContract{
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

	contractReq := k.keeper.Cdc.MustMarshalJSON(&types.QueryTestCaseEntriesSmartContract{
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

// QueryFullRequest implements the Query/QueryFullRequest gRPC method
func (k *Querier) QueryFullRequest(goCtx context.Context, req *types.QueryFullOracleReq) (*types.QueryFullOracleRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// id of the request
	id := req.GetRequestId()

	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "request id query cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	request, err := k.keeper.GetAIOracle(ctx, id)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRequestNotFound, err.Error())
	}
	// collect all the reports of a given request id
	reports := k.keeper.GetReports(ctx, id)

	// collect the result of a given request id

	result, err := k.keeper.GetResult(ctx, id)
	if err != nil {
		// init a nil result showing that the result does not have a result yet
		result = types.NewAIOracleResult(id, nil, types.RequestStatusPending)
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

// QueryMinFees allows user to check the minimum fees for a request that they are going to spend
func (k *Querier) QueryMinFees(goCtx context.Context, req *types.MinFeesReq) (*types.MinFeesRes, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	contractAddr, err := sdk.AccAddressFromBech32(req.GetContractAddr())
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrOScriptNotFound, err.Error())
	}

	minFees, err := k.keeper.calculateMinimumFees(ctx, req.GetTestOnly(), contractAddr, int(req.GetValNum()))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrQueryMinFees, err.Error())
	}

	return &types.MinFeesRes{
		MinimumFees: minFees.String(),
	}, nil
}

func (k *Querier) QueryMinGasPrices(goCtx context.Context, req *types.MinGasPricesReq) (*types.MinGasPricesRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MinGasPricesRes{
		MinGasPrices: ctx.MinGasPrices().String(),
	}, nil
}
