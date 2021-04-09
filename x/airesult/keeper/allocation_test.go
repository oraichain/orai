package keeper_test

import (
	"testing"

	airequestkeeper "github.com/oraichain/orai/x/airequest/keeper"
	airequest "github.com/oraichain/orai/x/airequest/types"
	"github.com/oraichain/orai/x/airesult/keeper"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	airesulttypes "github.com/oraichain/orai/x/airesult/types"

	websockettypes "github.com/oraichain/orai/x/websocket/types"
)

var (
	PKS = simapp.CreateTestPubKeys(5)

	valConsPk1 = PKS[0]
	valConsPk2 = PKS[1]
	valConsPk3 = PKS[2]

	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())

	authAcc  = authtypes.NewEmptyModuleAccount("auth")
	distrAcc = authtypes.NewEmptyModuleAccount("distribution")
	bankAcc  = authtypes.NewEmptyModuleAccount("bank")
)

func TestAllocateTokensToManyValidators(t *testing.T) {

	// define static variables

	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrs(app, ctx, 10, sdk.NewInt(10000000000))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 10% commission and 300,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[7], valConsPk1, sdk.NewInt(30000000), true)

	// create second validator with 10% commission and 250,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[8], valConsPk2, sdk.NewInt(25000000), true)

	// create second validator with 10% commission and 150,000 orai
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(10, 2), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[9], valConsPk3, sdk.NewInt(15000000), true)

	abciValA := abci.Validator{
		Address: valConsPk1.Address(),
		Power:   30,
	}
	abciValB := abci.Validator{
		Address: valConsPk2.Address(),
		Power:   25,
	}
	abciValC := abci.Validator{
		Address: valConsPk3.Address(),
		Power:   15,
	}

	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[7]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[8]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[9]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[7]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[8]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[9]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[7]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[8]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[9]).Rewards.IsZero())

	// allocate tokens as if both had voted and second was proposer
	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)))
	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, types.FeeCollectorName)
	require.NotNil(t, feeCollector)

	err := app.BankKeeper.SetBalances(ctx, feeCollector.GetAddress(), fees)

	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.GetModuleAccount(ctx, types.FeeCollectorName))

	votes := []abci.VoteInfo{
		{
			Validator:       abciValA,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValB,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValC,
			SignedLastBlock: true,
		},
	}

	airequestKeeper := airequestkeeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, nil)

	// init keeper to run custom allocate tokens
	// here we borrow staking module to store the reward in the replacement of airesult
	k := keeper.NewKeeper(app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.DistrKeeper, app.AccountKeeper, nil, airequestKeeper, types.FeeCollectorName)

	// wrap keeper in a test keeper for test functions
	testKeeper := keeper.NewTestKeeper(*k, app.AppCodec(), app.GetKey("staking"), nil, app.GetSubspace(stakingtypes.ModuleName), app.StakingKeeper, app.BankKeeper, app.DistrKeeper, app.AccountKeeper, nil, airequestKeeper, types.FeeCollectorName)

	id := ksuid.New().String()
	testKeeper.AiRequestKeeper.SetAIRequest(ctx, id, &airequest.AIRequest{RequestID: id})

	aiRequest, err := testKeeper.AiRequestKeeper.GetAIRequest(ctx, id)

	require.NoError(t, err)
	require.Equal(t, &airequest.AIRequest{RequestID: id}, aiRequest)

	// init reward
	reward := airesulttypes.DefaultReward(0)

	// set validators
	rewardRatio := sdk.NewDec(int64(1)).Sub(sdk.NewDecWithPrec(int64(60), 2))
	valFees, _ := sdk.NewDecCoinsFromCoins(reward.ProviderFees...).MulDec(rewardRatio).TruncateDecimal()
	validatorA := &websockettypes.Validator{valAddrs[7], abciValA.Power, "active"}
	validatorB := &websockettypes.Validator{valAddrs[8], abciValB.Power, "active"}
	validatorC := &websockettypes.Validator{valAddrs[9], abciValC.Power, "active"}
	reward.ValidatorFees = reward.ValidatorFees.Add(valFees...).Add(valFees...).Add(valFees...)
	reward.TotalPower = reward.TotalPower + validatorA.VotingPower + validatorB.VotingPower + validatorC.VotingPower
	reward.Validators = append(reward.Validators, *validatorA)
	reward.Validators = append(reward.Validators, *validatorB)
	reward.Validators = append(reward.Validators, *validatorC)

	temp := reward.ProviderFees

	// multiply by 3 because there are three validators
	reward.ProviderFees = reward.ProviderFees.Add(temp...).Add(temp...)

	// set reward
	testKeeper.Keeper.SetReward(ctx, reward)

	t.Logf("balance of provider 1: %v\n", app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))

	t.Logf("balance of provider 2: %v\n", app.BankKeeper.GetBalance(ctx, addrs[1], sdk.DefaultBondDenom))

	t.Logf("balance of provider 3: %v\n", app.BankKeeper.GetBalance(ctx, addrs[2], sdk.DefaultBondDenom))

	t.Logf("balance of provider 4: %v\n", app.BankKeeper.GetBalance(ctx, addrs[3], sdk.DefaultBondDenom))

	t.Logf("balance of provider 5: %v\n", app.BankKeeper.GetBalance(ctx, addrs[4], sdk.DefaultBondDenom))

	testKeeper.Keeper.AllocateTokens(ctx, votes, 1)
	rewardObj, err := testKeeper.Keeper.GetReward(ctx, 1-1)

	require.NoError(t, err)

	// confirm that each validator fee is 8.4 ORAI
	require.Equal(t, sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(8)}}, valFees)

	// provider fees must equal 21 * 3 ORAI (3 validators)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDecWithPrec(63, 0)}}, sdk.NewDecCoinsFromCoins(rewardObj.ProviderFees...))

	// validator fees must equal 24 ORAI
	require.Equal(t, sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(24)}}, rewardObj.ValidatorFees)

	// reward must equal to 87 ORAI
	require.Equal(t, sdk.Coins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(87)}}, rewardObj.ProviderFees.Add(rewardObj.ValidatorFees...))

	// total power must equal
	require.Equal(t, int64(rewardObj.TotalPower), validatorA.VotingPower+validatorB.VotingPower+validatorC.VotingPower)
	require.Equal(t, int64(rewardObj.TotalPower), int64(70))

	t.Logf("after allocation\n")

	t.Logf("outstanding reward of validators: %v\n", app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[8]))

	t.Logf("balance of provider 1: %v\n", app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))

	t.Logf("balance of provider 2: %v\n", app.BankKeeper.GetBalance(ctx, addrs[1], sdk.DefaultBondDenom))

	t.Logf("balance of provider 3: %v\n", app.BankKeeper.GetBalance(ctx, addrs[2], sdk.DefaultBondDenom))

	t.Logf("balance of provider 4: %v\n", app.BankKeeper.GetBalance(ctx, addrs[3], sdk.DefaultBondDenom))

	t.Logf("balance of provider 5: %v\n", app.BankKeeper.GetBalance(ctx, addrs[4], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000015))), app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000015))), app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000009))), app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000012))), app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))

	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000012))), app.BankKeeper.GetBalance(ctx, addrs[0], sdk.DefaultBondDenom))
}

