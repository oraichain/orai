package keeper_test

// import (
// 	"testing"

// 	aioraclekeeper "github.com/oraichain/orai/x/aioracle/keeper"
// 	"github.com/oraichain/orai/x/aioracle/keeper"
// 	aioraclekeeper "github.com/oraichain/orai/x/aioracle/keeper"
// 	aioracletypes "github.com/oraichain/orai/x/aioracle/types"
// 	"github.com/segmentio/ksuid"
// 	"github.com/stretchr/testify/require"
// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

// 	"github.com/cosmos/cosmos-sdk/simapp"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/auth/types"
// 	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
// )

// func TestResolveResult(t *testing.T) {

// 	// define static variables
// 	PKS := simapp.CreateTestPubKeys(5)

// 	valConsPk1 := PKS[0]
// 	valConsPk2 := PKS[1]
// 	valConsPk3 := PKS[2]

// 	// init sim app
// 	app := simapp.Setup(false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	addrs := simapp.AddTestAddrs(app, ctx, 10, sdk.NewInt(10000000000))
// 	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
// 	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

// 	// create validator with 10% commission and 300,000 orai
// 	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
// 	tstaking.CreateValidator(valAddrs[0], valConsPk1, sdk.NewInt(300000000), true)

// 	// create second validator with 10% commission and 250,000 orai
// 	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
// 	tstaking.CreateValidator(valAddrs[1], valConsPk2, sdk.NewInt(250000000), true)

// 	// create second validator with 10% commission and 150,000 orai
// 	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
// 	tstaking.CreateValidator(valAddrs[2], valConsPk3, sdk.NewInt(150000000), true)

// 	aioracleKeeper := aioraclekeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, nil)
// 	aioracleKeeper := aioraclekeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.StakingKeeper)

// 	// init keeper to run custom allocate tokens
// 	// here we borrow staking module to store the reward in the replacement of aioracle
// 	k := keeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.DistrKeeper, app.AccountKeeper, aioracleKeeper, aioracleKeeper, types.FeeCollectorName)

// 	// wrap keeper in a test keeper for test functions
// 	testKeeper := keeper.NewTestKeeper(*k, app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.DistrKeeper, app.AccountKeeper, aioracleKeeper, aioracleKeeper, types.FeeCollectorName)

// 	// init data source results
// 	dsResult1 := aioracletypes.NewDataSourceResult(&aioracletypes.EntryPoint{}, []byte{0x50}, "success")
// 	dsResult2 := aioracletypes.NewDataSourceResult(&aioracletypes.EntryPoint{}, []byte{0x50}, "success")
// 	dsResult3 := aioracletypes.NewDataSourceResult(&aioracletypes.EntryPoint{}, []byte{0x50}, "success")
// 	dsResults := []*aioracletypes.DataSourceResult{dsResult1, dsResult2, dsResult3}

// 	// init test case results
// 	tcResult1 := aioracletypes.NewTestCaseResult(&aioracletypes.EntryPoint{}, dsResults)
// 	tcResult2 := aioracletypes.NewTestCaseResult(&aioracletypes.EntryPoint{}, dsResults)
// 	tcResults := []*aioracletypes.TestCaseResult{tcResult1, tcResult2}

// 	// init reporter
// 	reporter := aioracletypes.NewReporter(addrs[0], "reporter", valAddrs[0])

// 	// init report
// 	id := ksuid.New().String()
// 	report := aioracletypes.NewReport(id, dsResults, tcResults, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, reporter, "")

// 	// verify report
// 	err := testKeeper.aioracleKeeper.AddReport(ctx, id, report)
// 	require.NoError(t, err)

// 	// test result with 2 validators aka two reports, total 70% reports to finish
// 	testKeeper.Keeper.ResolveResult(ctx, report, 100, 70)

// 	result, err := testKeeper.Keeper.GetResult(ctx, id)
// 	t.Logf("result after resolve with 2 validators: %v\n", result.Status)

// 	for i := 0; i < 68; i++ {
// 		// test result with 100 validators aka 100 reports, total 70% reports to finish
// 		testKeeper.Keeper.ResolveResult(ctx, report, 100, 70)
// 	}

// 	result, err = testKeeper.Keeper.GetResult(ctx, id)
// 	t.Logf("result after resolve with 100 validators and 68 reports: %v\n", result.Status)
// 	require.Equal(t, "pending", result.Status)

// 	// add another report to the list of 100 reports will lead to finished
// 	testKeeper.Keeper.ResolveResult(ctx, report, 100, 70)
// 	result, err = testKeeper.Keeper.GetResult(ctx, id)
// 	t.Logf("result after resolve with 100 validators and 69 reports: %v\n", result.Status)
// 	require.Equal(t, "finished", result.Status)

// 	id = ksuid.New().String()
// 	report = aioracletypes.NewReport(id, dsResults, tcResults, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, reporter, "")

// 	// verify report
// 	err = testKeeper.aioracleKeeper.AddReport(ctx, id, report)
// 	require.NoError(t, err)

// 	testKeeper.Keeper.ResolveResult(ctx, report, 1, 70)

// 	result, err = testKeeper.Keeper.GetResult(ctx, id)
// 	t.Logf("result after resolve with 1 validator: %v\n", result.Status)
// 	require.Equal(t, "finished", result.Status)
// }
