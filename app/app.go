package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kava-labs/kava/precompile/registry"
	ibchooks "github.com/osmosis-labs/osmosis/x/ibc-hooks"
	ibchookskeeper "github.com/osmosis-labs/osmosis/x/ibc-hooks/keeper"
	ibchookstypes "github.com/osmosis-labs/osmosis/x/ibc-hooks/types"

	packetforward "github.com/strangelove-ventures/packet-forward-middleware/v4/router"
	packetforwardkeeper "github.com/strangelove-ventures/packet-forward-middleware/v4/router/keeper"
	packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v4/router/types"

	// Token Factory
	"github.com/CosmWasm/token-factory/x/tokenfactory"
	// bindings "github.com/CosmWasm/token-factory/x/tokenfactory/bindings"
	tokenfactorykeeper "github.com/CosmWasm/token-factory/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/CosmWasm/token-factory/x/tokenfactory/types"

	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v4/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v4/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v4/modules/apps/29-fee/types"
	transfer "github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v4/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"

	// Note: please do your research before using this in production app, this is a demo and not an officially
	// supported IBC team implementation. It has no known issues, but do your own research before using it.
	intertx "github.com/cosmos/interchain-accounts/x/inter-tx"
	intertxkeeper "github.com/cosmos/interchain-accounts/x/inter-tx/keeper"
	intertxtypes "github.com/cosmos/interchain-accounts/x/inter-tx/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/oraichain/orai/doc/statik"

	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"

	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	ibcclientclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/client"
	appparams "github.com/oraichain/orai/app/params"
	appconfig "github.com/oraichain/orai/cmd/config"

	"github.com/CosmosContracts/juno/v15/x/clock"
	clockkeeper "github.com/CosmosContracts/juno/v15/x/clock/keeper"
	clocktypes "github.com/CosmosContracts/juno/v15/x/clock/types"

	"github.com/CosmosContracts/juno/v15/x/globalfee"
	globalfeekeeper "github.com/CosmosContracts/juno/v15/x/globalfee/keeper"
	globalfeetypes "github.com/CosmosContracts/juno/v15/x/globalfee/types"

	"github.com/CosmWasm/token-factory/x/tokenfactory/bindings"
	evmutil "github.com/kava-labs/kava/x/evmutil"
	evmutilkeeper "github.com/kava-labs/kava/x/evmutil/keeper"
	evmutiltypes "github.com/kava-labs/kava/x/evmutil/types"
	evmante "github.com/tharsis/ethermint/app/ante"
	ethermintconfig "github.com/tharsis/ethermint/server/config"
	"github.com/tharsis/ethermint/x/evm"
	evmrest "github.com/tharsis/ethermint/x/evm/client/rest"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"github.com/tharsis/ethermint/x/feemarket"
	feemarketkeeper "github.com/tharsis/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
)

const appName = "Oraichain"

// We pull these out so we can set them with LDFLAGS in the Makefile
var (
	NodeDir = ".oraid"

	BinaryVersion = "v0.42.4"

	// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
	// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
	ProposalsEnabled = "true"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	// https://github.com/oraichain/orai/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificProposals = ""
	EnabledCapabilities     = []string{
		tokenfactorytypes.EnableBurnFrom,
		// tokenfactorytypes.EnableForceTransfer,
		tokenfactorytypes.EnableSetMetadata,
	}
)

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

// These constants are derived from the above variables.
// These are the ones we will want to use in the code, based on
// any overrides above
var (
	// DefaultNodeHome default home directories for wasmd
	DefaultNodeHome = os.ExpandEnv("$PWD/") + NodeDir

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = appconfig.Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = appconfig.Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = appconfig.Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = appconfig.Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = appconfig.Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = appconfig.Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

var (
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			append(
				wasmclient.ProposalHandlers,
				paramsclient.ProposalHandler,
				distrclient.ProposalHandler,
				upgradeclient.ProposalHandler,
				upgradeclient.CancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler,
				ibcclientclient.UpgradeProposalHandler,
			)...,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		wasm.AppModuleBasic{},
		ica.AppModuleBasic{},
		intertx.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
		clock.AppModuleBasic{},
		globalfee.AppModuleBasic{},
		ibchooks.AppModuleBasic{},
		packetforward.AppModuleBasic{},
		evmutil.AppModuleBasic{},
		tokenfactory.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner}, // used for secure addition and subtraction of balance using module account
		evmutiltypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		ibcfeetypes.ModuleName:         nil,
		icatypes.ModuleName:            nil,
		wasm.ModuleName:                {authtypes.Burner},
		tokenfactorytypes.ModuleName:   {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}
)

