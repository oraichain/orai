package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// SetRngSeed sets the rolling seed value to be provided value.
func (k *Keeper) SetRngSeed(ctx sdk.Context, rollingSeed []byte) {
	ctx.KVStore(k.StoreKey).Set(types.SeedStoreKey(), rollingSeed)
}

// GetRngSeed returns the current rolling seed value.
func (k *Keeper) GetRngSeed(ctx sdk.Context) []byte {
	return ctx.KVStore(k.StoreKey).Get(types.SeedStoreKey())
}

// ResolveRngSeed resolves the seed for the Rng package
func (k *Keeper) ResolveRngSeed(ctx sdk.Context, req abci.RequestBeginBlock) {
	oldSeed := k.GetRngSeed(ctx)
	if len(oldSeed) == 0 {
		k.SetRngSeed(ctx, make([]byte, types.RngSeedSize))
	} else {
		newSeed := oldSeed[types.NumSeedRemoval:]
		hash := req.GetHash()
		// generate a new seed by removing the first byte of the previous seed, and add a new byte from the new hash.
		for i := 0; i < types.NumSeedRemoval; i++ {
			newSeed = append(newSeed, hash[i])
		}
		k.SetRngSeed(ctx, newSeed)
	}
}