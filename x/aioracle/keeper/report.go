package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/aioracle/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
	ctx.KVStore(k.StoreKey).Set(types.ReportStoreKey(id, string(rep.Reporter.Validator[:])), bz)
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
		if rep.GetBlockHeight() == ctx.BlockHeight() {
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

// ValidateReport validates if the report is valid to get rewards
func (k *Keeper) ValidateReport(ctx sdk.Context, reporter *types.Reporter, req *types.AIOracle) error {
	// Check if the validator is in the requested list of validators
	if !containsVal(req.GetValidators(), reporter.GetValidator()) {
		return sdkerrors.Wrap(types.ErrValidatorNotFound, fmt.Sprintln("failed to find the requested validator"))
	}
	return nil
}

// containsVal returns whether the given slice of validators contains the target validator.
func containsVal(vals []sdk.ValAddress, target sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(target) {
			return true
		}
	}
	return false
}

func (k *Keeper) ExecuteAIOracles(ctx sdk.Context, req abci.RequestBeginBlock) {
	aiOracles := k.GetAIOraclesBlockHeight(ctx)
	if len(aiOracles) == 0 {
		return
	}

	//clientCtx := client.Context{}

	clientCtx := ctx.Value(client.ClientContextKey).(*client.Context)
	fmt.Println("client ctx: ", clientCtx.GetFromAddress())

	// querier to interact with the wasm contract
	querier := NewQuerier(k)
	goCtx := sdk.WrapSDKContext(ctx)

	// execute each ai oracle request
	for _, aiOracle := range aiOracles {
		var dataSourceResults []*types.DataSourceResult
		var resultArr = []string{}
		// collect list entries to get entry length
		entries, err := querier.DataSourceEntries(goCtx, &types.QueryDataSourceEntriesContract{
			Contract: aiOracle.GetContract(),
			Request:  &types.EmptyParams{},
		})
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Cannot get data source entries for the given request contract: %v with error: %v", aiOracle.GetContract(), err))
			continue
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
			dataSourceResult := types.NewDataSourceResult(entry, result.GetData(), types.ResultSuccess)
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
		// report := types.NewReport(aiOracle.GetRequestID(), dataSourceResults, nil, ctx.BlockHeight(), aggregatedResult.Data, types.NewReporter())
	}
}
