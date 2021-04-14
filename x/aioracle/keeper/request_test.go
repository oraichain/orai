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
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestResolveRequestsFromReports(t *testing.T) {

	// define static variables
	PKS := simapp.CreateTestPubKeys(5)

	valConsPk1 := PKS[0]
	valConsPk2 := PKS[1]
	valConsPk3 := PKS[2]

	// init sim app
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrs(app, ctx, 10, sdk.NewInt(10000000000))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 10% commission and 300,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[0], valConsPk1, sdk.NewInt(300000000), true)

	// create second validator with 10% commission and 250,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[1], valConsPk2, sdk.NewInt(250000000), true)

	// create second validator with 10% commission and 150,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[2], valConsPk3, sdk.NewInt(150000000), true)

	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	testKeeper := aioraclekeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())

	// init ai request
	id := ksuid.New().String()
	aioracle := types.NewAIOracle(id, sdk.AccAddress{}, addrs[0], []sdk.ValAddress{valAddrs[0], valAddrs[1], valAddrs[2]}, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50})

	testKeeper.SetAIOracle(ctx, aioracle.RequestID, aioracle)
	//aioracleGet, err := testKeeper.GetAIOracle(ctx, id)

	// verify if set AI request is legit or not
	// require.NoError(t, err)
	// require.Equal(t, aioracleGet, aioracle)

	// init data source results
	dsResult1 := aioracletypes.NewDataSourceResult(&aioracletypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult2 := aioracletypes.NewDataSourceResult(&aioracletypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult3 := aioracletypes.NewDataSourceResult(&aioracletypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{0x50}, types.ResultSuccess)
	dsResults := []*aioracletypes.DataSourceResult{dsResult1, dsResult2, dsResult3}

	// init report
	report := types.NewReport(id, dsResults, 1, []byte{0x50}, valAddrs[0], "")

	// verify report
	err := testKeeper.AddReport(ctx, id, report)
	require.NoError(t, err)

	reward := aioracletypes.DefaultReward(1)
	testKeeper.ResolveRequestsFromReports(ctx, report, reward)

	// Remember that we are using coins, so the values will be truncated along the way. When applying the equations in the medium post, remember to truncate the final result.
	t.Logf("provider fees: %v\n", reward.ProviderFees)
	t.Logf("validator fees: %v\n", reward.ValidatorFees)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(13))), reward.ProviderFees)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(5))), reward.ValidatorFees)
}
