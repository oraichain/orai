package keeper

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/oraichain/orai/packages/filecache"
	"github.com/oraichain/orai/x/provider/types"
)

// Keeper of the provider store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       params.Subspace
	supplyKeeper     types.SupplyKeeper
	bankKeeper       types.BankKeeper
	stakingKeeper    types.StakingKeeper
	DistrKeeper      types.DistrKeeper
	fileCache        filecache.Cache
	feeCollectorName string
}

// NewKeeper creates a provider keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, providerSubspace params.Subspace, supplyKeeper types.SupplyKeeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, distrKeeper types.DistrKeeper, feeCollectorName string, cacheDir string) Keeper {
	if !providerSubspace.HasKeyTable() {
		// register parameters of the provider module into the param space
		providerSubspace = providerSubspace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       providerSubspace,
		fileCache:        filecache.New(cacheDir),
		supplyKeeper:     supplyKeeper,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		DistrKeeper:      distrKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

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

// GetAllAIDataSourceNames get an iterator of all key-value pairs in the store
func (k Keeper) GetAllAIDataSourceNames(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.DataSourceKeyPrefix))
}

// GetPaginatedAIDataSourceNames get an iterator of paginated key-value pairs in the store
func (k Keeper) GetPaginatedAIDataSourceNames(ctx sdk.Context, page, limit uint) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIteratorPaginated(store, []byte(types.DataSourceKeyPrefix), page, limit)
}

// SetOracleScript allows users to set a oScript into the store
func (k Keeper) SetOracleScript(ctx sdk.Context, name string, oScript types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(oScript)
	store.Set(types.OracleScriptStoreKey(name), bz)
}

// EditOracleScript allows users to edit a oScript in the store
func (k Keeper) EditOracleScript(ctx sdk.Context, oldName string, newName string, oScript types.OracleScript) {

	key := types.OracleScriptStoreKeyString(oldName)
	// if the user does not want to reuse the old name
	if oldName != newName {
		store := ctx.KVStore(k.storeKey)
		byteKey := []byte(key)
		store.Delete(byteKey)
		// delete the old file because it not pointed by any oScript
		k.fileCache.EraseFile(key)
	}

	k.SetOracleScript(ctx, newName, oScript)
}

// GetAIDataSource returns the data source object given the name of the data source
func (k Keeper) GetAIDataSource(ctx sdk.Context, name string) (types.AIDataSource, error) {
	store := ctx.KVStore(k.storeKey)
	var aiDataSource types.AIDataSource
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.DataSourceStoreKey(name)), &aiDataSource)
	if err != nil {
		return types.AIDataSource{}, err
	}
	return aiDataSource, nil
}

// GetAIDataSources returns list of data sources
func (k Keeper) GetAIDataSources(ctx sdk.Context, page, limit uint) ([]types.AIDataSource, error) {
	var dSources []types.AIDataSource

	iterator := k.GetPaginatedAIDataSourceNames(ctx, page, limit)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dSource types.AIDataSource
		err := k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &dSource)
		if err != nil {
			return []types.AIDataSource{}, err
		}
		dSources = append(dSources, dSource)
	}
	return dSources, nil
}

//IsNamePresent checks if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// SetAIDataSource allows users to set a data source into the store
func (k Keeper) SetAIDataSource(ctx sdk.Context, name string, aiDataSource types.AIDataSource) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(aiDataSource)
	store.Set(types.DataSourceStoreKey(name), bz)
}

// EditAIDataSource allows users to edit a data source in the store
func (k Keeper) EditAIDataSource(ctx sdk.Context, oldName string, newName string, aiDataSource types.AIDataSource) {
	key := types.DataSourceStoreKeyString(oldName)
	// if the user does not want to reuse the old name
	if oldName != newName {
		store := ctx.KVStore(k.storeKey)
		byteKey := []byte(key)
		store.Delete(byteKey)
		// delete the old file because it not pointed by any oScript
		k.fileCache.EraseFile(key)
	}

	k.SetAIDataSource(ctx, newName, aiDataSource)
}

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

// CreateTestCase allows users to set a test case into the store
func (k Keeper) CreateTestCase(ctx sdk.Context, name string, testCase types.TestCase) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(testCase)
	store.Set(types.TestCaseStoreKey(name), bz)
}

// EditTestCase allows users to edit a test case in the store
func (k Keeper) EditTestCase(ctx sdk.Context, oldName string, newName string, testCase types.TestCase) {
	key := types.TestCaseStoreKeyString(oldName)
	// if the user does not want to reuse the old name
	if oldName != newName {
		store := ctx.KVStore(k.storeKey)
		byteKey := []byte(key)
		store.Delete(byteKey)
		// delete the old file because it not pointed by any oScript
		k.fileCache.EraseFile(key)
	}

	k.CreateTestCase(ctx, newName, testCase)
}

// CreateStrategy allows users to create a new strategy into the store
func (k Keeper) CreateStrategy(ctx sdk.Context, name string, strategy types.Strategy) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(strategy)
	store.Set(types.StrategyStoreKey(strategy.StratID, name), bz)
}

// GetParam returns the parameter as specified by key as an uint64.
func (k Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
	k.paramSpace.Get(ctx, key, &res)
	return res
}

// SetParam saves the given key-value parameter to the store.
func (k Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
	k.paramSpace.Set(ctx, key, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// AddOracleScriptFile adds the file to filecache,
func (k Keeper) AddOracleScriptFile(file []byte, name string) {
	k.fileCache.AddFile(file, types.OracleScriptStoreKeyString(name))
}

// EraseOracleScriptFile removes the file in the filecache,
func (k Keeper) EraseOracleScriptFile(name string) {
	k.fileCache.EraseFile(types.OracleScriptStoreKeyString(name))
}

// AddAIDataSourceFile adds the file to filecache,
func (k Keeper) AddAIDataSourceFile(file []byte, name string) {
	k.fileCache.AddFile(file, types.DataSourceStoreKeyString(name))
}

// AddTestCaseFile adds the file to filecache,
func (k Keeper) AddTestCaseFile(file []byte, name string) {
	k.fileCache.AddFile(file, types.TestCaseStoreKeyString(name))
}

// EditOracleScriptFile edit a file in the filecache,
func (k Keeper) EditOracleScriptFile(file []byte, name string) {
	k.fileCache.EditFile(file, types.OracleScriptStoreKeyString(name))
}

// EditAIDataSourceFile edit a file in the filecache,
func (k Keeper) EditAIDataSourceFile(file []byte, name string) {
	k.fileCache.EditFile(file, types.DataSourceStoreKeyString(name))
}

// EditTestCaseFile edit a file in the filecache,
func (k Keeper) EditTestCaseFile(file []byte, name string) {
	k.fileCache.EditFile(file, types.TestCaseStoreKeyString(name))
}

// GetOracleScriptFile gets the oScript code stored in the file storage
func (k Keeper) GetOracleScriptFile(name string) []byte {
	code, err := k.fileCache.GetFile(types.OracleScriptStoreKeyString(name))
	if err != nil {
		return []byte{}
	}
	return code
}

// GetAIDataSourceFile gets the data source code stored in the file storage
func (k Keeper) GetAIDataSourceFile(name string) []byte {
	code, err := k.fileCache.GetFile(types.DataSourceStoreKeyString(name))
	if err != nil {
		return []byte{}
	}
	return code
}

// GetDNamesTcNames is a function that collects test case names and data source names from the oracle script
func (k Keeper) GetDNamesTcNames(oscriptPath string) ([]string, []string, error) {
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
