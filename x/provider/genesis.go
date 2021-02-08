package provider

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k *Keeper, genState *types.GenesisState) {
	k.SetParam(ctx, types.KeyMaximumCodeBytes, genState.Params.MaximumCodeBytes)
	k.SetParam(ctx, types.KeyOracleScriptRewardPercentage, genState.Params.OracleScriptRewardPercentage)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *Keeper) *types.GenesisState {
	return types.DefaultGenesisState()
}
