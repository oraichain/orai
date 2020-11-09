package keeper

import (
	"github.com/oraichain/orai/x/websocket/exported"
	"github.com/oraichain/orai/x/websocket/types"
)

// DefaultValResultI returns the default ai data source object
func (k Keeper) DefaultValResultI() exported.ValResultI {
	return k.DefaultValResult()
}

// DefaultValResult is a default constructor for the validator result
func (k Keeper) DefaultValResult() types.ValResult {
	return types.ValResult{
		Validator: nil,
		Result:    []byte{},
	}
}
