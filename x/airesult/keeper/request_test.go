package keeper_test

import (
	"fmt"
	"testing"

	airequestkeeper "github.com/oraichain/orai/x/airequest/keeper"
	airequest "github.com/oraichain/orai/x/airequest/types"
	"github.com/oraichain/orai/x/airesult/keeper"
	airesulttypes "github.com/oraichain/orai/x/airesult/types"
	"github.com/oraichain/orai/x/provider"
	providerkeeper "github.com/oraichain/orai/x/provider/keeper"
	websocketkeeper "github.com/oraichain/orai/x/websocket/keeper"
	websockettypes "github.com/oraichain/orai/x/websocket/types"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	providertypes "github.com/oraichain/orai/x/provider/types"
)

func TestResolveRequestsFromReports(t *testing.T) {

	// define static variables
	PKS := simapp.CreateTestPubKeys(5)

	valConsPk1 := PKS[0]
	valConsPk2 := PKS[1]
	valConsPk3 := PKS[2]
	valConsPk4 := PKS[4]

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

	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[3], valConsPk4, sdk.NewInt(300000000), true)

	providerKeeper := providerkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName))
	airequestKeeper := airequestkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, nil)
	websocketKeeper := websocketkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, nil, app.StakingKeeper, airequestKeeper)

	// init keeper to run custom allocate tokens
	// here we borrow staking module to store the reward in the replacement of airesult
	k := keeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, providerKeeper, app.BankKeeper, app.DistrKeeper, app.AccountKeeper, websocketKeeper, airequestKeeper, types.FeeCollectorName)

	// wrap keeper in a test keeper for test functions
	testKeeper := keeper.NewTestKeeper(*k, app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, providerKeeper, app.BankKeeper, app.DistrKeeper, app.AccountKeeper, websocketKeeper, airequestKeeper, types.FeeCollectorName)

	// init data sources
	firstDataSource := providertypes.NewAIDataSource("first data source", "abc", addrs[0], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5))), "none")
	secondDataSource := providertypes.NewAIDataSource("2nd data source", "abc", addrs[1], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5))), "none")
	thirdDataSource := providertypes.NewAIDataSource("3rd data source", "abc", addrs[2], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3))), "none")

	// init test cases
	firstTestCase := providertypes.NewTestCase("1st test case", "abc", addrs[3], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4))), "none")
	secondTestCase := providertypes.NewTestCase("2nd test case", "abc", addrs[4], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4))), "none")

	// store scripts into provider keeper. We must use the wrapper keeper since all the fields in the normal keeper are hidden by default, we cannot access them
	testKeeper.ProviderKeeper.SetAIDataSource(ctx, firstDataSource.GetName(), firstDataSource)
	testKeeper.ProviderKeeper.SetAIDataSource(ctx, secondDataSource.GetName(), secondDataSource)
	testKeeper.ProviderKeeper.SetAIDataSource(ctx, thirdDataSource.GetName(), thirdDataSource)
	testKeeper.ProviderKeeper.SetTestCase(ctx, firstTestCase.GetName(), firstTestCase)
	testKeeper.ProviderKeeper.SetTestCase(ctx, secondTestCase.GetName(), secondTestCase)

	// init oscript
	oscript := providertypes.NewOracleScript("oscript", "abc", addrs[0], "new oracle script", sdk.NewCoins(sdk.NewCoin(provider.Denom, sdk.NewInt(29))), []string{firstDataSource.Name, secondDataSource.Name, thirdDataSource.Name}, []string{firstDataSource.Name, secondTestCase.Name})

	// init ai request
	id := ksuid.New().String()
	aiRequest := airequest.NewAIRequest(id, oscript.Name, addrs[0], []sdk.ValAddress{valAddrs[0], valAddrs[1], valAddrs[2]}, 1, []provider.AIDataSource{*firstDataSource, *secondDataSource, *thirdDataSource}, []provider.TestCase{*firstTestCase, *secondTestCase}, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, []byte{0x49})

	testKeeper.AiRequestKeeper.SetAIRequest(ctx, aiRequest.RequestID, aiRequest)
	aiRequestGet, err := testKeeper.AiRequestKeeper.GetAIRequest(ctx, id)

	// verify if set AI request is legit or not
	require.NoError(t, err)
	require.Equal(t, aiRequestGet, aiRequest)

	// init data source results
	dsResult1 := websockettypes.NewDataSourceResult(firstDataSource.Name, []byte{0x50}, "success")
	dsResult2 := websockettypes.NewDataSourceResult(secondDataSource.Name, []byte{0x50}, "success")
	dsResult3 := websockettypes.NewDataSourceResult(thirdDataSource.Name, []byte{0x50}, "success")
	dsResults := []*websockettypes.DataSourceResult{dsResult1, dsResult2, dsResult3}

	// init test case results
	tcResult1 := websockettypes.NewTestCaseResult(firstTestCase.Name, dsResults)
	tcResult2 := websockettypes.NewTestCaseResult(firstTestCase.Name, dsResults)
	tcResults := []*websockettypes.TestCaseResult{tcResult1, tcResult2}

	// init reporter with validator 0
	reporter := websockettypes.NewReporter(addrs[0], "reporter", valAddrs[0])

	// init report
	report := websockettypes.NewReport(id, dsResults, tcResults, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, reporter, "")

	// verify report
	err = testKeeper.WebSocketKeeper.AddReport(ctx, id, report)
	require.NoError(t, err)

	reward := airesulttypes.DefaultReward(1)
	isValid, _ := testKeeper.Keeper.ResolveRequestsFromReports(ctx, report, reward, 1, 60)

	// if it's validator 0 then true
	require.Equal(t, true, isValid)
	fmt.Println("validator status: ", testKeeper.StakingKeeper.Validator(ctx, valAddrs[0]).GetStatus().String())
	fmt.Println("reward: ", reward.GetValidators())

	// init reporter with validator 0
	reporter = websockettypes.NewReporter(addrs[0], "reporter", valAddrs[3])

	// init report
	report = websockettypes.NewReport(id, dsResults, tcResults, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, reporter, "")

	// verify report
	err = testKeeper.WebSocketKeeper.AddReport(ctx, id, report)
	require.NoError(t, err)

	reward = airesulttypes.DefaultReward(1)
	isValid, _ = testKeeper.Keeper.ResolveRequestsFromReports(ctx, report, reward, 1, 60)

	// if it's validator 3 then should be false
	require.Equal(t, false, isValid)
}
