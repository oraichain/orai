package keeper

import (
	"context"
	"encoding/json"
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

// containsVal returns whether the given slice of validators contains the target validator.
func (k *Querier) ContainsVal(vals []sdk.ValAddress, target sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(target) {
			return true
		}
	}
	return false
}

func (k *Querier) ExecuteAIOracles(ctx sdk.Context, valAddress sdk.ValAddress) {
	aiOracles := k.keeper.GetAIOraclesBlockHeight(ctx)
	if len(aiOracles) == 0 {
		return
	}
	// execute each ai oracle request
	for _, aiOracle := range aiOracles {
		// if the ai oracle request does not include the validator address, then we skip
		if !isValidator(aiOracle.GetValidators(), valAddress) {
			continue
		}
		if aiOracle.TestOnly {
			k.executeAIOracleTestOnly(ctx, aiOracle, valAddress)
			continue
		}
		k.executeAIOracle(ctx, aiOracle, valAddress)
	}
}

func (k *Querier) executeAIOracle(ctx sdk.Context, aiOracle types.AIOracle, valAddress sdk.ValAddress) {
	// querier to interact with the wasm contract
	goCtx := sdk.WrapSDKContext(ctx)
	var dataSourceResults []*types.Result
	var resultArr = []string{}
	// collect list entries to get entry length
	entries, err := k.DataSourceEntries(goCtx, &types.QueryDataSourceEntriesContract{
		Contract: aiOracle.GetContract(),
		Request:  &types.EmptyParams{},
	})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot get data source entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
		return
	}

	// if there's no data source we stop executing and return no report
	if len(entries.GetData()) <= 0 {
		ctx.Logger().Error(fmt.Sprintf("The data source entry list is empty"))
		return
	}

	// loop to execute data source one by one
	for _, entry := range entries.GetData() {
		// run the data source script
		ctx.Logger().Info(fmt.Sprintf("Data source entrypoint: %v and input: %v", entry, string(aiOracle.GetInput())))
		result, err := k.DataSourceContract(goCtx, &types.QueryDataSourceContract{
			Contract: aiOracle.GetContract(),
			Request: &types.RequestDataSource{
				Dsource: entry,
				Input:   string(aiOracle.GetInput()),
			},
		})
		// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
		dataSourceResult := types.NewResult(entry, result.GetData(), types.ResultSuccess)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Cannot execute given data source %v with error: %v", entry.Url, err))
			// change status to fail so the datasource cannot be rewarded afterwards
			dataSourceResult.Status = types.ResultFailure
			dataSourceResult.Result = []byte(types.FailedResponseDs)
			continue
		}
		resultArr = append(resultArr, string(result.Data))
		// append an data source result into the list
		dataSourceResults = append(dataSourceResults, dataSourceResult)
	}

	// store report into blockchain as proof for executing AI requests
	report := types.NewReport(aiOracle.GetRequestID(), dataSourceResults, ctx.BlockHeight(), []byte{}, valAddress, types.ResultSuccess, sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(0))))

	ctx.Logger().Info(fmt.Sprintf("results collected from the data sources: %v", resultArr))
	aggregatedResult, err := k.OracleScriptContract(goCtx, &types.QueryOracleScriptContract{
		Contract: aiOracle.GetContract(),
		Request: &types.RequestOracleScript{
			Results: resultArr,
		},
	})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot execute oracle script contract %v with error: %v", aiOracle.GetContract(), err))
		report.AggregatedResult = []byte(types.FailedResponseOs)
		report.ResultStatus = types.ResultFailure
	} else {
		report.AggregatedResult = aggregatedResult.Data
	}
	ctx.Logger().Info(fmt.Sprintf("Oracle script final result: %v", aggregatedResult))
	reportBytes, err := json.Marshal(report)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot marshal report %v with error: %v", report, err))
	}
	// verify basic report information to make sure the report is stored in a valid manner
	if k.validateReportBasic(ctx, &aiOracle, report, report.BaseReport.GetBlockHeight()) {
		event := sdk.NewEvent(types.EventTypeReportWithData)
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeReport, string(reportBytes[:])),
		)
		ctx.EventManager().Events().AppendEvent(event)
		ctx.EventManager().EmitEvent(event)
		ctx.Logger().Info(fmt.Sprintf("Finish handling the AI oracles with report: %v", report))
	}
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

func (k *Querier) validateReportBasic(ctx sdk.Context, req *types.AIOracle, rep *types.Report, blockHeight int64) bool {
	if len(rep.GetDataSourceResults()) <= 0 || len(rep.GetAggregatedResult()) <= 0 {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Report results are invalid: %v", rep))
		return false
	}
	if rep.GetResultStatus() != types.ResultFailure && rep.GetResultStatus() != types.ResultSuccess {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Report result status is invalid: %v", rep.GetResultStatus()))
		return false
	}
	var dsResultSize int
	for _, dsResult := range rep.GetDataSourceResults() {
		if dsResult.GetStatus() != types.ResultFailure && dsResult.GetStatus() != types.ResultSuccess {
			k.keeper.Logger(ctx).Error(fmt.Sprintf("Data source result status is invalid: %v", dsResult.GetStatus()))
			return false
		}
		dsResultSize += len(dsResult.Result)
	}
	aggregatedResultSize := len(rep.GetAggregatedResult())
	finalLen := dsResultSize + aggregatedResultSize
	responseBytes := k.keeper.GetParam(ctx, types.KeyMaximumAIOracleResBytes)

	if finalLen >= int(responseBytes) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Report result size: %v cannot be larger than %v", finalLen, responseBytes))
		return false
	}

	return k.validateBasic(ctx, req, rep.BaseReport)
}

