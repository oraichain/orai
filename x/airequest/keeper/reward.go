package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airequest/types"
)

// SetReward saves the reward to the storage without performing validation.
func (k *Keeper) SetReward(ctx sdk.Context, rew *types.Reward) error {
	bz, err := k.Cdc.MarshalBinaryBare(rew)
	ctx.KVStore(k.StoreKey).Set(types.RewardStoreKey(rew.BaseReward.BlockHeight), bz)
	return err
}

// GetReward retrieves a specific reward given a block height
func (k *Keeper) GetReward(ctx sdk.Context, blockHeight int64) (*types.Reward, error) {
	store := ctx.KVStore(k.StoreKey)
	// check if there exists a reward in that block height or not
	hasReward := store.Has(types.RewardStoreKey(blockHeight))
	var err error
	if !hasReward {
		err = fmt.Errorf("")
		return nil, err
	}
	var reward types.Reward
	err = k.Cdc.UnmarshalBinaryBare(store.Get(types.RewardStoreKey(blockHeight)), &reward)
	if err != nil {
		return &types.Reward{
			BaseReward: &types.BaseReward{
				BlockHeight: int64(-1),
			},
			Results: []*types.Result{},
		}, err
	}
	return &reward, nil
}

// ProcessReward collects all the information needed to create a new Reward object
func (k *Keeper) ProcessReward(ctx sdk.Context) {
	reports := k.GetReportsBlockHeight(ctx)
	reward := types.DefaultReward(ctx.BlockHeight())

	// get param reward percentage oracle
	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range reports {
		isValid, valCount := k.ResolveRequestsFromReports(ctx, &report, reward)
		// if we can resolve the requests from the reports successfully, we resolve its result
		if isValid {
			// collect param for the the total reports needed to be considered finished
			totalReportsPercentage := k.GetParam(ctx, types.KeyReportPercentages)
			if totalReportsPercentage < uint64(0) || totalReportsPercentage > uint64(1) {
				totalReportsPercentage = types.DefaultReportPercentages
			}
			k.ResolveResult(ctx, &report, valCount, totalReportsPercentage)
		}
	}

	tcReports := k.GetTestCaseReportsBlockHeight(ctx)
	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range tcReports {
		k.ResolveRequestsFromTestCaseReports(ctx, &report, reward)
	}
	// check if the reward is empty or not
	if len(reward.BaseReward.Validators) > 0 {
		k.Logger(ctx).Info(fmt.Sprintf("reward information: %v", reward))
		k.SetReward(ctx, reward)
	}
}
