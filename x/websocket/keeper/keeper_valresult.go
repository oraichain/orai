package keeper

import (
	"github.com/oraichain/orai/x/websocket/types"
)

// DefaultValResultI returns the default ai data source object
func (k Keeper) DefaultValResultI() types.ValResult {
	return k.DefaultValResult()
}

// DefaultValResult is a default constructor for the validator result
func (k Keeper) DefaultValResult() types.ValResult {
	return types.ValResult{
		Validator: &types.Validator{},
		Result:    []byte{},
	}
}

// GetKeyResultSuccess is a getter to collect the result success key for validator result verification using by other modules.
func (k Keeper) GetKeyResultSuccess() string {
	return types.ResultSuccess
}

// NewValResult is a wrapper function of the websocket module that allow others to initiate a new valresult entity through the keeper
func (k Keeper) NewValResult(val *types.Validator, result []byte, status string) *types.ValResult {
	return &types.ValResult{
		Validator:    val,
		Result:       result,
		ResultStatus: status,
	}
}
