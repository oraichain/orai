package keeper_test

// import (
// 	"testing"

// 	aioraclekeeper "github.com/oraichain/orai/x/aioracle/keeper"
// 	aioracle "github.com/oraichain/orai/x/aioracle/types"
// 	"github.com/oraichain/orai/x/aioracle/keeper"
// 	aioracletypes "github.com/oraichain/orai/x/aioracle/types"
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

// func TestResolveRequestsFromReports(t *testing.T) {

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

// 	// init ai request
// 	id := ksuid.New().String()
// 	aioracle := aioracle.Newaioracle(id, sdk.AccAddress{}, addrs[0], []sdk.ValAddress{valAddrs[0], valAddrs[1], valAddrs[2]}, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, []byte{0x49})

// 	testKeeper.aioracleKeeper.Setaioracle(ctx, aioracle.RequestID, aioracle)
// 	aioracleGet, err := testKeeper.aioracleKeeper.Getaioracle(ctx, id)

// 	// verify if set AI request is legit or not
// 	require.NoError(t, err)
// 	require.Equal(t, aioracleGet, aioracle)

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
// 	report := aioracletypes.NewReport(id, dsResults, tcResults, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, reporter, "")

// 	// verify report
// 	err = testKeeper.aioracleKeeper.AddReport(ctx, id, report)
// 	require.NoError(t, err)

// 	reward := aioracletypes.DefaultReward(1)
// 	testKeeper.Keeper.ResolveRequestsFromReports(ctx, report, reward, 1, 60)

// 	// Remember that we are using coins, so the values will be truncated along the way. When applying the equations in the medium post, remember to truncate the final result.
// 	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(21))), reward.ProviderFees)
// 	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(8))), reward.ValidatorFees)
// }
