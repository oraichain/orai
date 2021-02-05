package keeper_test

import (
	"testing"

	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/oraichain/orai/x/provider"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	providertypes "github.com/oraichain/orai/x/provider/types"
)

var (
	PKS = simapp.CreateTestPubKeys(5)

	valConsPk1 = PKS[0]
	valConsPk2 = PKS[1]
	valConsPk3 = PKS[2]

	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())
)

func TestCollectRequestFees(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrs(app, ctx, 10, sdk.NewInt(10000000000))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 10% commission and 300,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[0], valConsPk1, sdk.NewInt(30000000), true)

	// create second validator with 10% commission and 250,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[1], valConsPk2, sdk.NewInt(25000000), true)

	// create second validator with 10% commission and 150,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[2], valConsPk3, sdk.NewInt(15000000), true)

	// init keeper to run custom allocate tokens
	// here we borrow staking module to store the reward in the replacement of airesult
	k := keeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, nil)

	// init data sources
	firstDataSource := providertypes.NewAIDataSource("first data source", "abc", addrs[0], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5))), "none")
	secondDataSource := providertypes.NewAIDataSource("2nd data source", "abc", addrs[1], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5))), "none")
	thirdDataSource := providertypes.NewAIDataSource("3rd data source", "abc", addrs[2], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3))), "none")

	// init test cases
	firstTestCase := providertypes.NewTestCase("1st test case", "abc", addrs[3], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4))), "none")
	secondTestCase := providertypes.NewTestCase("2nd test case", "abc", addrs[4], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4))), "none")

	// init oscript
	oscript := providertypes.NewOracleScript("oscript", "abc", addrs[0], "new oracle script", sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(29))), []string{firstDataSource.Name, secondDataSource.Name, thirdDataSource.Name}, []string{firstDataSource.Name, secondTestCase.Name})

	aiRequest := types.NewAIRequest(ksuid.New().String(), oscript.Name, addrs[0], []sdk.ValAddress{valAddrs[0], valAddrs[1], valAddrs[2]}, 1, []provider.AIDataSource{*firstDataSource, *secondDataSource, *thirdDataSource}, []provider.TestCase{*firstTestCase, *secondTestCase}, sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(29))), []byte{0x50}, []byte{0x49})

	k.SetAIRequest(ctx, aiRequest.RequestID, aiRequest)

	coins := k.CollectRequestFees(ctx, 1)

	require.Equal(t, sdk.NewCoins(sdk.NewCoin("orai", sdk.NewInt(29))), coins)
}
