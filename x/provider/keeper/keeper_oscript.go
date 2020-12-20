package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

// GetOracleScript returns the oScript object given the name of the oScript
func (k Keeper) GetOracleScript(ctx sdk.Context, name string) (types.OracleScript, error) {
	store := ctx.KVStore(k.storeKey)
	var oScript types.OracleScript
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.OracleScriptStoreKey(name)), &oScript)
	if err != nil {
		return types.OracleScript{}, err
	}
	return oScript, nil
}

// SetOracleScript allows users to set a oScript into the store
func (k Keeper) SetOracleScript(ctx sdk.Context, name string, oScript types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(oScript)
	store.Set(types.OracleScriptStoreKey(name), bz)
}

// GetOracleScripts returns list of oracle scripts
func (k Keeper) GetOracleScripts(ctx sdk.Context, page, limit uint) ([]types.OracleScript, error) {
	var oScripts []types.OracleScript

	iterator := k.GetPaginatedOracleScriptNames(ctx, page, limit)
	//iterator := k.GetAllOracleScriptNames(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oScript types.OracleScript
		err := k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &oScript)
		if err != nil {
			return []types.OracleScript{}, err
		}
		oScripts = append(oScripts, oScript)
	}
	return oScripts, nil
}

// GetAllOracleScriptNames get an iterator of all key-value pairs in the store
func (k Keeper) GetAllOracleScriptNames(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.OScriptKeyPrefix))
}

// GetPaginatedOracleScriptNames get an iterator of paginated key-value pairs in the store
func (k Keeper) GetPaginatedOracleScriptNames(ctx sdk.Context, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIteratorPaginated(store, []byte(types.OScriptKeyPrefix), page, limit)
}

// EditOracleScript allows users to edit a oScript in the store
func (k Keeper) EditOracleScript(ctx sdk.Context, oldName, newName string, code []byte, oScript types.OracleScript) {

	key := types.OracleScriptStoreKey(oldName)
	// if the user does not want to reuse the old name
	if oldName != newName {
		store := ctx.KVStore(k.storeKey)
		byteKey := []byte(key)
		store.Delete(byteKey)
		// delete the old file because it not pointed by any oScript
		k.EraseOracleScriptFile(oldName)
		k.AddOracleScriptFile(code, newName)
	} else {
		k.EditOracleScriptFile(code, oldName)
	}
	k.SetOracleScript(ctx, newName, oScript)
}

// AddOracleScriptFile adds the file to filecache,
func (k Keeper) AddOracleScriptFile(file []byte, name string) {
	k.fileCache.AddFile(file, types.OracleScriptFileString(name))
}

// EraseOracleScriptFile removes the file in the filecache,
func (k Keeper) EraseOracleScriptFile(name string) {
	k.fileCache.EraseFile(types.OracleScriptFileString(name))
}

// EditOracleScriptFile edit a file in the filecache,
func (k Keeper) EditOracleScriptFile(file []byte, name string) {
	k.fileCache.EditFile(file, types.OracleScriptFileString(name))
}

// GetOracleScriptFile gets the oScript code stored in the file storage
func (k Keeper) GetOracleScriptFile(name string) ([]byte, error) {
	code, err := k.fileCache.GetFile(types.OracleScriptFileString(name))
	if err != nil {
		return []byte{}, err
	}
	return code, nil
}

// GetOScriptPath is the path for the complete oracle script file location
func (k Keeper) GetOScriptPath(oScriptName string) string {
	// collect data source name from the oScript script
	return types.ScriptPath + types.OracleScriptFileString(oScriptName)
}

// GetDNamesTcNames - an utility function for retriving data source and test case names from the oracle script
func (k Keeper) GetDNamesTcNames(ctx sdk.Context, oScript string) ([]string, []string, error) {
	// get data source and test case names from the oracle script
	oracleScript, err := k.GetOracleScript(ctx, oScript)
	if err != nil {
		return nil, nil, err
	}
	aiDataSources := oracleScript.GetDSources()
	testCases := oracleScript.GetTCases()
	return aiDataSources, testCases, nil
}

// GetKeyOracleScriptRewardPercentage returns the oracle script reward percentage from the provider module
func (k Keeper) GetKeyOracleScriptRewardPercentage(ctx sdk.Context) int64 {
	percentage := k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)
	return int64(percentage)
}
