package app

import (
	"io"
	"os"

	aiRequest "github.com/oraichain/orai/x/airequest"
	aiResult "github.com/oraichain/orai/x/airesult"
	"github.com/oraichain/orai/x/provider"
	webSocket "github.com/oraichain/orai/x/websocket"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const (
	AppName          = "OraiApp"
	Bech32MainPrefix = "orai"
)

var (
	// TODO: rename your cli

	// DefaultCLIHome default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$PWD/.oraicli")

	// TODO: rename your daemon

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$PWD/.oraid")

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		supply.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		//gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler),
		params.AppModuleBasic{},
		//crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		//evidence.AppModuleBasic{},
		provider.AppModuleBasic{},
		aiRequest.AppModuleBasic{},
		webSocket.AppModuleBasic{},
		aiResult.AppModuleBasic{},
		// TODO: Add your module(s) AppModuleBasic
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		//gov.ModuleName:            {supply.Burner},
	}
)

// SetBech32AddressPrefixes sets the global Bech32 prefixes
func SetBech32AddressPrefixes(config *sdk.Config) {
	accountPrefix := Bech32MainPrefix
	validatorPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	consensusPrefix := Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	config.SetBech32PrefixForAccount(accountPrefix, accountPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(validatorPrefix, validatorPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(consensusPrefix, consensusPrefix+sdk.PrefixPublic)
}

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	//codec.RegisterEvidences(cdc)

	return cdc.Seal()
}

// NewApp extended ABCI application
type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	accountKeeper  auth.AccountKeeper // auth of accounts and transactions
	bankKeeper     bank.Keeper        // token transfer functions
	supplyKeeper   supply.Keeper      // total supply of the chain
	stakingKeeper  staking.Keeper     // pos layer, interacting with validators
	slashingKeeper slashing.Keeper    // validator punishment mechanisms
	mintKeeper     mint.Keeper        // create new token unit
	distrKeeper    distr.Keeper       // fee distribution
	//govKeeper      gov.Keeper         // on chain proposal and voting for token holders
	//crisisKeeper   crisis.Keeper      // halt the blockchain under some occasions
	paramsKeeper params.Keeper // global available param store
	//evidenceKeeper evidence.Keeper    // handling evidence of misbihaviour
	providerKeeper  provider.Keeper  // our own provider which provides a marketplace for AI providers to provide their AI services
	aiRequestKeeper aiRequest.Keeper // our own airequest which collects results from different AI models and aggregate them
	aiResultKeeper  aiResult.Keeper
	// TODO: Add your module(s)
	webSocketKeeper webSocket.Keeper

	// Module Manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// verify app interface at compile time
var _ simapp.App = (*NewApp)(nil)

// NewOraichainApp is a constructor function for dexaiApp
func NewOraichainApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, home string, baseAppOptions ...func(*bam.BaseApp),
) *NewApp {
	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(AppName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	// TODO: Add the keys that module requires
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		supply.StoreKey,
		staking.StoreKey,
		mint.StoreKey,
		distr.StoreKey,
		slashing.StoreKey,
		//gov.StoreKey,
		params.StoreKey,
		//evidence.StoreKey,
		provider.StoreKey,
		aiRequest.StoreKey,
		webSocket.StoreKey,
		aiResult.StoreKey,
	)

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &NewApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	// Set specific supspaces auth, bank, staking, distr, slashing, mint, evidence, gov, crisis
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)

	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)

	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)

	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)

	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)

	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)

	app.subspaces[provider.ModuleName] = app.paramsKeeper.Subspace(provider.DefaultParamspace)

	app.subspaces[aiRequest.ModuleName] = app.paramsKeeper.Subspace(aiRequest.DefaultParamspace)

	app.subspaces[webSocket.ModuleName] = app.paramsKeeper.Subspace(webSocket.DefaultParamspace)

	app.subspaces[aiResult.ModuleName] = app.paramsKeeper.Subspace(aiResult.DefaultParamspace)

	//app.subspaces[evidence.ModuleName] = app.paramsKeeper.Subspace(evidence.DefaultParamspace)

	//app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace)

	//app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.supplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		keys[distr.StoreKey],
		app.subspaces[distr.ModuleName],
		&stakingKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		app.subspaces[slashing.ModuleName],
	)

	app.mintKeeper = mint.NewKeeper(
		app.cdc,
		keys[mint.StoreKey],
		app.subspaces[mint.ModuleName],
		&stakingKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)

	app.distrKeeper = distr.NewKeeper(
		cdc, keys[distr.StoreKey],
		app.subspaces[distr.ModuleName],
		&stakingKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)

	// app.crisisKeeper = crisis.NewKeeper(
	// 	app.subspaces[crisis.ModuleName],
	// 	invCheckPeriod,
	// 	app.supplyKeeper,
	// 	auth.FeeCollectorName,
	// )

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks()),
	)

	// TODO: Add your module(s) keepers

	app.providerKeeper = provider.NewKeeper(
		app.cdc,
		keys[provider.StoreKey],
		app.subspaces[provider.ModuleName],
		".oraifiles",
	)

	app.aiRequestKeeper = aiRequest.NewKeeper(
		app.cdc,
		keys[aiRequest.StoreKey],
		app.subspaces[aiRequest.ModuleName],
		&stakingKeeper,
		app.providerKeeper,
	)

	app.webSocketKeeper = webSocket.NewKeeper(
		app.cdc,
		keys[webSocket.StoreKey],
		&stakingKeeper,
	)

	app.aiResultKeeper = aiResult.NewKeeper(
		app.cdc,
		keys[aiResult.StoreKey],
		app.subspaces[aiResult.ModuleName],
		app.supplyKeeper,
		app.bankKeeper,
		&stakingKeeper,
		app.distrKeeper,
		app.providerKeeper,
		app.webSocketKeeper,
		app.aiRequestKeeper,
		auth.FeeCollectorName,
	)

	// // Register the proposal types.
	// govRouter := gov.NewRouter()
	// govRouter.
	// 	AddRoute(gov.RouterKey, gov.ProposalHandler).
	// 	AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
	// 	AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))

	// app.govKeeper = gov.NewKeeper(cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.supplyKeeper, &stakingKeeper, govRouter)

	// // Create evidence keeper with evidence router.
	// evKeeper := evidence.NewKeeper(cdc, keys[evidence.StoreKey], app.subspaces[evidence.ModuleName], &stakingKeeper, app.slashingKeeper)

	// evidenceRouter := evidence.NewRouter()

	// evKeeper.SetRouter(evidenceRouter)

	// app.evidenceKeeper = *evKeeper

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),

		auth.NewAppModule(app.accountKeeper),

		bank.NewAppModule(app.bankKeeper, app.accountKeeper),

		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),

		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),

		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),

		// TODO: Add your module(s)

		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),

		mint.NewAppModule(app.mintKeeper),

		provider.NewAppModule(app.providerKeeper, app.subspaces[provider.ModuleName]),

		aiRequest.NewAppModule(app.aiRequestKeeper, app.stakingKeeper, app.providerKeeper, app.subspaces[aiRequest.ModuleName]),

		webSocket.NewAppModule(app.webSocketKeeper, app.stakingKeeper, app.subspaces[webSocket.ModuleName]),

		aiResult.NewAppModule(app.aiResultKeeper, app.supplyKeeper, app.bankKeeper, app.stakingKeeper, app.distrKeeper, app.providerKeeper, app.webSocketKeeper, app.aiRequestKeeper, app.subspaces[aiResult.ModuleName]),

		//gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),

		//crisis.NewAppModule(&app.crisisKeeper),

		//evidence.NewAppModule(app.evidenceKeeper),
	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.

	app.mm.SetOrderBeginBlockers(aiRequest.ModuleName, aiResult.ModuleName, distr.ModuleName, slashing.ModuleName, staking.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName, aiResult.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// app.mm.SetOrderInitGenesis(
	// 	auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName, supply.ModuleName,
	// 	slashing.ModuleName, gov.ModuleName, mint.ModuleName, provider.ModuleName, crisis.ModuleName,
	// 	genutil.ModuleName, evidence.ModuleName,
	// )

	app.mm.SetOrderInitGenesis(
		auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName, supply.ModuleName,
		slashing.ModuleName, mint.ModuleName, provider.ModuleName, aiRequest.ModuleName, webSocket.ModuleName, aiResult.ModuleName,
		genutil.ModuleName,
	)

	// register all module routes and module queriers
	//app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// // NewDefaultGenesisState generates the default state for the application.
// func NewDefaultGenesisState() GenesisState {
// 	return ModuleBasics.DefaultGenesis()
// }

// InitChainer application update at chain initialization
func (app *NewApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState

	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// BeginBlocker application updates every begin block
func (app *NewApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *NewApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// LoadHeight loads a particular height
func (app *NewApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *NewApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *NewApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *NewApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

// GetKey returns the KVStoreKey for the provided store key
func (app *NewApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *NewApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tKeys[storeKey]
}
