package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// HasReport checks if the report of this ID triple exists in the storage.
func (k *Keeper) HasReport(ctx sdk.Context, id string, val sdk.ValAddress) bool {
	return ctx.KVStore(k.StoreKey).Has(types.ReportStoreKey(id, string(val[:])))
}

// SetReport saves the report to the storage without performing validation.
func (k *Keeper) SetReport(ctx sdk.Context, id string, rep *types.Report) error {
	bz, err := k.Cdc.MarshalBinaryBare(rep)
	if err != nil {
		return err
	}
	ctx.KVStore(k.StoreKey).Set(types.ReportStoreKey(id, string(rep.BaseReport.GetValidatorAddress()[:])), bz)
	return nil
}

// SetTestCaseReport saves the test case report to the storage without performing validation.
func (k *Keeper) SetTestCaseReport(ctx sdk.Context, id string, rep *types.TestCaseReport) error {
	bz, err := k.Cdc.MarshalBinaryBare(rep)
	if err != nil {
		return err
	}
	ctx.KVStore(k.StoreKey).Set(types.TestCaseReportStoreKey(id, string(rep.BaseReport.GetValidatorAddress()[:])), bz)
	return nil
}

// AddReport performs sanity checks and adds a new batch from one validator to one request
// to the store. Note that we expect each validator to report to all raw data requests at once.
func (k *Keeper) AddReport(ctx sdk.Context, rid string, rep *types.Report) error {

	k.SetReport(ctx, rid, rep)
	return nil
}

// GetReportIterator returns the iterator for all reports of the given request ID.
func (k *Keeper) GetReportIterator(ctx sdk.Context, rid string) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.ReportStoreKeyPrefix(rid))
}

// GetReportCount returns the number of reports for the given request ID.
func (k *Keeper) GetReportCount(ctx sdk.Context, rid string) (count uint64) {
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return count
}

// GetReports returns all reports for the given request ID, or nil if there is none.
func (k *Keeper) GetReports(ctx sdk.Context, rid string) (reports []types.Report) {
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &rep)
		reports = append(reports, rep)
	}
	return reports
}

// GetReportsBlockHeight returns all reports for the given block height, or nil if there is none.
func (k *Keeper) GetReportsBlockHeight(ctx sdk.Context) (reports []types.Report) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.ReportStoreKeyPrefixAll())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &rep)
		// check if block height is equal or not
		if rep.BaseReport.GetBlockHeight() == ctx.BlockHeight() {
			reports = append(reports, rep)
		}
	}
	return reports
}

// GetTestCaseReportsBlockHeight returns all reports for the given block height, or nil if there is none.
func (k *Keeper) GetTestCaseReportsBlockHeight(ctx sdk.Context) (reports []types.TestCaseReport) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.TestCaseReportStoreKeyPrefixAll())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.TestCaseReport
		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &rep)
		// check if block height is equal or not
		if rep.BaseReport.GetBlockHeight() == ctx.BlockHeight() {
			reports = append(reports, rep)
		}
	}
	return reports
}

// DeleteReports removes all reports for the given request ID.
func (k *Keeper) DeleteReports(ctx sdk.Context, rid string) {
	var keys [][]byte
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
	}
	for _, key := range keys {
		ctx.KVStore(k.StoreKey).Delete(key)
	}
}

// GetAllReports returns all the reports of every requests ever
func (k *Keeper) GetAllReports(ctx sdk.Context) (reports []types.Report) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.ReportStoreKeyPrefixAll())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &rep)
		reports = append(reports, rep)
	}
	return reports
}

// containsVal returns whether the given slice of validators contains the target validator.
func (k *Keeper) ContainsVal(vals []sdk.ValAddress, target sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(target) {
			return true
		}
	}
	return false
}

func (k *Keeper) ExecuteAIOracles(ctx sdk.Context, valAddress sdk.ValAddress) {
	aiOracles := k.GetAIOraclesBlockHeight(ctx)
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

func isValidator(vals []sdk.ValAddress, valAddr sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(valAddr) {
			return true
		}
	}
	return false
}

func (k *Keeper) executeAIOracle(ctx sdk.Context, aiOracle types.AIOracle, valAddress sdk.ValAddress) {
	// querier to interact with the wasm contract
	querier := NewQuerier(k)
	goCtx := sdk.WrapSDKContext(ctx)
	var dataSourceResults []*types.Result
	var resultArr = []string{}
	// collect list entries to get entry length
	entries, err := querier.DataSourceEntries(goCtx, &types.QueryDataSourceEntriesContract{
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
		result, err := querier.DataSourceContract(goCtx, &types.QueryDataSourceContract{
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

	ctx.Logger().Info(fmt.Sprintf("results collected from the data sources: %v", resultArr))
	aggregatedResult, err := querier.OracleScriptContract(goCtx, &types.QueryOracleScriptContract{
		Contract: aiOracle.GetContract(),
		Request: &types.RequestOracleScript{
			Results: resultArr,
		},
	})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot execute oracle script contract %v with error: %v", aiOracle.GetContract(), err))
	}
	ctx.Logger().Info(fmt.Sprintf("Oracle script final result: %v", aggregatedResult))
	// store report into blockchain as proof for executing AI requests
	report := types.NewReport(aiOracle.GetRequestID(), dataSourceResults, ctx.BlockHeight(), aggregatedResult.Data, valAddress, types.ResultSuccess)
	err = k.SetReport(ctx, aiOracle.GetRequestID(), report)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot store report with request id %v of validator %v with error: %v", aiOracle.GetRequestID(), valAddress.String(), err))
	}
	ctx.Logger().Info(fmt.Sprintf("finish handling the AI oracles with report: %v", report))
}

func (k *Keeper) executeAIOracleTestOnly(ctx sdk.Context, aiOracle types.AIOracle, valAddress sdk.ValAddress) {
	// querier to interact with the wasm contract
	querier := NewQuerier(k)
	goCtx := sdk.WrapSDKContext(ctx)
	var resultsWithTc []*types.ResultWithTestCase
	// collect list entries to get entry length
	entries, err := querier.DataSourceEntries(goCtx, &types.QueryDataSourceEntriesContract{
		Contract: aiOracle.GetContract(),
		Request:  &types.EmptyParams{},
	})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot get data source entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
		return
	}
	tCaseEntries, err := querier.TestCaseEntries(goCtx, &types.QueryTestCaseEntriesContract{
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
			result, err := querier.TestCaseContract(goCtx, &types.QueryTestCaseContract{
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
	report := types.NewTestCaseReport(aiOracle.GetRequestID(), resultsWithTc, ctx.BlockHeight(), valAddress)
	err = k.SetTestCaseReport(ctx, aiOracle.GetRequestID(), report)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Cannot store report with request id %v of validator %v with error: %v", aiOracle.GetRequestID(), valAddress.String(), err))
	}
	ctx.Logger().Info(fmt.Sprintf("finish handling the test AI oracles with report: %v", report))
}
