package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// HasReport checks if the report of this ID triple exists in the storage.
func (k *Keeper) HasReport(ctx sdk.Context, id string, val sdk.ValAddress) bool {
	return ctx.KVStore(k.StoreKey).Has(types.ReportStoreKey(id, string(val)))
}

// HasTestCaseReport checks if the test case report of this ID triple exists in the storage.
func (k *Keeper) HasTestCaseReport(ctx sdk.Context, id string, val sdk.ValAddress) bool {
	return ctx.KVStore(k.StoreKey).Has(types.TestCaseReportStoreKey(id, string(val[:])))
}

// SetReport saves the report to the storage without performing validation.
func (k *Keeper) SetReport(ctx sdk.Context, id string, rep *types.Report) error {
	bz, err := k.Cdc.MarshalBinaryBare(rep)
	if err != nil {
		return err
	}
	ctx.KVStore(k.StoreKey).Set(types.ReportStoreKey(id, string(rep.BaseReport.GetValidatorAddress())), bz)
	return nil
}

// SetTestCaseReport saves the test case report to the storage without performing validation.
func (k *Keeper) SetTestCaseReport(ctx sdk.Context, id string, rep *types.TestCaseReport) error {
	bz, err := k.Cdc.MarshalBinaryBare(rep)
	if err != nil {
		return err
	}
	ctx.KVStore(k.StoreKey).Set(types.TestCaseReportStoreKey(id, string(rep.BaseReport.GetValidatorAddress())), bz)
	return nil
}

// GetReportIterator returns the iterator for all reports of the given request ID.
func (k *Keeper) GetReportIterator(ctx sdk.Context, rid string) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.ReportStoreKeyPrefix(rid))
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

// GetTestCaseReports returns test case all reports for the given request ID, or nil if there is none.
func (k *Keeper) GetTestCaseReports(ctx sdk.Context, rid string) (reports []types.TestCaseReport) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.TestCaseReportStoreKeyPrefix(rid))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.TestCaseReport
		err := k.Cdc.UnmarshalBinaryBare(iterator.Value(), &rep)
		// incase panic or nil value return
		if err == nil {
			reports = append(reports, rep)
		}
	}
	return reports
}

// GetReportsBlockHeight returns all reports for the given block height, or nil if there is none.
func (k *Keeper) GetReportsBlockHeight(ctx sdk.Context) (reports []types.Report) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.ReportStoreKeyPrefixAll())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		err := k.Cdc.UnmarshalBinaryBare(iterator.Value(), &rep)
		// check if block height is equal or not
		if err == nil && rep.BaseReport.GetBlockHeight() == ctx.BlockHeight() {
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
		err := k.Cdc.UnmarshalBinaryBare(iterator.Value(), &rep)
		// check if block height is equal or not
		if err == nil && rep.BaseReport.GetBlockHeight() == ctx.BlockHeight() {
			reports = append(reports, rep)
		}
	}
	return reports
}

func (k *Querier) validateReportBasic(ctx sdk.Context, req *types.AIOracle, rep *types.Report) bool {
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

func (k *Querier) validateTestCaseReportBasic(ctx sdk.Context, req *types.AIOracle, rep *types.TestCaseReport) bool {
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