func TestAllocateTokensTruncation(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(1234))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 10% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[0], valConsPk1, sdk.NewInt(110), true)

	// create second validator with 10% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[1], valConsPk2, sdk.NewInt(100), true)

	// create third validator with 10% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[2], valConsPk3, sdk.NewInt(100), true)

	abciValA := abci.Validator{
		Address: valConsPk1.Address(),
		Power:   11,
	}
	abciValB := abci.Validator{
		Address: valConsPk2.Address(),
		Power:   10,
	}
	abciValС := abci.Validator{
		Address: valConsPk3.Address(),
		Power:   10,
	}

	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[0]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[1]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[0]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[1]).Rewards.IsZero())

	// allocate tokens as if both had voted and second was proposer
	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(634195840)))

	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, types.FeeCollectorName)
	require.NotNil(t, feeCollector)

	err := app.BankKeeper.SetBalances(ctx, feeCollector.GetAddress(), fees)
	require.NoError(t, err)

	app.AccountKeeper.SetAccount(ctx, feeCollector)

	votes := []abci.VoteInfo{
		{
			Validator:       abciValA,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValB,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValС,
			SignedLastBlock: true,
		},
	}
	app.DistrKeeper.AllocateTokens(ctx, 31, 31, sdk.ConsAddress(valConsPk2.Address()), votes)

	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0]).Rewards.IsValid())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards.IsValid())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[2]).Rewards.IsValid())
}
