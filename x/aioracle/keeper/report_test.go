package keeper_test

import (
	"testing"

	aioraclekeeper "github.com/oraichain/orai/x/aioracle/keeper"
	"github.com/oraichain/orai/x/aioracle/types"
	aioracletypes "github.com/oraichain/orai/x/aioracle/types"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestGetReportsBlockHeight(t *testing.T) {
	// init sim app
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	testKeeper := aioraclekeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())

	// init ai request
	id := ksuid.New().String()

	// init data source results
	dsResult1 := aioracletypes.NewResult(&aioracletypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult2 := aioracletypes.NewResult(&aioracletypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult3 := aioracletypes.NewResult(&aioracletypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{0x50}, types.ResultSuccess)
	dsResults := []*aioracletypes.Result{dsResult1, dsResult2, dsResult3}

	// init report
	report := types.NewReport(id, dsResults, 0, []byte{0x50}, nil, "")

	// verify report
	err := testKeeper.SetReport(ctx, id, report)
	require.NoError(t, err)

	reports := testKeeper.GetReportsBlockHeight(ctx)

	require.Equal(t, 1, len(reports))
}
