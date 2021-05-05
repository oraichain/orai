package keeper_test

import (
	"testing"

	aioraclekeeper "github.com/oraichain/orai/x/aioracle/keeper"
	"github.com/oraichain/orai/x/aioracle/types"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestResolveResult(t *testing.T) {

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

	// init data source results
	dsResult1 := types.NewResult(&types.EntryPoint{}, []byte{0x50}, "success")
	dsResult2 := types.NewResult(&types.EntryPoint{}, []byte{0x50}, "success")
	dsResult3 := types.NewResult(&types.EntryPoint{}, []byte{0x50}, "success")
	dsResults := []*types.Result{dsResult1, dsResult2, dsResult3}

	// init report
	id := ksuid.New().String()
	report := types.NewReport(id, dsResults, 1, []byte{0x50}, valAddrs[0], types.ResultSuccess)

	// verify report
	err := testKeeper.AddReport(ctx, id, report)
	require.NoError(t, err)

	// test result with 2 validators aka two reports, total 70% reports to finish
	testKeeper.ResolveResult(ctx, report, 100, 70)

	result, err := testKeeper.GetResult(ctx, id)
	t.Logf("result after resolve with 2 validators: %v\n", result.Status)

	for i := 0; i < 68; i++ {
		// test result with 100 validators aka 100 reports, total 70% reports to finish
		testKeeper.ResolveResult(ctx, report, 100, 70)
	}

	result, err = testKeeper.GetResult(ctx, id)
	t.Logf("result after resolve with 100 validators and 68 reports: %v\n", result.Status)
	require.Equal(t, "pending", result.Status)

	// add another report to the list of 100 reports will lead to finished
	testKeeper.ResolveResult(ctx, report, 100, 70)
	result, err = testKeeper.GetResult(ctx, id)
	t.Logf("result after resolve with 100 validators and 69 reports: %v\n", result.Status)
	require.Equal(t, "finished", result.Status)

	id = ksuid.New().String()
	report = types.NewReport(id, dsResults, 1, []byte{0x50}, valAddrs[0], types.ResultSuccess)

	// verify report
	err = testKeeper.AddReport(ctx, id, report)
	require.NoError(t, err)

	testKeeper.ResolveResult(ctx, report, 1, 70)

	result, err = testKeeper.GetResult(ctx, id)
	t.Logf("result after resolve with 1 validator: %v\n", result.Status)
	require.Equal(t, "finished", result.Status)

}