func (k *Querier) executeAIOracleTestOnly(ctx sdk.Context, aiOracle types.AIOracle, valAddress sdk.ValAddress) {
	// querier to interact with the wasm contract
	goCtx := sdk.WrapSDKContext(ctx)
	var resultsWithTc []*types.ResultWithTestCase
	// collect list entries to get entry length
	entries, err := k.DataSourceEntries(goCtx, &types.QueryDataSourceEntriesContract{
		Contract: aiOracle.GetContract(),
		Request:  &types.EmptyParams{},
	})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot get data source entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
		return
	}
	tCaseEntries, err := k.TestCaseEntries(goCtx, &types.QueryTestCaseEntriesContract{
		Contract: aiOracle.GetContract(),
		Request:  &types.EmptyParams{},
	})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot get test case entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
		return
	}

	// if there's no data source or test case then we stop executing
	if len(entries.GetData()) <= 0 || len(tCaseEntries.GetData()) <= 0 {
		ctx.Logger().Error(fmt.Sprintf("The data source or test case entry list is empty"))
		return
	}

	// loop to execute data source one by one
	for _, entry := range entries.GetData() {
		var results []*types.Result
		for _, tCaseEntry := range tCaseEntries.GetData() {
			// run the data source script
			ctx.Logger().Info(fmt.Sprintf("Data source entrypoint: %v and input: %v", entry, string(aiOracle.GetInput())))
			ctx.Logger().Info(fmt.Sprintf("Testcase entrypoint: %v", tCaseEntry))
			result, err := k.TestCaseContract(goCtx, &types.QueryTestCaseContract{
				Contract: aiOracle.GetContract(),
				Request: &types.RequestTestCase{
					Tcase: tCaseEntry,
					Input: entry,
				},
			})
			tCaseResult := types.NewResult(tCaseEntry, []byte{}, types.ResultSuccess)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("Cannot execute test case %v with error: %v", tCaseEntry.Url, err))
				tCaseResult.Result = []byte(types.FailedResponseTc)
				tCaseResult.Status = types.ResultFailure
				results = append(results, tCaseResult)
				continue
			}
			tCaseResult.Result = result.GetData()
			results = append(results, tCaseResult)
		}
		resultWithTc := types.NewResultWithTestCase(entry, results, types.ResultSuccess)
		resultsWithTc = append(resultsWithTc, resultWithTc)
	}
	// store report into blockchain as proof for executing AI requests
	report := types.NewTestCaseReport(aiOracle.GetRequestID(), resultsWithTc, ctx.BlockHeight(), valAddress, sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.NewInt(0))))
	reportBytes, err := json.Marshal(report)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot marshal report %v with error: %v", report, err))
	}
	if k.validateTestCaseReportBasic(ctx, &aiOracle, report, report.BaseReport.GetBlockHeight()) {
		event := sdk.NewEvent(types.EventTypeReportWithData)
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeTestCaseReport, string(reportBytes[:])),
		)
		ctx.EventManager().Events().AppendEvent(event)
		ctx.EventManager().EmitEvent(event)
		ctx.Logger().Info(fmt.Sprintf("Finish handling the test AI oracles with report: %v", report))
	}
}

func (k *Querier) validateTestCaseReportBasic(ctx sdk.Context, req *types.AIOracle, rep *types.TestCaseReport, blockHeight int64) bool {
	if len(rep.GetResultsWithTestCase()) <= 0 {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Report results are invalid: %v", rep))
		return false
	}
	var tcResultSize int
	for _, result := range rep.GetResultsWithTestCase() {
		for _, tcResult := range result.GetTestCaseResults() {
			if tcResult.GetStatus() != types.ResultFailure && tcResult.GetStatus() != types.ResultSuccess {
				k.keeper.Logger(ctx).Error(fmt.Sprintf("Report result status is invalid: %v", tcResult.GetStatus()))
				return false
			}
			tcResultSize += len(tcResult.GetResult())
		}
	}
	responseBytes := k.keeper.GetParam(ctx, types.KeyMaximumAIOracleResBytes)

	if tcResultSize >= int(responseBytes) {
		k.keeper.Logger(ctx).Error(fmt.Sprintf("Report result size: %v cannot be larger than %v", tcResultSize, int(responseBytes)))
		return false
	}
	return k.validateBasic(ctx, req, rep.BaseReport)
}

func isValidator(vals []sdk.ValAddress, valAddr sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(valAddr) {
			return true
		}
	}
	return false
}
