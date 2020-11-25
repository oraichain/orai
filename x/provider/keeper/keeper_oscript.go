package keeper

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	key := types.OracleScriptStoreKeyString(oldName)
	// if the user does not want to reuse the old name
	if oldName != newName {
		store := ctx.KVStore(k.storeKey)
		byteKey := []byte(key)
		store.Delete(byteKey)
		// delete the old file because it not pointed by any oScript
		k.fileCache.EraseFile(key)
		k.AddOracleScriptFile(code, newName)
	} else {
		k.EditOracleScriptFile(code, oldName)
	}
	k.SetOracleScript(ctx, newName, oScript)
}

// AddOracleScriptFile adds the file to filecache,
func (k Keeper) AddOracleScriptFile(file []byte, name string) {
	k.fileCache.AddFile(file, types.OracleScriptStoreKeyString(name))
}

// EraseOracleScriptFile removes the file in the filecache,
func (k Keeper) EraseOracleScriptFile(name string) {
	k.fileCache.EraseFile(types.OracleScriptStoreKeyString(name))
}

// EditOracleScriptFile edit a file in the filecache,
func (k Keeper) EditOracleScriptFile(file []byte, name string) {
	k.fileCache.EditFile(file, types.OracleScriptStoreKeyString(name))
}

// GetOracleScriptFile gets the oScript code stored in the file storage
func (k Keeper) GetOracleScriptFile(name string) []byte {
	code, err := k.fileCache.GetFile(types.OracleScriptStoreKeyString(name))
	if err != nil {
		return []byte{}
	}
	return code
}

// GetDSourceTCasesScripts is a function that collects test cases and data sources from the oracle script file
func (k Keeper) GetDSourceTCasesScripts(oScript string) ([]string, []string, error) {
	// collect data source name from the oScript script
	oscriptPath := types.ScriptPath + types.OracleScriptStoreKeyString(oScript)
	//use "data source" as an argument to collect the data source script name
	cmd := exec.Command("bash", oscriptPath, "aiDataSource")
	cmd.Stdin = strings.NewReader("some input")
	var dataSourceName bytes.Buffer
	cmd.Stdout = &dataSourceName
	err := cmd.Run()
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrFailedToOpenFile, err.Error())
	}

	// collect data source result from the script
	result := strings.TrimSuffix(dataSourceName.String(), "\n")

	aiDataSources := strings.Fields(result)

	//use "test case" as an argument to collect the test case script name
	cmd = exec.Command("bash", oscriptPath, "testcase")
	cmd.Stdin = strings.NewReader("some input")
	var testCaseName bytes.Buffer
	cmd.Stdout = &testCaseName
	err = cmd.Run()
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrFailedToOpenFile, fmt.Sprintf("failed to collect test case name: %s", result))
	}

	// collect data source result from the script
	result = strings.TrimSuffix(testCaseName.String(), "\n")

	testCases := strings.Fields(result)

	return aiDataSources, testCases, nil
}

// GetOScriptPath is the path for the complete oracle script file location
func (k Keeper) GetOScriptPath(oScriptName string) string {
	// collect data source name from the oScript script
	return types.ScriptPath + types.OracleScriptStoreKeyString(oScriptName)
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
