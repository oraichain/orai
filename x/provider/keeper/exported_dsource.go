package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/exported"
)

// This file implements all the exported functions of the Data Source struct that are exported by the Keeper

// GetAIDataSourceI returns the data source object given the name of the data source
func (k Keeper) GetAIDataSourceI(ctx sdk.Context, name string) (exported.AIDataSourceI, error) {
	dataSource, err := k.GetAIDataSource(ctx, name)
	if err != nil {
		return nil, err
	}
	return dataSource, nil
}

// DefaultAIDataSourceI returns the default ai data source object
func (k Keeper) DefaultAIDataSourceI() exported.AIDataSourceI {
	return k.DefaultAIDataSource()
}
