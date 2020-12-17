package websocket

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/websocket/types"
)

// InitGenesis initialize default parameters, can assign some coins in the chain here
// and the k's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	// TODO: Define logic for when you would like to initialize a new genesis
	// Init params for the airequest module
}

// DefaultGenesisState returns the default provider genesis state.
func DefaultGenesisState() GenesisState {
	return types.DefaultGenesisState()
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	// TODO: Define logic for exporting state
	return types.GenesisState{
		Reports:   []types.Report{},
		Reporters: []types.Reporter{},
		//Params:     k.GetParams(ctx),
	}
}
