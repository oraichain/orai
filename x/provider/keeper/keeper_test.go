package keeper_test

import (
	"testing"

	"github.com/oraichain/orai/x/provider/keeper"
	"github.com/oraichain/orai/x/provider/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	providertypes "github.com/oraichain/orai/x/provider/types"
)

var (
	PKS = simapp.CreateTestPubKeys(5)
)

func TestGetMinimumFees(t *testing.T) {

	valCount := 10
	minFees := 294
	rewardPercentage := int64(60)

	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrs(app, ctx, 10, sdk.NewInt(10000000000))

	k := keeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName))

	// init data sources
	firstDataSource := providertypes.NewAIDataSource("first data source", "abc", addrs[0], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5))), "none")
	secondDataSource := providertypes.NewAIDataSource("2nd data source", "abc", addrs[1], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(5))), "none")
	thirdDataSource := providertypes.NewAIDataSource("3rd data source", "abc", addrs[2], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3))), "none")

	// init test cases
	firstTestCase := providertypes.NewTestCase("1st test case", "abc", addrs[3], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4))), "none")
	secondTestCase := providertypes.NewTestCase("2nd test case", "abc", addrs[4], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(4))), "none")

	k.SetAIDataSource(ctx, firstDataSource.GetName(), firstDataSource)
	k.SetAIDataSource(ctx, secondDataSource.GetName(), secondDataSource)
	k.SetAIDataSource(ctx, thirdDataSource.GetName(), thirdDataSource)
	k.SetTestCase(ctx, firstTestCase.GetName(), firstTestCase)
	k.SetTestCase(ctx, secondTestCase.GetName(), secondTestCase)

	coins, err := k.GetMinimumFees(ctx, []string{firstDataSource.GetName(), secondDataSource.GetName(), thirdDataSource.GetName()}, []string{firstTestCase.Name, secondTestCase.Name}, valCount, rewardPercentage)

	// provider fees must equal 29 ORAI
	require.NoError(t, err)
	require.Equal(t, sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(int64(minFees))}}, coins)

	// must error due to wrong name
	_, err = k.GetMinimumFees(ctx, []string{"firstDataSource.GetName()", secondDataSource.GetName(), thirdDataSource.GetName()}, []string{firstTestCase.Name, secondTestCase.Name}, 1, rewardPercentage)

	require.Error(t, err, sdkerrors.Wrap(types.ErrDataSourceNotFound, ""))

	// must error due to wrong name
	_, err = k.GetMinimumFees(ctx, []string{"firstDataSource.GetName()", secondDataSource.GetName(), thirdDataSource.GetName()}, []string{firstTestCase.Name, "ggg"}, 1, rewardPercentage)

	require.Error(t, err, sdkerrors.Wrap(types.ErrTestCaseNotFound, ""))
}
