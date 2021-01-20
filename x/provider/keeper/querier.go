package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/provider/types"
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

// DataSourceInfo implements the Query/DataSourceInfo gRPC method
func (k *Querier) DataSourceInfo(goCtx context.Context, req *types.DataSourceInfoReq) (*types.DataSourceInfoRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "data source name query cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	aiDataSource, err := k.keeper.GetAIDataSource(ctx, req.Name)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}

	return &types.DataSourceInfoRes{
		Name:        aiDataSource.GetName(),
		Owner:       aiDataSource.GetOwner(),
		Contract:    aiDataSource.GetContract(),
		Description: aiDataSource.GetDescription(),
		Fees:        aiDataSource.GetFees(),
	}, nil
}

func (k *Querier) ListDataSources(goCtx context.Context, req *types.ListDataSourcesReq) (*types.ListDataSourcesRes, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var queryResAIDSources []types.AIDataSource

	dSources, err := k.keeper.GetAIDataSources(ctx, uint(req.Page), uint(req.Limit))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}

	// get the total number of data sources
	var count int64 = 0
	iterator := k.keeper.GetAllAIDataSourceNames(ctx)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	// get code of the each dSource
	for _, dSource := range dSources {
		if req.Name == "" || strings.Contains(dSource.Name, req.Name) {
			queryResAIDSources = append(queryResAIDSources, dSource)
		}
	}

	return &types.ListDataSourcesRes{
		AIDataSources: queryResAIDSources,
		Count:         count,
	}, nil

}

func (k *Querier) OracleScriptInfo(goCtx context.Context, req *types.OracleScriptInfoReq) (*types.OracleScriptInfoRes, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "oracle script name query cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	oScript, err := k.keeper.GetOracleScript(ctx, req.Name)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}

	return &types.OracleScriptInfoRes{
		Name:        oScript.GetName(),
		Owner:       oScript.GetOwner(),
		Contract:    oScript.GetContract(),
		Description: oScript.GetDescription(),
		Fees:        oScript.GetMinimumFees(),
		DSources:    oScript.DSources,
		TCases:      oScript.TCases,
	}, nil
}

func (k *Querier) ListOracleScripts(goCtx context.Context, req *types.ListOracleScriptsReq) (*types.ListOracleScriptsRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var queryResOScripts []types.OracleScript

	// collect all the oracle scripts based on the pagination parameters
	oScripts, err := k.keeper.GetOracleScripts(ctx, uint(req.Page), uint(req.Limit))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}

	// get the total number of oracle scripts
	var count int64 = 0
	iterator := k.keeper.GetAllOracleScriptNames(ctx)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	// get code of the each oScript
	for _, oScript := range oScripts {
		if req.Name == "" || strings.Contains(oScript.Name, req.Name) {
			queryResOScripts = append(queryResOScripts, oScript)
		}
	}

	return &types.ListOracleScriptsRes{
		OracleScripts: queryResOScripts,
		Count:         count,
	}, nil
}

func (k *Querier) ListTestCases(goCtx context.Context, req *types.ListTestCasesReq) (*types.ListTestCasesRes, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var queryResTestCases []types.TestCase

	tCases, err := k.keeper.GetTestCases(ctx, uint(req.Page), uint(req.Limit))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, err.Error())
	}

	// get the total number of test cases
	var count int64 = 0
	iterator := k.keeper.GetAllTestCaseNames(ctx)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	// get code of the each tCase
	for _, tCase := range tCases {
		if req.Name == "" || strings.Contains(tCase.Name, req.Name) {
			// create a new queryResOracleScript
			queryResTestCases = append(queryResTestCases, tCase)
		}
	}

	return &types.ListTestCasesRes{
		TestCases: queryResTestCases,
		Count:     count,
	}, nil

}

func (k *Querier) TestCaseInfo(goCtx context.Context, req *types.TestCaseInfoReq) (*types.TestCaseInfoRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "test case name query cannot be empty")
	}

	testCase, err := k.keeper.GetTestCase(ctx, req.Name)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, err.Error())
	}

	return &types.TestCaseInfoRes{
		Name:        testCase.GetName(),
		Owner:       testCase.GetOwner(),
		Contract:    testCase.GetContract(),
		Description: testCase.GetDescription(),
		Fees:        testCase.GetFees(),
	}, nil
}

func (k *Querier) QueryMinFees(goCtx context.Context, req *types.MinFeesReq) (*types.MinFeesRes, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// id of the request
	_, err := k.keeper.GetOracleScript(ctx, req.OracleScriptName)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, err.Error())
	}
	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.keeper.GetDNamesTcNames(ctx, req.OracleScriptName)
	if err != nil {
		return nil, err
	}

	minimumFees, err := k.keeper.GetMinimumFees(ctx, aiDataSources, testCases, int(req.ValNum))
	if err != nil {
		return nil, err
	}

	return &types.MinFeesRes{
		MinFees: minimumFees.AmountOf(types.Denom).String(),
	}, nil
}
