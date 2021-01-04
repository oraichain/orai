package provider

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/types"
)

// InitGenesis initialize default parameters, can assign some coins in the chain here
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// TODO: Define logic for when you would like to initialize a new genesis
	for _, record := range data.OracleScripts {
		keeper.SetOracleScript(ctx, record.Name, record)
	}
	for _, record := range data.AIDataSources {
		keeper.SetAIDataSource(ctx, record.Name, record)
	}
	for _, record := range data.TestCases {
		keeper.SetTestCase(ctx, record.Name, record)
	}

	// Init params for the provider module
	keeper.SetParam(ctx, types.KeyOracleScriptRewardPercentage, data.Params.OracleScriptRewardPercentage)
	keeper.SetParam(ctx, types.KeyMaximumCodeBytes, data.Params.MaximumCodeBytes)
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
	var oScripts []OracleScript
	var aiDataSources []AIDataSource
	var testCases []TestCase
	oscriptIterator := k.GetAllOracleScriptNames(ctx)
	for ; oscriptIterator.Valid(); oscriptIterator.Next() {

		name := string(oscriptIterator.Key())
		oScript, _ := k.GetOracleScript(ctx, name)
		oScripts = append(oScripts, oScript)
	}

	dataSourceIterator := k.GetAllAIDataSourceNames(ctx)
	for ; dataSourceIterator.Valid(); dataSourceIterator.Next() {

		name := string(dataSourceIterator.Key())
		aiDataSource, _ := k.GetAIDataSource(ctx, name)
		aiDataSources = append(aiDataSources, aiDataSource)
	}

	testCaseIterator := k.GetAllTestCaseNames(ctx)
	for ; testCaseIterator.Valid(); testCaseIterator.Next() {

		name := string(testCaseIterator.Key())
		testCase, _ := k.GetTestCase(ctx, name)
		testCases = append(testCases, testCase)
	}

	return types.GenesisState{
		OracleScripts: oScripts,
		AIDataSources: aiDataSources,
		TestCases:     testCases,
		Params:        k.GetParams(ctx),
	}
}
