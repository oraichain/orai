package keeper_test

import (
	"testing"

	aioraclekeeper "github.com/oraichain/orai/x/aioracle/keeper"
	"github.com/oraichain/orai/x/aioracle/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestGetReward(t *testing.T) {

	// init sim app
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	testKeeper := aioraclekeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.AccountKeeper, app.DistrKeeper, feeCollector.GetName())
	test := types.Reward{
		Results: []*types.Result{{
			Status: "success",
		}},
	}

	rewardItem, _ := testKeeper.Cdc.MarshalBinaryBare(&test)
	t.Log("result", rewardItem)

	reward := &types.Reward{}
	// // var rewardItem []byte = nil
	err := testKeeper.Cdc.UnmarshalBinaryBare(rewardItem, reward)

	t.Log("result", err, reward)

}
