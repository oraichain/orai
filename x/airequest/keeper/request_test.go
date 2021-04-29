package keeper_test

import (
	"fmt"
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
	testKeeper := aiRequestkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())

	// init ai request
	id := ksuid.New().String()
	airequest := types.NewAIRequest(id, sdk.AccAddress{}, addrs[0], []sdk.ValAddress{valAddrs[0], valAddrs[1], valAddrs[2]}, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, false)

	testKeeper.SetAIRequest(ctx, airequest.RequestID, airequest)
	//aiRequestGet, err := testKeeper.GetAIRequest(ctx, id)

	// verify if set AI request is legit or not
	// require.NoError(t, err)
	// require.Equal(t, aiRequestGet, airequest)

	// init data source results
	dsResult1 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult2 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult3 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{0x50}, types.ResultSuccess)
	dsResults := []*aiRequesttypes.Result{dsResult1, dsResult2, dsResult3}

	// init report
	report := types.NewReport(id, dsResults, 1, []byte{0x50}, valAddrs[0], types.ResultSuccess, nil)

	// verify report
	err := testKeeper.SetReport(ctx, id, report)
	require.NoError(t, err)

	reward := aiRequesttypes.DefaultReward(1)
	testKeeper.ResolveRequestsFromReports(ctx, report, reward)

	// Remember that we are using coins, so the values will be truncated along the way. When applying the equations in the medium post, remember to truncate the final result.
	t.Logf("provider fees: %v\n", reward.BaseReward.ProviderFees)
	t.Logf("validator fees: %v\n", reward.BaseReward.ValidatorFees)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(13))), reward.BaseReward.ProviderFees)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(5))), reward.BaseReward.ValidatorFees)
}

func TestResolveTestCaseRequestsFromReports(t *testing.T) {

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
	testKeeper := aiRequestkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())

	// init ai request
	id := ksuid.New().String()
	airequest := types.NewAIRequest(id, sdk.AccAddress{}, addrs[0], []sdk.ValAddress{valAddrs[0], valAddrs[1], valAddrs[2]}, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, false)

	testKeeper.SetAIRequest(ctx, airequest.RequestID, airequest)
	//aiRequestGet, err := testKeeper.GetAIRequest(ctx, id)

	// verify if set AI request is legit or not
	// require.NoError(t, err)
	// require.Equal(t, aiRequestGet, airequest)

	// init data source results
	tcResult1 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	tcResult2 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	tcResult3 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{0x50}, types.ResultSuccess)
	tcResults := []*aiRequesttypes.Result{tcResult1, tcResult2, tcResult3}
	dsResult := aiRequesttypes.NewResultWithTestCase(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, tcResults, types.ResultSuccess)
	dsResults := []*aiRequesttypes.ResultWithTestCase{dsResult, dsResult}

	// init report
	report := types.NewTestCaseReport(id, dsResults, 1, valAddrs[0], nil)

	// verify report
	err := testKeeper.SetTestCaseReport(ctx, id, report)
	require.NoError(t, err)

	reward := aiRequesttypes.DefaultReward(1)
	testKeeper.ResolveRequestsFromTestCaseReports(ctx, report, reward)

	// Remember that we are using coins, so the values will be truncated along the way. When applying the equations in the medium post, remember to truncate the final result.
	t.Logf("provider fees: %v\n", reward.BaseReward.ProviderFees)
	t.Logf("validator fees: %v\n", reward.BaseReward.ValidatorFees)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(26))), reward.BaseReward.ProviderFees)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(10))), reward.BaseReward.ValidatorFees)
}

func TestResolveBothRequestsFromReports(t *testing.T) {

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
	testKeeper := aiRequestkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())

	// init ai request
	id := ksuid.New().String()
	airequest := types.NewAIRequest(id, sdk.AccAddress{}, addrs[0], []sdk.ValAddress{valAddrs[0], valAddrs[1], valAddrs[2]}, 1, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(29))), []byte{0x50}, false)

	testKeeper.SetAIRequest(ctx, airequest.RequestID, airequest)
	//aiRequestGet, err := testKeeper.GetAIRequest(ctx, id)

	// verify if set AI request is legit or not
	// require.NoError(t, err)
	// require.Equal(t, aiRequestGet, airequest)

	// init data source results
	tcResult1 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	tcResult2 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	tcResult3 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{0x50}, types.ResultSuccess)
	tcResults := []*aiRequesttypes.Result{tcResult1, tcResult2, tcResult3}
	dsResult := aiRequesttypes.NewResultWithTestCase(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, tcResults, types.ResultSuccess)
	tcResultFinal := []*aiRequesttypes.ResultWithTestCase{dsResult, dsResult}

	// init data source results
	dsResult1 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult2 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5)))}, []byte{0x50}, types.ResultSuccess)
	dsResult3 := aiRequesttypes.NewResult(&aiRequesttypes.EntryPoint{ProviderFees: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3)))}, []byte{0x50}, types.ResultSuccess)
	dsResults := []*aiRequesttypes.Result{dsResult1, dsResult2, dsResult3}

	require.Equal(t, 3, len(dsResults))
	require.Equal(t, 2, len(tcResultFinal))

	count := 0
	for _, _ = range tcResultFinal {
		for _, _ = range tcResults {
			count++
		}
	}
	require.Equal(t, 6, count)

	// init report
	report := types.NewReport(id, dsResults, 1, []byte{0x50}, valAddrs[0], types.ResultFailure, nil)

	// verify report
	err := testKeeper.SetReport(ctx, id, report)
	require.NoError(t, err)

	// init report
	tcReport := types.NewTestCaseReport(id, tcResultFinal, 1, valAddrs[0], nil)

	// verify report
	err = testKeeper.SetTestCaseReport(ctx, id, tcReport)
	require.NoError(t, err)

	reward := aiRequesttypes.DefaultReward(1)
	testKeeper.ResolveRequestsFromReports(ctx, report, reward)
	testKeeper.ResolveRequestsFromTestCaseReports(ctx, tcReport, reward)

	fmt.Println("reward: ", len(reward.Results))
	require.Equal(t, 9, len(reward.Results))

	var totalFees sdk.Coins
	for _, result := range reward.Results {
		t.Logf("result fees: %v", result.EntryPoint.ProviderFees)
		totalFees = totalFees.Add(result.EntryPoint.ProviderFees...)
	}

	// Remember that we are using coins, so the values will be truncated along the way. When applying the equations in the medium post, remember to truncate the final result.
	t.Logf("provider fees: %v\n", reward.BaseReward.ProviderFees)
	t.Logf("validator fees: %v\n", reward.BaseReward.ValidatorFees)
	require.Equal(t, totalFees, reward.BaseReward.ProviderFees)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(15))), reward.BaseReward.ValidatorFees)
}
