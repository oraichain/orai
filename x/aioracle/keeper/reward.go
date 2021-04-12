package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/aioracle/types"
)

// SetReward saves the reward to the storage without performing validation.
func (k Keeper) SetReward(ctx sdk.Context, rew *types.Reward) error {
	bz, err := k.Cdc.MarshalBinaryBare(rew)
	ctx.KVStore(k.StoreKey).Set(types.RewardStoreKey(rew.BlockHeight), bz)
	return err
}

// GetReward retrieves a specific reward given a block height
func (k Keeper) GetReward(ctx sdk.Context, blockHeight int64) (*types.Reward, error) {
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
			BlockHeight: int64(-1),
		}, err
	}
	return &reward, nil
}

// ProcessReward collects all the information needed to create a new Reward object
func (k Keeper) ProcessReward(ctx sdk.Context) {
	reports := k.GetReportsBlockHeight(ctx)
	// if there's no report from any validators, we skip
	if len(reports) == 0 {
		return
	}
	reward := types.DefaultReward(ctx.BlockHeight())

	// get param reward percentage oracle
	rewardPercentage := int64(70)
	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range reports {
		isValid, valCount := k.ResolveRequestsFromReports(ctx, &report, reward, rewardPercentage)
		// if we can resolve the requests from the reports successfully, we resolve its result
		if isValid {
			// collect param for the the total reports needed to be considered finished
			totalReportsPercentage := k.GetParam(ctx, types.KeyReportPercentages)
			if totalReportsPercentage <= uint64(0) {
				totalReportsPercentage = uint64(70)
			}
			k.ResolveResult(ctx, &report, valCount, totalReportsPercentage)
		}

	}
	// check if the reward is empty or not
	if len(reward.Validators) > 0 {
		k.Logger(ctx).Info(fmt.Sprintf("block for reward: %v\n", ctx.BlockHeight()))
		k.SetReward(ctx, reward)
	}
}
