package airesult

import (
	"encoding/json"

	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	aiRequest "github.com/oraichain/orai/x/airequest"
	"github.com/oraichain/orai/x/airesult/client/cli"
	"github.com/oraichain/orai/x/airesult/client/rest"
	"github.com/oraichain/orai/x/airesult/keeper"
	"github.com/oraichain/orai/x/airesult/types"
	"github.com/oraichain/orai/x/provider"
	webSocket "github.com/oraichain/orai/x/websocket"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the airequest module.
type AppModuleBasic struct{}

// Name returns the airequest module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the airequest module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the airequest
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the airequest module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the airequest module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the airequest module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the airequest module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the airequest module.
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
	// TODO: Add keepers that your application depends on
	supplyKeeper    supply.Keeper
	bankKeeper      bank.Keeper
	stakingKeeper   staking.Keeper
	distrKeeper     distr.Keeper
	providerKeeper  provider.Keeper
	webSocketKeeper webSocket.Keeper
	aiRequestKeeper aiRequest.Keeper
	params          params.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, s supply.Keeper, b bank.Keeper, staking staking.Keeper, distr distr.Keeper, p provider.Keeper, ws webSocket.Keeper, ar aiRequest.Keeper, params params.Subspace /*TODO: Add Keepers that your application depends on*/) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		keeper:          k,
		supplyKeeper:    s,
		bankKeeper:      b,
		stakingKeeper:   staking,
		distrKeeper:     distr,
		providerKeeper:  p,
		webSocketKeeper: ws,
		aiRequestKeeper: ar,
		params:          params,
		// TODO: Add keepers that your application depends on
	}
}

// Name returns the airequest module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the airequest module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the airequest module.
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the airequest module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the airequest module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the airequest module sdk.Querier.
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

// ExportGenesis returns the exported genesis state as raw bytes for the airequest
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the airequest module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	// Run the begin blocker custom function to execute tasks when a new block starts
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the airequest module. It returns no validator
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
	k.AllocateTokens(ctx, req.GetLastCommitInfo().Votes)
	//k.DirectAllocateTokens(ctx, req.GetLastCommitInfo().Votes)
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
	// 	TODO: fill out if your application requires endblock, if not you can delete this function
	k.ProcessReward(ctx)
}
