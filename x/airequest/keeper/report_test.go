package keeper_test

import (
	"testing"

	aiRequestkeeper "github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
	aiRequesttypes "github.com/oraichain/orai/x/airequest/types"
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
	testKeeper := aiRequestkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())

	// init ai request
	id := ksuid.New().String()

	// init data source results
	dsResult1 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult2 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult3 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{0x50}, types.ResultSuccess)
	dsResults := []*aiRequesttypes.Result{dsResult1, dsResult2, dsResult3}

	// init report
	report := types.NewReport(id, dsResults, 0, []byte{0x50}, nil, "")

	// verify report
	err := testKeeper.SetReport(ctx, id, report)
	require.NoError(t, err)

	reports := testKeeper.GetReportsBlockHeight(ctx)

	require.Equal(t, 1, len(reports))
}
