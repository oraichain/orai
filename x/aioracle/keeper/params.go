package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// GetParam returns the parameter as specified by key as an uint64.
func (k *Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
	k.ParamSpace.Get(ctx, key, &res)
	return res
}

// SetParam saves the given key-value parameter to the store.
func (k *Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
	k.ParamSpace.Set(ctx, key, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ParamSpace.GetParamSet(ctx, &params)
	return params
}