// Verify app interface at compile time
var (
	_ simapp.App              = (*OraichainApp)(nil)
	_ servertypes.Application = (*OraichainApp)(nil)
)

// Options bundles several configuration params for an App.
type EvmOptions struct {
	MempoolEnableAuth    bool
	MempoolAuthAddresses []sdk.AccAddress
	EVMTrace             string
	EVMMaxGasWanted      uint64
}

// DefaultOptions is a sensible default Options value.
var DefaultEvmOptions = EvmOptions{
	EVMTrace:        ethermintconfig.DefaultEVMTracer,
	EVMMaxGasWanted: ethermintconfig.DefaultMaxTxGasWanted,
}

// OraichainApp extended ABCI application
type OraichainApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper      authkeeper.AccountKeeper
	BankKeeper         bankkeeper.BaseKeeper
	CapabilityKeeper   *capabilitykeeper.Keeper
	StakingKeeper      stakingkeeper.Keeper
	SlashingKeeper     slashingkeeper.Keeper
	MintKeeper         mintkeeper.Keeper
	DistrKeeper        distrkeeper.Keeper
	GovKeeper          govkeeper.Keeper
	CrisisKeeper       crisiskeeper.Keeper
	UpgradeKeeper      upgradekeeper.Keeper
	ParamsKeeper       paramskeeper.Keeper
	IbcKeeper          *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvmKeeper          *evmkeeper.Keeper
	EvmutilKeeper      evmutilkeeper.Keeper
	FeeMarketKeeper    feemarketkeeper.Keeper
	EvidenceKeeper     evidencekeeper.Keeper
	TransferKeeper     ibctransferkeeper.Keeper
	WasmKeeper         wasm.Keeper
	FeeGrantKeeper     feegrantkeeper.Keeper
	AuthzKeeper        authzkeeper.Keeper
	ContractKeeper     *wasmkeeper.PermissionedKeeper
	ClockKeeper        clockkeeper.Keeper
	GlobalFeeKeeper    globalfeekeeper.Keeper
	TokenFactoryKeeper tokenfactorykeeper.Keeper

	IbcFeeKeeper        ibcfeekeeper.Keeper
	IBCHooksKeeper      *ibchookskeeper.Keeper
	IcaControllerKeeper icacontrollerkeeper.Keeper
	IcaHostKeeper       icahostkeeper.Keeper
	InterTxKeeper       intertxkeeper.Keeper
	// Middleware wrapper
	Ics20WasmHooks      *ibchooks.WasmHooks
	HooksICS4Wrapper    ibchooks.ICS4Middleware
	PacketForwardKeeper *packetforwardkeeper.Keeper

	// custom modules here

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper
	// scopedIBCFeeKeeper        capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedInterTxKeeper       capabilitykeeper.ScopedKeeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// module configurator
	configurator module.Configurator
}

