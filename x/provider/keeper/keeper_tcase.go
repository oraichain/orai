package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

// GetAllTestCaseNames get an iterator of all key-value pairs in the store
func (k Keeper) GetAllTestCaseNames(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.TestCaseKeyPrefix))
}

// GetPaginatedTestCaseNames get an iterator of paginated key-value pairs in the store
func (k Keeper) GetPaginatedTestCaseNames(ctx sdk.Context, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIteratorPaginated(store, []byte(types.TestCaseKeyPrefix), page, limit)
}

// GetTestCase returns the the AI test case of a given request
func (k Keeper) GetTestCase(ctx sdk.Context, name string) (types.TestCase, error) {
	store := ctx.KVStore(k.storeKey)
	var testCase types.TestCase
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.TestCaseStoreKey(name)), &testCase)
	if err != nil {
		return types.TestCase{}, err
	}
	return testCase, nil
}

// GetTestCases returns list of test cases
func (k Keeper) GetTestCases(ctx sdk.Context, page, limit uint) ([]types.TestCase, error) {
	var tCases []types.TestCase

	iterator := k.GetPaginatedTestCaseNames(ctx, page, limit)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tCase types.TestCase
		err := k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &tCase)
		if err != nil {
			return []types.TestCase{}, err
		}
		tCases = append(tCases, tCase)
	}
	return tCases, nil
}

// SetTestCase allows users to set a test case into the store
func (k Keeper) SetTestCase(ctx sdk.Context, name string, testCase types.TestCase) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(testCase)
	store.Set(types.TestCaseStoreKey(name), bz)
}

// DefaultTestCase creates an empty Test Case struct
func (k Keeper) DefaultTestCase() types.TestCase {
	return types.TestCase{}
}

// EditTestCase allows users to edit a test case in the store
func (k Keeper) EditTestCase(ctx sdk.Context, oldName, newName string, code []byte, testCase types.TestCase) {
	key := types.TestCaseStoreKey(oldName)
	// if the user does not want to reuse the old name
	if oldName != newName {
		store := ctx.KVStore(k.storeKey)
		byteKey := []byte(key)
		store.Delete(byteKey)
		// delete the old file because it not pointed by any oScript
		k.EraseTestCaseFile(oldName)
		k.AddTestCaseFile(code, newName)
	} else {
		// edit the file instead since old name = new name
		k.EditTestCaseFile(code, oldName)
	}
	k.SetTestCase(ctx, newName, testCase)
}

// GetTestCaseFile gets the test case code stored in the file storage
func (k Keeper) GetTestCaseFile(name string) ([]byte, error) {
	code, err := k.fileCache.GetFile(types.TestCaseStoreFileString(name))
	if err != nil {
		return []byte{}, err
	}
	return code, nil
}

// AddTestCaseFile adds the file to filecache,
func (k Keeper) AddTestCaseFile(file []byte, name string) {
	k.fileCache.AddFile(file, types.TestCaseStoreFileString(name))
}

// EditTestCaseFile edit a file in the filecache,
func (k Keeper) EditTestCaseFile(file []byte, name string) {
	k.fileCache.EditFile(file, types.TestCaseStoreFileString(name))
}

// EraseTestCaseFile removes the file in the filecache,
func (k Keeper) EraseTestCaseFile(name string) {
	k.fileCache.EraseFile(types.TestCaseStoreFileString(name))
}
