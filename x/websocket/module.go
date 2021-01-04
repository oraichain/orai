package websocket

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/oraichain/orai/x/websocket/client/cli"
	"github.com/oraichain/orai/x/websocket/client/rest"
	"github.com/oraichain/orai/x/websocket/keeper"
	"github.com/oraichain/orai/x/websocket/types"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}

	InSignals  = make(chan bool)
	OutSignals = make(chan bool)
)

// AppModuleBasic defines the basic application module used by the websocket module.
type AppModuleBasic struct{}

// Name returns the websocket module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the websocket module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the websocket
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the websocket module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the websocket module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, InSignals, OutSignals)	
}

// GetTxCmd returns the root tx command for the websocket module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the websocket module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the websocket module.
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
	// TODO: Add keepers that your application depends on
	stakingKeeper staking.Keeper
	params        params.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, s staking.Keeper, params params.Subspace /*TODO: Add Keepers that your application depends on*/) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		stakingKeeper:  s,
		params:         params,
		// TODO: Add keepers that your application depends on
	}
}

// Name returns the websocket module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the websocket module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the websocket module.
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the websocket module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the websocket module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the websocket module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the nameservice module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the websocket
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the websocket module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	// Run the begin blocker custom function to execute tasks when a new block starts
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the websocket module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// Run the end blocker custom function to execute tasks when end block
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
	//k.DirectAllocateTokens(ctx, req.GetLastCommitInfo().Votes)
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
	// 	TODO: fill out if your application requires endblock, if not you can delete this function
}