// NewOraichainApp returns a reference to an initialized OraichainApp.
func NewOraichainApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	skipUpgradeHeights map[int64]bool, homePath string, invCheckPeriod uint, encodingConfig appparams.EncodingConfig, enabledProposals []wasm.ProposalType,
	appOpts servertypes.AppOptions, wasmOpts []wasm.Option, options EvmOptions, baseAppOptions ...func(*baseapp.BaseApp)) *OraichainApp {

	appCodec, legacyAmino := encodingConfig.Codec, encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, evmtypes.StoreKey, feemarkettypes.StoreKey, capabilitytypes.StoreKey,
		wasm.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey, icahosttypes.StoreKey,
		icacontrollertypes.StoreKey, intertxtypes.StoreKey, ibcfeetypes.StoreKey, globalfeetypes.StoreKey,
		ibchookstypes.StoreKey, clocktypes.StoreKey, packetforwardtypes.StoreKey, evmutiltypes.StoreKey, tokenfactorytypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &OraichainApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedInterTxKeeper := app.CapabilityKeeper.ScopeToModule(intertxtypes.ModuleName)
	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.getSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.getSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	validateKeeper(app.AccountKeeper, app.BankKeeper)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.getSubspace(stakingtypes.ModuleName),
	)
	validateKeeper(stakingKeeper)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.getSubspace(minttypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.getSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.getSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.getSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, bApp)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.BaseApp.MsgServiceRouter(),
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	// NOTE: upgrade feature comment out for new build after proposal done
	// app.UpgradeKeeper.SetUpgradeHandler("ai-oracle", func(ctx sdk.Context, plan upgradetypes.Plan) {
	// 	// TODO: Add some modification logic here
	// })

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	validateKeeper(app.DistrKeeper, app.SlashingKeeper)
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// Create IBC Keeper
	validateKeeper(scopedIBCKeeper, app.UpgradeKeeper)
	app.IbcKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibchost.StoreKey],
		app.getSubspace(ibchost.ModuleName),
		stakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// Create Ethermint keepers
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, keys[feemarkettypes.StoreKey], app.getSubspace(feemarkettypes.ModuleName),
	)

	app.EvmutilKeeper = evmutilkeeper.NewKeeper(
		appCodec,
		keys[evmutiltypes.StoreKey],
		app.getSubspace(evmutiltypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
	)

	validateKeeper(app.FeeMarketKeeper)
	validateKeeper(app.EvmutilKeeper)
	evmBankKeeper := evmutilkeeper.NewEvmBankKeeperWithDenoms(app.EvmutilKeeper, app.BankKeeper, app.AccountKeeper, appconfig.EvmDenom, appconfig.CosmosDenom)
	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey], app.getSubspace(evmtypes.ModuleName),
		app.AccountKeeper, evmBankKeeper, app.StakingKeeper, app.FeeMarketKeeper,
		options.EVMTrace,
	)
	app.EvmutilKeeper.SetEvmKeeper(app.EvmKeeper)

	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		app.keys[ibchookstypes.StoreKey],
	)
	validateKeeper(hooksKeeper)
	app.IBCHooksKeeper = &hooksKeeper

	validateKeeper(app.IBCHooksKeeper)
	wasmHooks := ibchooks.NewWasmHooks(app.IBCHooksKeeper, nil, sdk.GetConfig().GetBech32AccountAddrPrefix()) // The contract keeper needs to be set later
	app.Ics20WasmHooks = &wasmHooks
	validateKeeper(app.IbcKeeper, app.IbcKeeper.ChannelKeeper, app.Ics20WasmHooks)
	app.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		app.IbcKeeper.ChannelKeeper,
		app.Ics20WasmHooks,
	)

	// IBC Fee Module keeper
	validateKeeper(app.HooksICS4Wrapper, app.IbcKeeper.PortKeeper)
	app.IbcFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec, keys[ibcfeetypes.StoreKey], app.getSubspace(ibcfeetypes.ModuleName),
		app.HooksICS4Wrapper, // may be replaced with IBC middleware
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	// Initialize packet forward middleware router
	validateKeeper(app.IbcFeeKeeper)
	app.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec, app.keys[packetforwardtypes.StoreKey],
		app.getSubspace(packetforwardtypes.ModuleName),
		nil, // Will be zero-value here. Reference is set later on with SetTransferKeeper.
		app.IbcKeeper.ChannelKeeper,
		app.DistrKeeper,
		app.BankKeeper,
		// The ICS4Wrapper is replaced by the IBCFeeKeeper instead of the channel so that sending can be overridden by the middleware
		&app.IbcFeeKeeper,
	)

	// Create Transfer Keepers
	validateKeeper(app.PacketForwardKeeper, scopedTransferKeeper)
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.getSubspace(ibctransfertypes.ModuleName),
		app.PacketForwardKeeper,
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	app.PacketForwardKeeper.SetTransferKeeper(app.TransferKeeper)

	tokenFactoryKeeper := tokenfactorykeeper.NewKeeper(
		keys[tokenfactorytypes.StoreKey],
		app.getSubspace(tokenfactorytypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		EnabledCapabilities,
	)
	app.TokenFactoryKeeper = tokenFactoryKeeper

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	app.EvidenceKeeper = *evidenceKeeper

	validateKeeper(scopedICAHostKeeper)
	app.IcaHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		keys[icahosttypes.StoreKey],
		app.getSubspace(icahosttypes.SubModuleName),
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)
	validateKeeper(scopedICAControllerKeeper)
	app.IcaControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		keys[icacontrollertypes.StoreKey],
		app.getSubspace(icacontrollertypes.SubModuleName),
		app.IbcFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		scopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)

	validateKeeper(app.IcaControllerKeeper, scopedInterTxKeeper)
	// For wasmd we use the demo controller from https://github.com/cosmos/interchain-accounts but see notes below
	app.InterTxKeeper = intertxkeeper.NewKeeper(appCodec, keys[intertxtypes.StoreKey], app.IcaControllerKeeper, scopedInterTxKeeper)

	// set the contract keeper for the Ics20WasmHooks
	// just re-use the full router - do we want to limit this more?
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	validateKeeper(scopedWasmKeeper, app.TransferKeeper)
	// Setup wasm bindings
	wasmOpts = append(bindings.RegisterCustomPlugins(&app.BankKeeper, &app.TokenFactoryKeeper), wasmOpts...)
	wasmOpts = append(RegisterStargateQueries(*bApp.GRPCQueryRouter(), appCodec), wasmOpts...)

	app.WasmKeeper = wasm.NewKeeper(
		appCodec,
		keys[wasm.StoreKey],
		app.getSubspace(wasm.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		filepath.Join(homePath, "wasm"),
		wasmConfig,
		strings.Join(AllCapabilities(), ","),
		wasmOpts...,
	)

	// set the contract keeper for the Ics20WasmHooks
	validateKeeper(app.WasmKeeper)
	app.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	validateKeeper(app.ContractKeeper)
	app.Ics20WasmHooks.ContractKeeper = app.ContractKeeper

	app.ClockKeeper = clockkeeper.NewKeeper(
		app.keys[clocktypes.StoreKey],
		appCodec,
		*app.ContractKeeper,
	)

	app.GlobalFeeKeeper = globalfeekeeper.NewKeeper(
		appCodec,
		app.keys[globalfeetypes.StoreKey],
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IbcKeeper.ClientKeeper)).
		AddRoute(clocktypes.RouterKey, clockkeeper.NewClockProposalHandler(app.ClockKeeper))

	// The gov proposal types can be individually enabled
	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, enabledProposals))
	}

	// Create Transfer Stack
	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		app.PacketForwardKeeper,
		1,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, app.IbcFeeKeeper)
	transferStack = ibchooks.NewIBCMiddleware(transferStack, &app.HooksICS4Wrapper)

	// Create Interchain Accounts Stack
	// SendPacket, since it is originating from the application to core IBC:
	// icaAuthModuleKeeper.SendTx -> icaController.SendPacket -> fee.SendPacket -> channel.SendPacket

	// Note: please do your research before using this in production app, this is a demo and not an officially
	// supported IBC team implementation. Do your own research before using it.
	var icaControllerStack porttypes.IBCModule
	// You will likely want to use your own reviewed and maintained ica auth module
	icaControllerStack = intertx.NewIBCModule(app.InterTxKeeper)
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, app.IcaControllerKeeper)
	icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerStack, app.IbcFeeKeeper)

	// RecvPacket, message that originates from core IBC and goes down to app, the flow is:
	// channel.RecvPacket -> fee.OnRecvPacket -> icaHost.OnRecvPacket
	var icaHostStack porttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(app.IcaHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, app.IbcFeeKeeper)

	// Create fee enabled wasm ibc Stack
	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(app.WasmKeeper, app.IbcKeeper.ChannelKeeper, app.IbcFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, app.IbcFeeKeeper)

	// create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasm.ModuleName, wasmStack).
		AddRoute(intertxtypes.ModuleName, icaControllerStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostStack)

	app.IbcKeeper.SetRouter(ibcRouter)

	validateKeeper(govRouter)
	app.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.getSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		govRouter,
	)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IbcKeeper),
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		transfer.NewAppModule(app.TransferKeeper),
		ibcfee.NewAppModule(app.IbcFeeKeeper),
		ica.NewAppModule(&app.IcaControllerKeeper, &app.IcaHostKeeper),
		intertx.NewAppModule(appCodec, app.InterTxKeeper),
		clock.NewAppModule(appCodec, app.ClockKeeper),
		globalfee.NewAppModule(appCodec, app.GlobalFeeKeeper, appconfig.MinimalDenom),
		ibchooks.NewAppModule(app.AccountKeeper),
		packetforward.NewAppModule(app.PacketForwardKeeper),
		evmutil.NewAppModule(app.EvmutilKeeper, app.BankKeeper),
		tokenfactory.NewAppModule(app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		// additional non simd modules
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		intertxtypes.ModuleName,
		wasm.ModuleName,
		ibchookstypes.ModuleName,
		clocktypes.ModuleName,
		globalfee.ModuleName,
		evmutiltypes.ModuleName,
		tokenfactorytypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		evmtypes.ModuleName,
		// fee market module must go after evm module in order to retrieve the block gas used.
		feemarkettypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		// additional non simd modules
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		intertxtypes.ModuleName,
		wasm.ModuleName,
		ibchookstypes.ModuleName,
		clocktypes.ModuleName,
		globalfee.ModuleName,
		evmutiltypes.ModuleName,
		tokenfactorytypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	// NOTE: wasm module should be at the end as it can call other module functionality direct or via message dispatching during
	// genesis phase. For example bank transfer, auth account check, staking, ...
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		// additional non simd modules
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		intertxtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		evmutiltypes.ModuleName,
		// wasm after ibc transfer
		wasm.ModuleName,
		ibchookstypes.ModuleName,
		clocktypes.ModuleName,
		globalfee.ModuleName,
		tokenfactorytypes.ModuleName,

		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// add test gRPC service for testing gRPC queries in isolation
	// testdata.RegisterTestServiceServer(app.GRPCQueryRouter(), testdata.QueryImpl{}) // TODO: this is testdata !!!!

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IbcKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		tokenfactory.NewAppModule(app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			AccountKeeper: app.AccountKeeper,
			BankKeeper:    app.BankKeeper,
			EvmKeeper:     app.EvmKeeper,

			GlobalFeeKeeper: app.GlobalFeeKeeper,
			StakingKeeper:   app.StakingKeeper,

			FeegrantKeeper:       app.FeeGrantKeeper,
			FeeMarketKeeper:      app.FeeMarketKeeper,
			SignModeHandler:      encodingConfig.TxConfig.SignModeHandler(),
			SigGasConsumer:       evmante.DefaultSigVerificationGasConsumer,
			MaxTxGasWanted:       options.EVMMaxGasWanted,
			IBCKeeper:            app.IbcKeeper,
			TxCounterStoreKey:    keys[wasm.StoreKey],
			WasmConfig:           wasmConfig,
			Cdc:                  appCodec,
			BypassMinFeeMsgTypes: []string{},
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}
	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	// handle statesync for cosmwasm
	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic("failed to register snapshot extension: " + err.Error())
		}
	}

	// set upgrade module
	app.upgradeHandler()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}
		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedICAControllerKeeper = scopedICAControllerKeeper
	app.ScopedInterTxKeeper = scopedInterTxKeeper
	clockkeeper.RegisterProposalTypes()

	// register wasm keeper
	registry.InitializePrecompiles(app.ContractKeeper, app.WasmKeeper, app.EvmKeeper)

	return app
}

