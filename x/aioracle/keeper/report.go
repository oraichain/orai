package keeper

import (
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
