package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

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

// DefaultAIDataSource creates an empty Data Source struct
func (k Keeper) DefaultAIDataSource() types.AIDataSource {
	return types.AIDataSource{}
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

// SetAIDataSource allows users to set a data source into the store
func (k Keeper) SetAIDataSource(ctx sdk.Context, name string, aiDataSource types.AIDataSource) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(aiDataSource)
	store.Set(types.DataSourceStoreKey(name), bz)
}

// EditAIDataSource allows users to edit a data source in the store
func (k Keeper) EditAIDataSource(ctx sdk.Context, oldName, newName string, code []byte, aiDataSource types.AIDataSource) {
	key := types.DataSourceStoreKey(oldName)
	// if the user does not want to reuse the old name
	if oldName != newName {
		store := ctx.KVStore(k.storeKey)
		byteKey := []byte(key)
		store.Delete(byteKey)
		// delete the old file because it not pointed by any oScript
		k.EraseAIDataSourceFile(oldName)
		// add new file
		k.AddAIDataSourceFile(code, newName)
	} else {
		// edit the file instead since old name = new name
		k.EditAIDataSourceFile(code, oldName)
	}
	k.SetAIDataSource(ctx, newName, aiDataSource)
}

// AddAIDataSourceFile adds the file to filecache,
func (k Keeper) AddAIDataSourceFile(file []byte, name string) {
	k.fileCache.AddFile(file, types.DataSourceStoreFileString(name))
}

// EditAIDataSourceFile edit a file in the filecache,
func (k Keeper) EditAIDataSourceFile(file []byte, name string) {
	k.fileCache.EditFile(file, types.DataSourceStoreFileString(name))
}

// EraseAIDataSourceFile removes the file in the filecache,
func (k Keeper) EraseAIDataSourceFile(name string) {
	k.fileCache.EraseFile(types.DataSourceStoreFileString(name))
}

// GetAIDataSourceFile gets the data source code stored in the file storage
func (k Keeper) GetAIDataSourceFile(name string) ([]byte, error) {
	code, err := k.fileCache.GetFile(types.DataSourceStoreFileString(name))
	if err != nil {
		return []byte{}, err
	}
	return code, err
}
