package keeper

import (
	// this line is used by starport scaffolding # 1

	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/provider/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {

		switch path[0] {
		// this line is used by starport scaffolding # 2
		case types.QueryOracleScript:
			return queryOracleScript(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryOracleScripts:
			return queryOracleScripts(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryDataSource:
			return queryDataSource(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryDataSources:
			return queryDataSources(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryOracleScriptNames:
			return queryOracleScriptNames(ctx, keeper, legacyQuerierCdc)
		case types.QueryDataSourceNames:
			return queryDataSourceNames(ctx, keeper, legacyQuerierCdc)
		// case types.QueryAIRequest:
		// 	return queryAIRequest(ctx, path[1:], keeper)
		case types.QueryTestCase:
			return queryTestCase(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryTestCases:
			return queryTestCases(ctx, req, keeper, legacyQuerierCdc)
		// case types.QueryAIRequestIDs:
		// 	return queryAIRequestIDs(ctx, keeper)
		case types.QueryTestCaseNames:
			return queryTestCaseNames(ctx, keeper, legacyQuerierCdc)
		// case types.QueryFullRequest:
		// 	return queryFullRequestByID(ctx, path[1:], keeper)
		case types.QueryMinFees:
			return queryMinFees(ctx, req, keeper, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

	}
}

// queryOracleScript queries a oScript given its name
func queryOracleScript(ctx sdk.Context, req abci.RequestQuery, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var params types.QueryOracleScriptRequest

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// need check Name ?
	oScript, err := keeper.GetOracleScript(ctx, params.Name)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}

	result := types.NewQueryResOracleScript(oScript.GetName(), oScript.GetOwner(),
		oScript.GetContract(), oScript.GetDescription(),
		oScript.GetMinimumFees(), oScript.DSources, oScript.TCases)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, result)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryOracleScript queries a list of oracle scripts
func queryOracleScripts(ctx sdk.Context, req abci.RequestQuery, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var params types.QueryOracleScriptsRequest

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var queryResOScripts []types.QueryResOracleScript

	// collect all the oracle scripts based on the pagination parameters
	oScripts, err := keeper.GetOracleScripts(ctx, uint(params.Page), uint(params.Limit))
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}

	// get the total number of oracle scripts
	count := 0
	iterator := keeper.GetAllOracleScriptNames(ctx)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	// get code of the each oScript
	for _, oScript := range oScripts {
		if params.Name == "" || strings.Contains(oScript.Name, params.Name) {
			// create a new queryResOracleScript
			oScriptRes := types.NewQueryResOracleScript(oScript.GetName(), oScript.GetOwner(), oScript.GetContract(),
				oScript.GetDescription(), oScript.GetMinimumFees(), oScript.DSources, oScript.TCases)
			queryResOScripts = append(queryResOScripts, oScriptRes)
		}
	}

	// return the query to the command
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.NewQueryResOracleScripts(queryResOScripts, count))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryDataSource queries a complete Whois struct returned to the user in []byte
func queryDataSource(ctx sdk.Context, req abci.RequestQuery, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var params types.QueryDataSourceRequest
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	aiDataSource, err := keeper.GetAIDataSource(ctx, params.Name)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}

	aiDataSourceRes := types.NewQueryResAIDataSource(aiDataSource.GetName(), aiDataSource.GetOwner(), aiDataSource.GetContract(), aiDataSource.GetDescription(), aiDataSource.Fees)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, aiDataSourceRes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryDataSources queries a list of data sources
func queryDataSources(ctx sdk.Context, req abci.RequestQuery, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDataSourcesRequest

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var queryResAIDSources []types.QueryResAIDataSource

	dSources, err := keeper.GetAIDataSources(ctx, uint(params.Page), uint(params.Limit))
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}

	// get the total number of data sources
	count := 0
	iterator := keeper.GetAllAIDataSourceNames(ctx)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	// get code of the each dSource
	for _, dSource := range dSources {
		if params.Name == "" || strings.Contains(dSource.Name, params.Name) {

			// create a new queryResOracleScript
			queryResAIDSource := types.NewQueryResAIDataSource(dSource.GetName(), dSource.GetOwner(), dSource.GetContract(),
				dSource.GetDescription(), dSource.Fees)

			queryResAIDSources = append(queryResAIDSources, queryResAIDSource)
		}
	}

	// return the query to the command
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.NewQueryResAIDataSources(queryResAIDSources, count))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryOracleScriptNames returns all the oScript names in the store
func queryOracleScriptNames(ctx sdk.Context, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var namesList types.QueryResOracleScriptNames

	iterator := keeper.GetAllOracleScriptNames(ctx)

	for ; iterator.Valid(); iterator.Next() {
		namesList = append(namesList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, namesList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryDataSourceNames returns all the data source names in the store
func queryDataSourceNames(ctx sdk.Context, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var namesList types.QueryResAIDataSourceNames

	iterator := keeper.GetAllAIDataSourceNames(ctx)

	for ; iterator.Valid(); iterator.Next() {
		namesList = append(namesList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, namesList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryTestCase queries an AI request test case
func queryTestCase(ctx sdk.Context, req abci.RequestQuery, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTestCaseRequest
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	testCase, err := keeper.GetTestCase(ctx, params.Name)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrTestCaseNotFound, err.Error())
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc,
		types.NewQueryResTestCase(testCase.GetName(), testCase.GetOwner(), testCase.GetContract(),
			testCase.GetDescription(), testCase.Fees))

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryTestCases queries a list of test cases
func queryTestCases(ctx sdk.Context, req abci.RequestQuery, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTestCasesRequest

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var queryResTestCases []types.QueryResTestCase

	tCases, err := keeper.GetTestCases(ctx, uint(params.Page), uint(params.Limit))
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrTestCaseNotFound, err.Error())
	}

	// get the total number of test cases
	count := 0
	iterator := keeper.GetAllOracleScriptNames(ctx)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	// get code of the each tCase
	for _, tCase := range tCases {
		if params.Name == "" || strings.Contains(tCase.Name, params.Name) {

			// create a new queryResOracleScript
			queryResTestCases = append(queryResTestCases,
				types.NewQueryResTestCase(tCase.GetDescription(), tCase.GetOwner(), tCase.GetContract(),
					tCase.GetDescription(), tCase.Fees))
		}
	}

	// return the query to the command
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.NewQueryResTestCases(queryResTestCases, count))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryTestCaseNames returns all the test case names in the store
func queryTestCaseNames(ctx sdk.Context, keeper *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var testCaseNames types.QueryResTestCaseNames

	iterator := keeper.GetAllTestCaseNames(ctx)

	for ; iterator.Valid(); iterator.Next() {
		testCaseNames = append(testCaseNames, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, testCaseNames)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryMinFees(ctx sdk.Context, req abci.RequestQuery, k *Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var params types.MinFeesRequest

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// number of validator

	// id of the request
	_, err := k.GetOracleScript(ctx, params.OracleScriptName)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, err.Error())
	}
	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.GetDNamesTcNames(ctx, params.OracleScriptName)
	if err != nil {
		return nil, err
	}

	minimumFees, err := k.GetMinimumFees(ctx, aiDataSources, testCases, params.ValNum)
	if err != nil {
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, types.NewQueryResMinFees(minimumFees.AmountOf(types.Denom).String()))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
