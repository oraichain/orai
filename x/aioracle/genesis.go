package aioracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// InitGenesis initialize default parameters, can assign some coins in the chain here
// and the k's address to pubkey map
func InitGenesis(ctx sdk.Context, k *Keeper, data *GenesisState) {
	// TODO: Define logic for when you would like to initialize a new genesis
	// Init params for the aioracle module
	k.SetRngSeed(ctx, make([]byte, types.RngSeedSize))
	k.SetParam(ctx, types.KeyMaximumAIOracleReqBytes, data.Params.MaximumAiOracleRequestBytes)
	k.SetParam(ctx, types.KeyMaximumAIOracleResBytes, data.Params.MaximumAiOracleResponseBytes)
	k.SetParam(ctx, types.KeyAIOracleRewardPercentages, data.Params.RewardAiOraclePercentages)
	k.SetParam(ctx, types.KeyReportPercentages, data.Params.ReportsPercentages)
}

// DefaultGenesisState returns the default aioracle genesis state.
func DefaultGenesisState() *GenesisState {
	return types.DefaultGenesisState()
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k *Keeper) (data *GenesisState) {
	// TODO: Define logic for exporting state
	return &types.GenesisState{
		AIOracles: []types.AIOracle{},
		Params:    k.GetParams(ctx),
	}
}
