package airequest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airequest/types"
)

// InitGenesis initialize default parameters, can assign some coins in the chain here
// and the k's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	// TODO: Define logic for when you would like to initialize a new genesis
	// Init params for the airequest module
	k.SetRngSeed(ctx, make([]byte, types.RngSeedSize))
	k.SetParam(ctx, types.KeyMaximumRequestBytes, data.Params.MaximumRequestBytes)
}

// DefaultGenesisState returns the default airequest genesis state.
func DefaultGenesisState() GenesisState {
	return types.DefaultGenesisState()
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	// TODO: Define logic for exporting state
	return types.GenesisState{
		AIRequests: []types.AIRequest{},
		Params:     k.GetParams(ctx),
	}
}
