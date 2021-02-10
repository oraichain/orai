package keeper

import (
	// this line is used by starport scaffolding # 1

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/provider/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

// Query statements for data sources
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryDataSource        = "datasource"
	QueryDataSources       = "datasources"
	QueryDataSourceNames   = "dnames"
	QueryOracleScript      = "oscript"
	QueryOracleScripts     = "oscripts"
	QueryOracleScriptNames = "onames"
	QueryMinFees           = "min_fees"
	QueryTestCase          = "testcase"
	QueryTestCases         = "testcases"
	QueryTestCaseNames     = "tcnames"
)

// NewLegacyQuerier implementing for legacy querier
func NewLegacyQuerier(querier *Querier, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {

		var (
			rsp interface{}
			err error
		)
		goCtx := ctx.Context()

		switch path[0] {

		case QueryOracleScript:
			params := &types.OracleScriptInfoReq{}
			if err := legacyQuerierCdc.UnmarshalJSON(req.Data, params); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			rsp, err = querier.OracleScriptInfo(goCtx, params)
		case QueryOracleScripts:
			params := &types.ListOracleScriptsReq{}
			if err := legacyQuerierCdc.UnmarshalJSON(req.Data, params); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			rsp, err = querier.ListOracleScripts(goCtx, params)
		case QueryDataSource:
			params := &types.DataSourceInfoReq{}
			if err := legacyQuerierCdc.UnmarshalJSON(req.Data, params); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			rsp, err = querier.DataSourceInfo(goCtx, params)
		case QueryDataSources:
			params := &types.ListDataSourcesReq{}
			if err := legacyQuerierCdc.UnmarshalJSON(req.Data, params); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			rsp, err = querier.ListDataSources(goCtx, params)

		case QueryTestCase:
			params := &types.TestCaseInfoReq{}
			if err := legacyQuerierCdc.UnmarshalJSON(req.Data, params); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			rsp, err = querier.TestCaseInfo(goCtx, params)
		case QueryTestCases:
			params := &types.ListTestCasesReq{}
			if err := legacyQuerierCdc.UnmarshalJSON(req.Data, params); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			rsp, err = querier.ListTestCases(goCtx, params)

		case QueryMinFees:
			params := &types.MinFeesReq{}
			if err := legacyQuerierCdc.UnmarshalJSON(req.Data, params); err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
			}
			rsp, err = querier.QueryMinFees(goCtx, params)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		// wrap serialization
		if err != nil {
			return nil, err
		}
		bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, rsp)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}
}
