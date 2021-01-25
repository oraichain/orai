package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryContract return data from input of smart contract, should be struct with json serialized
func (k *Keeper) QueryContract(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	return k.wasmKeeper.QuerySmart(ctx, contractAddr, req)
}
