package keeper

import (
	"context"
	"fmt"
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

	contractReq, err := k.keeper.Cdc.MarshalJSON(&types.QueryDataSourceSmartContract{
		Get: req.Request,
	})
	if err != nil {
		return nil, err
	}
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

	contractReq, err := k.keeper.Cdc.MarshalJSON(&types.QueryTestCaseSmartContract{
		Test: req.Request,
	})
	if err != nil {
		return nil, err
	}
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

	contractReq, err := k.keeper.Cdc.MarshalJSON(&types.QueryOracleScriptSmartContract{
		Aggregate: req.Request,
	})
	if err != nil {
		return nil, err
	}
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

	contractReq, err := k.keeper.Cdc.MarshalJSON(&types.QueryDataSourceEntriesSmartContract{
		GetDataSources: req.Request,
	})
	if err != nil {
		return nil, err
	}
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

	contractReq, err := k.keeper.Cdc.MarshalJSON(&types.QueryTestCaseEntriesSmartContract{
		GetTestCases: req.Request,
	})
	if err != nil {
		return nil, err
	}
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
	if request.TestOnly {

	}
	// collect all the reports of a given request id
	reports := k.keeper.GetReports(ctx, id)
	if len(reports) == 0 {
		reports = []types.Report{}
	}
	// collect all the reports of a given request id
	tcReports := k.keeper.GetTestCaseReports(ctx, id)
	if len(tcReports) == 0 {
		tcReports = []types.TestCaseReport{}
	}

	// collect the result of a given request id

	result, err := k.keeper.GetResult(ctx, id)
	if err != nil {
		// init a nil result showing that the result does not have a result yet
		result = types.NewAIOracleResult(id, nil, types.RequestStatusPending)
	}

	return types.NewQueryFullRequestRes(*request, reports, tcReports, *result), nil
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

	minFees, err := k.calculateMinimumFees(ctx, req.GetTestOnly(), contractAddr, int(req.GetValNum()))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrQueryMinFees, err.Error())
	}

	return &types.MinFeesRes{
		MinimumFees: minFees.String(),
	}, nil
}

func (k *Querier) calculateMinimumFees(ctx sdk.Context, isTestOnly bool, contractAddr sdk.AccAddress, numVal int) (sdk.Coins, error) {

	goCtx := sdk.WrapSDKContext(ctx)
	entries := &types.ResponseEntryPointContract{}
	var err error
	if isTestOnly {
		entries, err = k.TestCaseEntries(goCtx, &types.QueryTestCaseEntriesContract{
			Contract: contractAddr,
			Request:  &types.EmptyParams{},
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err = k.DataSourceEntries(goCtx, &types.QueryDataSourceEntriesContract{
			Contract: contractAddr,
			Request:  &types.EmptyParams{},
		})
		if err != nil {
			return nil, err
		}
	}
	minFees := sdk.Coins{}
	for _, entry := range entries.GetData() {
		minFees = minFees.Add(entry.GetProviderFees()...)
	}
	valFees := k.keeper.CalculateValidatorFees(ctx, minFees)
	minFees = minFees.Add(valFees...)
	// since there are more than 1 validator, we need to multiply those fees
	minFees, _ = sdk.NewDecCoinsFromCoins(minFees...).MulDec(sdk.NewDec(int64(numVal))).TruncateDecimal()
	return minFees, nil
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

// Params queries the staking parameters
func (k *Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// Params queries the staking parameters
func (k *Querier) Param(c context.Context, paramReq *types.QueryParamRequest) (*types.QueryParamResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	resp := k.keeper.GetParam(ctx, []byte(paramReq.Param))
	return &types.QueryParamResponse{Param: resp}, nil
}

// containsVal returns whether the given slice of validators contains the target validator.
func (k *Querier) ContainsVal(vals []sdk.ValAddress, target sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(target) {
			return true
		}
	}
	return false
}

func (k *Querier) validateBasic(ctx sdk.Context, req *types.AIOracle, rep *types.BaseReport) bool {
	// Check if validator exists and active
	_, isExist := k.keeper.StakingKeeper.GetValidator(ctx, rep.GetValidatorAddress())
	if !isExist {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("error in validating the report: validator does not exist"))
		return false
	}
	if !k.ContainsVal(req.GetValidators(), rep.GetValidatorAddress()) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Validator %v does not exist in the list of request validators", rep.GetValidatorAddress().String()))
		return false
	}
	if len(rep.GetValidatorAddress()) <= 0 || len(rep.GetRequestId()) <= 0 || rep.GetBlockHeight() <= 0 {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Report basic information is invalid: %v", rep))
		return false
	}
	return true
}
