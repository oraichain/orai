package keeper

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/provider/types"
)

// NewQuerier creates a new querier for provider clients.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		// TODO: Put the modules query routes
		case types.QueryOracleScript:
			return queryOracleScript(ctx, path[1:], keeper)
		case types.QueryOracleScripts:
			return queryOracleScripts(ctx, keeper, req)
		case types.QueryDataSource:
			return queryDataSource(ctx, path[1:], keeper)
		case types.QueryDataSources:
			return queryDataSources(ctx, keeper, req)
		case types.QueryOracleScriptNames:
			return queryOracleScriptNames(ctx, keeper)
		case types.QueryDataSourceNames:
			return queryDataSourceNames(ctx, keeper)
		case types.QueryAIRequest:
			return queryAIRequest(ctx, path[1:], keeper)
		case types.QueryTestCase:
			return queryTestCase(ctx, path[1:], keeper)
		case types.QueryTestCases:
			return queryTestCases(ctx, keeper, req)
		case types.QueryAIRequestIDs:
			return queryAIRequestIDs(ctx, keeper)
		case types.QueryTestCaseNames:
			return queryTestCaseNames(ctx, keeper)
		case types.QueryFullRequest:
			return queryFullRequestByID(ctx, path[1:], keeper)
		case types.QueryMinFees:
			return queryMinFees(ctx, path[1:], keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown provider query")
		}
	}
}

// queryOracleScript queries a oScript given its name
func queryOracleScript(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	// tsao cho nay lai lay path[0] ?
	oScript, err := keeper.GetOracleScript(ctx, path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}

	// get code of the oScript
	code, err := keeper.fileCache.GetFile(types.OracleScriptStoreKeyString(oScript.Name))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCodeNotFound, err.Error())
	}

	executable := base64.StdEncoding.EncodeToString(code)

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResOracleScript(oScript.Name, oScript.Owner, executable, oScript.Description, oScript.MinimumFees))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryOracleScript queries a list of oracle scripts
func queryOracleScripts(ctx sdk.Context, keeper Keeper, req abci.RequestQuery) ([]byte, error) {
	// tsao cho nay lai lay path[0] ?

	var queryResOScripts []types.QueryResOracleScript

	// parse limit and offset from the query message data
	pagiSlice := strings.Split(string(req.GetData()[:]), "-")
	page, err := strconv.Atoi(pagiSlice[0])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrPaginationInputInvalid, err.Error())
	}
	limit, err := strconv.Atoi(pagiSlice[1])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrPaginationInputInvalid, err.Error())
	}

	// collect all the oracle scripts based on the pagination parameters
	oScripts, err := keeper.GetOracleScripts(ctx, uint(page), uint(limit))
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrOracleScriptNotFound, err.Error())
	}

	// get code of the each oScript
	for _, oScript := range oScripts {
		code, err := keeper.fileCache.GetFile(types.OracleScriptStoreKeyString(oScript.Name))
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrCodeNotFound, err.Error())
		}
		executable := base64.StdEncoding.EncodeToString(code)

		// create a new queryResOracleScript
		queryResOScripts = append(queryResOScripts, types.NewQueryResOracleScript(oScript.Name, oScript.Owner, executable, oScript.Description, oScript.MinimumFees))
	}

	// return the query to the command
	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResOracleScripts(queryResOScripts, len(oScripts)))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryDataSource queries a complete Whois struct returned to the user in []byte
func queryDataSource(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	aiDataSource, err := keeper.GetAIDataSource(ctx, path[0])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}

	// get code of the data source
	code, err := keeper.fileCache.GetFile(types.DataSourceStoreKeyString(aiDataSource.Name))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCodeNotFound, err.Error())
	}

	executable := base64.StdEncoding.EncodeToString(code)

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResAIDataSource(aiDataSource.Name, aiDataSource.Owner, executable, aiDataSource.Description))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryDataSources queries a list of data sources
func queryDataSources(ctx sdk.Context, keeper Keeper, req abci.RequestQuery) ([]byte, error) {
	// tsao cho nay lai lay path[0] ?

	var queryResAIDSources []types.QueryResAIDataSource

	// parse limit and offset from the query message data
	pagiSlice := strings.Split(string(req.GetData()[:]), "-")
	page, err := strconv.Atoi(pagiSlice[0])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrPaginationInputInvalid, err.Error())
	}
	limit, err := strconv.Atoi(pagiSlice[1])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrPaginationInputInvalid, err.Error())
	}

	dSources, err := keeper.GetAIDataSources(ctx, uint(page), uint(limit))
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrDataSourceNotFound, err.Error())
	}

	// get code of the each dSource
	for _, dSource := range dSources {
		code, err := keeper.fileCache.GetFile(types.DataSourceStoreKeyString(dSource.Name))
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrCodeNotFound, err.Error())
		}
		executable := base64.StdEncoding.EncodeToString(code)

		// create a new queryResOracleScript
		queryResAIDSources = append(queryResAIDSources, types.NewQueryResAIDataSource(dSource.Name, dSource.Owner, executable, dSource.Description))
	}

	// return the query to the command
	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResAIDataSources(queryResAIDSources, len(dSources)))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryOracleScriptNames returns all the oScript names in the store
