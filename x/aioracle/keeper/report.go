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

// AddReport performs sanity checks and adds a new batch from one validator to one request
// to the store. Note that we expect each validator to report to all raw data requests at once.
func (k *Keeper) AddReport(ctx sdk.Context, rid string, rep *types.Report) error {
	return k.SetReport(ctx, rid, rep)
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

// DeleteReports removes all reports for the given request ID.
func (k *Keeper) DeleteReports(ctx sdk.Context, rid string) {
	store := ctx.KVStore(k.StoreKey)
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// GetAllReports returns all the reports of every requests ever
func (k *Keeper) GetAllReports(ctx sdk.Context) (reports []types.Report) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.StoreKey), types.ReportStoreKeyPrefixAll())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		err := k.Cdc.UnmarshalBinaryBare(iterator.Value(), &rep)
		if err == nil {
			reports = append(reports, rep)
		}
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

func isValidator(vals []sdk.ValAddress, valAddr sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(valAddr) {
			return true
		}
	}
	return false
}

func (k *Keeper) validateBasic(ctx sdk.Context, req *types.AIOracle, rep *types.Report, blockHeight int64) bool {
	// Check if validator exists and active
	_, isExist := k.StakingKeeper.GetValidator(ctx, rep.BaseReport.GetValidatorAddress())
	if !isExist {
		k.Logger(ctx).Error(fmt.Sprintf("error in validating the report: validator does not exist"))
		return false
	}
	if !k.ContainsVal(req.GetValidators(), rep.BaseReport.GetValidatorAddress()) {
		k.Logger(ctx).Error(fmt.Sprintf("Validator %v does not exist in the list of request validators", rep.BaseReport.GetValidatorAddress().String()))
		return false
	}
	if len(rep.BaseReport.GetValidatorAddress()) == 0 || len(rep.BaseReport.GetRequestId()) == 0 || rep.BaseReport.GetBlockHeight() <= 0 {
		k.Logger(ctx).Error(fmt.Sprintf("Report basic information is invalid: %v", rep))
		return false
	}
	return true
}

func (k *Keeper) validateReportBasic(ctx sdk.Context, req *types.AIOracle, rep *types.Report, blockHeight int64) bool {
	if rep.ResultStatus == types.ResultFailure {
		k.Logger(ctx).Error("result status is fail")
		return false
	}
	if len(rep.GetDataSourceResults()) == 0 || len(rep.GetAggregatedResult()) == 0 {
		k.Logger(ctx).Error(fmt.Sprintf("Report results are invalid: %v", rep))
		return false
	}
	if rep.GetResultStatus() != types.ResultFailure && rep.GetResultStatus() != types.ResultSuccess {
		k.Logger(ctx).Error(fmt.Sprintf("Report result status is invalid: %v", rep.GetResultStatus()))
		return false
	}
	var dsResultSize int
	for _, dsResult := range rep.GetDataSourceResults() {
		if dsResult.GetStatus() != types.ResultFailure && dsResult.GetStatus() != types.ResultSuccess {
			k.Logger(ctx).Error(fmt.Sprintf("Data source result status is invalid: %v", dsResult.GetStatus()))
			return false
		}
		dsResultSize += len(dsResult.Result)
	}
	aggregatedResultSize := len(rep.GetAggregatedResult())
	finalLen := dsResultSize + aggregatedResultSize
	responseBytes := k.GetParam(ctx, types.KeyMaximumAIOracleResBytes)

	if finalLen >= int(responseBytes) {
		k.Logger(ctx).Error(fmt.Sprintf("Report result size: %v cannot be larger than %v", finalLen, responseBytes))
		return false
	}

	return k.validateBasic(ctx, req, rep, blockHeight)
}

func (k *Keeper) validateTestCaseReportBasic(ctx sdk.Context, req *types.AIOracle, rep *types.TestCaseReport, blockHeight int64) bool {
	if len(rep.GetResultsWithTestCase()) <= 0 {
		k.Logger(ctx).Error(fmt.Sprintf("Report results are invalid: %v", rep))
		return false
	}
	var tcResultSize int
	for _, result := range rep.GetResultsWithTestCase() {
		for _, tcResult := range result.GetTestCaseResults() {
			if tcResult.GetStatus() != types.ResultFailure && tcResult.GetStatus() != types.ResultSuccess {
				k.Logger(ctx).Error(fmt.Sprintf("Report result status is invalid: %v", tcResult.GetStatus()))
				return false
			}
			tcResultSize += len(tcResult.GetResult())
		}
	}
	responseBytes := k.GetParam(ctx, types.KeyMaximumAIOracleResBytes)

	if tcResultSize >= int(responseBytes) {
		k.Logger(ctx).Error(fmt.Sprintf("Report result size: %v cannot be larger than %v", tcResultSize, int(responseBytes)))
		return false
	}
	return true
}
