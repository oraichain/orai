package app

import (
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
)

type TestSupport struct {
	t   *testing.T
	app *OraichainApp
}

func NewTestSupport(t *testing.T, app *OraichainApp) *TestSupport {
	return &TestSupport{t: t, app: app}
}

func (s TestSupport) IBCKeeper() ibckeeper.Keeper {
	return *s.app.ibcKeeper
}

func (s TestSupport) WasmKeeper() wasm.Keeper {
	return s.app.wasmKeeper
}

func (s TestSupport) AppCodec() codec.Codec {
	return s.app.appCodec
}

func (s TestSupport) StakingKeeper() stakingkeeper.Keeper {
	return s.app.stakingKeeper
}

func (s TestSupport) BankKeeper() bankkeeper.Keeper {
	return s.app.bankKeeper
}

func (s TestSupport) TransferKeeper() ibctransferkeeper.Keeper {
	return s.app.transferKeeper
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}