// Name returns the name of the App
func (app *OraichainApp) Name() string { return app.BaseApp.Name() }

// application updates every begin block
func (app *OraichainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *OraichainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *OraichainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *OraichainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *OraichainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *OraichainApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *OraichainApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// SimulationManager implements the SimulationApp interface
func (app *OraichainApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// getSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *OraichainApp) getSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *OraichainApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	evmrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}

}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *OraichainApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *OraichainApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func (app *OraichainApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName).WithKeyTable(packetforwardtypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)
	paramsKeeper.Subspace(ibchookstypes.ModuleName)
	paramsKeeper.Subspace(clocktypes.ModuleName)
	paramsKeeper.Subspace(globalfeetypes.ModuleName)
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	paramsKeeper.Subspace(evmutiltypes.ModuleName)
	paramsKeeper.Subspace(tokenfactorytypes.ModuleName)

	return paramsKeeper
}

func (app *OraichainApp) upgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(BinaryVersion, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		app.GlobalFeeKeeper.SetParams(ctx, globalfeetypes.Params{MinimumGasPrices: sdk.NewDecCoins(sdk.NewDecCoin(appconfig.MinimalDenom, sdk.NewInt(100)))})
		response, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
		return response, err
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == BinaryVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storetypes.StoreUpgrades{
			Added: []string{globalfeetypes.StoreKey},
		}))
	}
}

// AllCapabilities returns all capabilities available with the current wasmvm
// See https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
// This functionality is going to be moved upstream: https://github.com/CosmWasm/wasmvm/issues/425
func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
		"cosmwasm_1_3",
		"cosmwasm_1_4",
	}
}
