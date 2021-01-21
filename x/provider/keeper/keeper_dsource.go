package keeper

import (
	"context"

	"github.com/oraichain/orai/x/provider/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ListDataSources(ctx context.Context, req *types.ListDataSourcesReq) (*types.ListDataSourcesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDataSources not implemented")
}

func (k Keeper) ListOracleScripts(ctx context.Context, req *types.ListOracleScriptsReq) (*types.ListOracleScriptsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDataSources not implemented")
}

func (k Keeper) ListTestCases(ctx context.Context, req *types.ListTestCasesReq) (*types.ListTestCasesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDataSources not implemented")
}

func (k Keeper) TestCaseInfo(ctx context.Context, req *types.TestCaseInfoReq) (*types.TestCaseInfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDataSources not implemented")
}

func (k Keeper) OracleScriptInfo(ctx context.Context, req *types.OracleScriptInfoReq) (*types.OracleScriptInfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDataSources not implemented")
}
