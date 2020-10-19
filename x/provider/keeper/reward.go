package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ducphamle2/dexai/x/provider/types"
)

// SetReward saves the reward to the storage without performing validation.
func (k Keeper) SetReward(ctx sdk.Context, blockHeight int64, rew types.Reward) {
	ctx.KVStore(k.storeKey).Set(types.RewardStoreKey(blockHeight), k.cdc.MustMarshalBinaryLengthPrefixed(rew))
}

// GetReward retrieves a specific reward given a block height
func (k Keeper) GetReward(ctx sdk.Context, blockHeight int64) (types.Reward, error) {
	store := ctx.KVStore(k.storeKey)
	var reward types.Reward
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.RewardStoreKey(blockHeight)), &reward)
	if err != nil {
		return types.Reward{
			BlockHeight: int64(-1),
		}, err
	}
	return reward, nil
}

// ProcessReward collects all the information needed to create a new Reward object
func (k Keeper) ProcessReward(ctx sdk.Context) {
	blockHeight := ctx.BlockHeight()

	reports := k.GetReportsBlockHeight(ctx, blockHeight)

	// if there's no report from any validators, we skip
	if len(reports) == 0 {
		return
	}

	reward := types.NewReward([]types.Validator{}, []sdk.AccAddress{}, []sdk.AccAddress{}, blockHeight, int64(0))

	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range k.GetReportsBlockHeight(ctx, blockHeight) {
		vals, epOwners, tcOwners, votingPower := k.ResolveRequest(ctx, report, ctx.BlockHeight())
		reward.Validators = append(reward.Validators, vals...)
		reward.DataSourceOwners = append(reward.DataSourceOwners, epOwners...)
		reward.TestCaseOwners = append(reward.TestCaseOwners, tcOwners...)
		reward.TotalPower = reward.TotalPower + votingPower
	}

	k.SetReward(ctx, blockHeight, reward)
}