func queryOracleScriptNames(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var namesList types.QueryResOracleScriptNames

	iterator := keeper.GetAllOracleScriptNames(ctx)

	for ; iterator.Valid(); iterator.Next() {
		namesList = append(namesList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, namesList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryDataSourceNames returns all the data source names in the store
func queryDataSourceNames(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var namesList types.QueryResAIDataSourceNames

	iterator := keeper.GetAllAIDataSourceNames(ctx)

	for ; iterator.Valid(); iterator.Next() {
		namesList = append(namesList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, namesList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryAIRequest queries an AI request given its request ID
func queryAIRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	// tsao cho nay lai lay path[0] ?
	request, err := keeper.GetAIRequest(ctx, path[0])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrRequestNotFound, err.Error())
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResAIRequest(request.RequestID, request.Creator, request.OracleScriptName, request.Validators, request.BlockHeight, request.AIDataSources, request.TestCases, request.Fees.String()))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryFullRequestByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "error")
	}
	// id of the request
	id := path[0]
	request, err := k.GetAIRequest(ctx, id)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRequestNotFound, err.Error())
	}
	// collect all the reports of a given request id
	reports := k.GetReports(ctx, id)

	// collect the result of a given request id

	result, err := k.GetResult(ctx, id)
	if err != nil {
		// init a nil result showing that the result does not have a result yet
		result = types.NewAIRequestResult(id, nil, types.RequestStatusPending)
	}

	res, err := codec.MarshalJSONIndent(k.cdc, types.NewQueryResFullRequest(request, reports, result))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryAIRequestIDs returns all the AI Request IDs in the store
func queryAIRequestIDs(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var requestIDs types.QueryResAIRequestIDs

	iterator := keeper.GetAllAIRequestIDs(ctx)

	for ; iterator.Valid(); iterator.Next() {
		requestIDs = append(requestIDs, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, requestIDs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryTestCase queries an AI request test case
func queryTestCase(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	// tsao cho nay lai lay path[0] ?
	testCase, err := keeper.GetTestCase(ctx, path[0])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrRequestNotFound, err.Error())
	}

	// get code of the test case
	code, err := keeper.fileCache.GetFile(types.TestCaseStoreKeyString(testCase.Name))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCodeNotFound, err.Error())
	}

	executable := base64.StdEncoding.EncodeToString(code)

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResTestCase(testCase.Name, testCase.Owner, executable, testCase.Description))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryTestCases queries a list of test cases
func queryTestCases(ctx sdk.Context, keeper Keeper, req abci.RequestQuery) ([]byte, error) {
	// tsao cho nay lai lay path[0] ?

	var queryResTestCases []types.QueryResTestCase

	// parse limit and offset from the query message data
	pagiSlice := strings.Split(string(req.GetData()[:]), "-")
	page, err := strconv.Atoi(pagiSlice[0])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrPaginationInputInvalid, err.Error())
	}
	limit, err := strconv.Atoi(pagiSlice[1])
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrPaginationInputInvalid, err.Error())
	}

	tCases, err := keeper.GetTestCases(ctx, uint(page), uint(limit))
	if err != nil {
		return []byte{}, sdkerrors.Wrap(types.ErrTestCaseNotFound, err.Error())
	}

	// get code of the each tCase
	for _, tCase := range tCases {
		code, err := keeper.fileCache.GetFile(types.TestCaseStoreKeyString(tCase.Name))
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrCodeNotFound, err.Error())
		}
		executable := base64.StdEncoding.EncodeToString(code)

		// create a new queryResOracleScript
		queryResTestCases = append(queryResTestCases, types.NewQueryResTestCase(tCase.Name, tCase.Owner, executable, tCase.Description))
	}

	// return the query to the command
	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewQueryResTestCases(queryResTestCases, len(tCases)))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryTestCaseNames returns all the test case names in the store
func queryTestCaseNames(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var testCaseNames types.QueryResTestCaseNames

	iterator := keeper.GetAllTestCaseNames(ctx)

	for ; iterator.Valid(); iterator.Next() {
		testCaseNames = append(testCaseNames, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, testCaseNames)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryMinFees(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "error")
	}
	// id of the request
	oScriptName := path[0]
	_, err := k.GetOracleScript(ctx, oScriptName)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, err.Error())
	}

	// collect data source name from the oScript
	oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(oScriptName)

	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.GetDNamesTcNames(oscriptPath)
	if err != nil {
		return nil, err
	}

	minimumFees, err := k.GetMinimumFees(ctx, aiDataSources, testCases)
	if err != nil {
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(k.cdc, types.NewQueryResMinFees(minimumFees.AmountOf(types.Denom).String()))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// GetMinimumFees collects minimum fees needed of an oracle script
func (k Keeper) GetMinimumFees(ctx sdk.Context, dNames, tcNames []string) (sdk.Coins, error) {
	var totalFees sdk.Coins

	// we have different test cases, so we need to loop through them
	for i := 0; i < len(tcNames); i++ {
		// loop to run the test case
		// collect all the test cases object to store in the ai request
		testCase, err := k.GetTestCase(ctx, tcNames[i])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrTestCaseNotFound, fmt.Sprintf("failed to get test case: %s", err.Error()))
		}
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(testCase.Fees...)
	}

	for j := 0; j < len(dNames); j++ {
		// collect all the data source objects to store in the ai request
		aiDataSource, err := k.GetAIDataSource(ctx, dNames[j])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrDataSourceNotFound, fmt.Sprintf("failed to get data source: %s", err.Error()))
		}
		// Aggregate the required fees for an AI request
		totalFees = totalFees.Add(aiDataSource.Fees...)
	}
	rewardRatio := sdk.NewDecWithPrec(int64(k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)), 2)
	minimumFees, _ := sdk.NewDecCoinsFromCoins(totalFees...).QuoDec(rewardRatio).TruncateDecimal()

	return minimumFees, nil
}
