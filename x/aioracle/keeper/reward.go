package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// SetReward saves the reward to the storage without performing validation.
func (k *Keeper) SetReward(ctx sdk.Context, rew *types.Reward) error {
	bz, err := k.Cdc.MarshalBinaryBare(rew)
	if err == nil {
		ctx.KVStore(k.StoreKey).Set(types.RewardStoreKey(rew.BaseReward.BlockHeight), bz)
	}
	return err
}

// GetReward retrieves a specific reward given a block height
func (k *Keeper) GetReward(ctx sdk.Context, blockHeight int64) (*types.Reward, error) {
	store := ctx.KVStore(k.StoreKey)
	rewardKey := types.RewardStoreKey(blockHeight)
	// check if there exists a reward in that block height or not
	rewardItem := store.Get(rewardKey)
	if rewardItem == nil {
		return nil, fmt.Errorf("no reward item at key: %s", rewardKey)
	}
	reward := &types.Reward{}
	err := k.Cdc.UnmarshalBinaryBare(rewardItem, reward)
	return reward, err
}

// ProcessReward collects all the information needed to create a new Reward object
func (k *Querier) ProcessReward(ctx sdk.Context) {
	reports := k.keeper.GetReportsBlockHeight(ctx)
	reward := types.DefaultReward(ctx.BlockHeight())

	// get param reward percentage oracle
	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range reports {
		isValid, valCount := k.ResolveRequestsFromReports(ctx, &report, reward)
		// if we can resolve the requests from the reports successfully, we resolve its result
		if isValid {
			// collect param for the the total reports needed to be considered finished
			totalReportsPercentage := k.keeper.GetParam(ctx, types.KeyReportPercentages)
			if totalReportsPercentage < uint64(0) || totalReportsPercentage > uint64(1) {
				totalReportsPercentage = types.DefaultReportPercentages
			}
			k.keeper.ResolveResult(ctx, &report, valCount, totalReportsPercentage)
		}
	}

	tcReports := k.keeper.GetTestCaseReportsBlockHeight(ctx)
	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range tcReports {
		k.ResolveRequestsFromTestCaseReports(ctx, &report, reward)
	}
	// check if the reward is empty or not
	if len(reward.BaseReward.Validators) > 0 {
		k.keeper.Logger(ctx).Info(fmt.Sprintf("reward information: %v", reward))
		k.keeper.SetReward(ctx, reward)
	}
}
