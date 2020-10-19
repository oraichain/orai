package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ducphamle2/dexai/x/provider/types"
)

// HasReport checks if the report of this ID triple exists in the storage.
func (k Keeper) HasReport(ctx sdk.Context, id string, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.ReportStoreKey(id, string(val[:])))
}

// SetReport saves the report to the storage without performing validation.
func (k Keeper) SetReport(ctx sdk.Context, id string, rep types.Report) {
	ctx.KVStore(k.storeKey).Set(types.ReportStoreKey(id, string(rep.Validator.Address[:])), k.cdc.MustMarshalBinaryBare(rep))
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

func (k Keeper) validateReport(ctx sdk.Context, rid string, rep types.Report, req types.AIRequest) error {
	// Check if the validator is in the requested list of validators
	if !containsVal(req.Validators, rep.Validator.Address) {
		return sdkerrors.Wrap(types.ErrCannotFindValidator, fmt.Sprintln("failed to find the requested validator"))
	}
	// Check if the validator has reported or not
	if k.HasReport(ctx, rid, rep.Validator.Address) {
		return sdkerrors.Wrap(types.ErrValidatorAlreadyReported, fmt.Sprintf("Validator already reported"))
	}
	// if len(rep.RawReports) != len(req.RawRequests) {
	// 	return types.ErrInvalidReportSize
	// }
	// for _, rep := range rep.RawReports {
	// 	// Here we can safely assume that external IDs are unique, as this has already been
	// 	// checked by ValidateBasic performed in baseapp's runTx function.
	// 	if !ContainsEID(req.RawRequests, rep.ExternalID) {
	// 		return sdkerrors.Wrapf(
	// 			types.ErrRawRequestNotFound, "reqID: %d, extID: %d", rid, rep.ExternalID)
	// 	}
	// }
	return nil
}

// AddReport performs sanity checks and adds a new batch from one validator to one request
// to the store. Note that we expect each validator to report to all raw data requests at once.
func (k Keeper) AddReport(ctx sdk.Context, rid string, rep types.Report) error {
	req, err := k.GetAIRequest(ctx, rid)
	if err != nil {
		return err
	}
	// We need to validate the report before storing into the blockchain
	err = k.validateReport(ctx, rid, rep, req)
	if err != nil {
		return err
	}
	k.SetReport(ctx, rid, rep)
	return nil
}

// GetReportIterator returns the iterator for all reports of the given request ID.
func (k Keeper) GetReportIterator(ctx sdk.Context, rid string) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.ReportStoreKeyPrefix(rid))
}

// GetReportCount returns the number of reports for the given request ID.
func (k Keeper) GetReportCount(ctx sdk.Context, rid string) (count uint64) {
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return count
}

// GetReports returns all reports for the given request ID, or nil if there is none.
func (k Keeper) GetReports(ctx sdk.Context, rid string) (reports []types.Report) {
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rep)
		reports = append(reports, rep)
	}
	return reports
}

// GetReportsBlockHeight returns all reports for the given block height, or nil if there is none.
func (k Keeper) GetReportsBlockHeight(ctx sdk.Context, blockHeight int64) (reports []types.Report) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.ReportStoreKeyPrefixAll())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rep)
		// check if block height is equal or not
		if rep.BlockHeight == blockHeight {
			reports = append(reports, rep)
		}
	}
	return reports
}

// DeleteReports removes all reports for the given request ID.
func (k Keeper) DeleteReports(ctx sdk.Context, rid string) {
	var keys [][]byte
	iterator := k.GetReportIterator(ctx, rid)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
	}
	for _, key := range keys {
		ctx.KVStore(k.storeKey).Delete(key)
	}
}

// GetAllReports returns all the reports of every requests ever
func (k Keeper) GetAllReports(ctx sdk.Context) (reports []types.Report) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.ReportStoreKeyPrefixAll())

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rep types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rep)
		reports = append(reports, rep)
	}
	return reports
}
