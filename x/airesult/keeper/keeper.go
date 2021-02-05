package keeper

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	params "github.com/cosmos/cosmos-sdk/x/params/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/oraichain/orai/x/airequest"
	"github.com/oraichain/orai/x/airesult/types"
	"github.com/oraichain/orai/x/provider"
	"github.com/oraichain/orai/x/websocket"
	"github.com/tendermint/tendermint/libs/log"
)

// always clone keeper to make it immutable
type (
	Keeper struct {
		cdc              codec.Marshaler
		storeKey         sdk.StoreKey
		wasmKeeper       *wasm.Keeper
		paramSpace       params.Subspace
		stakingKeeper    staking.Keeper
		providerKeeper   *provider.Keeper
		webSocketKeeper  *websocket.Keeper
		aiRequestKeeper  *airequest.Keeper
		bankKeeper       bank.Keeper
		distrKeeper      distr.Keeper
		authKeeper       auth.AccountKeeper
		feeCollectorName string
	}

	// TestKeeper is created solely for unit test
	TestKeeper struct {
		Keeper           Keeper
		Cdc              codec.Marshaler
		StoreKey         sdk.StoreKey
		WasmKeeper       *wasm.Keeper
		ParamSpace       params.Subspace
		StakingKeeper    staking.Keeper
		ProviderKeeper   *provider.Keeper
		WebSocketKeeper  *websocket.Keeper
		AiRequestKeeper  *airequest.Keeper
		BankKeeper       bank.Keeper
		DistrKeeper      distr.Keeper
		AuthKeeper       auth.AccountKeeper
		FeeCollectorName string
	}
)

// NewKeeper creates a airequest keeper
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, wasmKeeper *wasm.Keeper, subspace params.Subspace, stakingKeeper staking.Keeper, providerKeeper *provider.Keeper, bankKeeper bank.Keeper, distrKeeper distr.Keeper, authKeeper auth.AccountKeeper, webSocketKeeper *websocket.Keeper, aiRequestKeeper *airequest.Keeper, feeCollectorName string) *Keeper {
	if !subspace.HasKeyTable() {
		// register parameters of the airequest module into the param space
		subspace = subspace.WithKeyTable(types.ParamKeyTable())
	}
	return &Keeper{
		storeKey:         key,
		cdc:              cdc,
		wasmKeeper:       wasmKeeper,
		paramSpace:       subspace,
		stakingKeeper:    stakingKeeper,
		providerKeeper:   providerKeeper,
		bankKeeper:       bankKeeper,
		distrKeeper:      distrKeeper,
		authKeeper:       authKeeper,
		webSocketKeeper:  webSocketKeeper,
		aiRequestKeeper:  aiRequestKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// NewTestKeeper creates a airequest keeper for testing
func NewTestKeeper(keeper Keeper, cdc codec.Marshaler, key sdk.StoreKey, wasmKeeper *wasm.Keeper, subspace params.Subspace, stakingKeeper staking.Keeper, providerKeeper *provider.Keeper, bankKeeper bank.Keeper, distrKeeper distr.Keeper, authKeeper auth.AccountKeeper, webSocketKeeper *websocket.Keeper, aiRequestKeeper *airequest.Keeper, feeCollectorName string) *TestKeeper {
	if !subspace.HasKeyTable() {
		// register parameters of the airequest module into the param space
		subspace = subspace.WithKeyTable(types.ParamKeyTable())
	}
	return &TestKeeper{
		Keeper:           keeper,
		StoreKey:         key,
		Cdc:              cdc,
		WasmKeeper:       wasmKeeper,
		ParamSpace:       subspace,
		StakingKeeper:    stakingKeeper,
		ProviderKeeper:   providerKeeper,
		BankKeeper:       bankKeeper,
		DistrKeeper:      distrKeeper,
		AuthKeeper:       authKeeper,
		WebSocketKeeper:  webSocketKeeper,
		AiRequestKeeper:  aiRequestKeeper,
		FeeCollectorName: feeCollectorName,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetExpirationCountParam returns the expiration count from param
func (k *Keeper) GetExpirationCountParam(ctx sdk.Context) int64 {
	// TODO
	percentage := k.GetParam(ctx, types.KeyExpirationCount)
	return int64(percentage)
}

// GetTotalReportsParam returns the total reports from param
func (k *Keeper) GetTotalReportsParam(ctx sdk.Context) int64 {
	// TODO
	percentage := k.GetParam(ctx, types.KeyTotalReports)
	return int64(percentage)
}
