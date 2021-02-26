package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airesult/types"
)

// SetReward saves the reward to the storage without performing validation.
func (k Keeper) SetReward(ctx sdk.Context, rew *types.Reward) error {
	bz, err := k.cdc.MarshalBinaryBare(rew)
	ctx.KVStore(k.storeKey).Set(types.RewardStoreKey(rew.BlockHeight), bz)
	return err
}

// GetReward retrieves a specific reward given a block height
func (k Keeper) GetReward(ctx sdk.Context, blockHeight int64) (*types.Reward, error) {
	store := ctx.KVStore(k.storeKey)
	// check if there exists a reward in that block height or not
	hasReward := store.Has(types.RewardStoreKey(blockHeight))
	var err error
	if !hasReward {
		err = fmt.Errorf("")
		return nil, err
	}
	var reward types.Reward
	err = k.cdc.UnmarshalBinaryBare(store.Get(types.RewardStoreKey(blockHeight)), &reward)
	if err != nil {
		return &types.Reward{
			BlockHeight: int64(-1),
		}, err
	}
	return &reward, nil
}

// ProcessReward collects all the information needed to create a new Reward object
func (k Keeper) ProcessReward(ctx sdk.Context) {
	blockHeight := ctx.BlockHeight()
	reports := k.webSocketKeeper.GetReportsBlockHeight(ctx, blockHeight)
	// if there's no report from any validators, we skip
	if len(reports) == 0 {
		return
	}
	reward := types.DefaultReward(blockHeight)

	// get param reward percentage oracle
	rewardPercentage := k.providerKeeper.GetOracleScriptRewardPercentageParam(ctx)
	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range k.webSocketKeeper.GetReportsBlockHeight(ctx, blockHeight) {
		isValid, valCount := k.ResolveRequestsFromReports(ctx, &report, reward, blockHeight, rewardPercentage)
		// if we can resolve the requests from the reports successfully, we resolve its result
		if isValid {
			// collect param for the the total reports needed to be considered finished
			totalReportsPercentage := k.GetTotalReportsParam(ctx)
			if k.GetTotalReportsParam(ctx) <= int64(0) {
				totalReportsPercentage = int64(70)
			}
			k.ResolveResult(ctx, &report, valCount, totalReportsPercentage)
		}

	}
	// check if the reward is empty or not
	if len(reward.Validators) > 0 {
		k.Logger(ctx).Info(fmt.Sprintf("block for reward: %v\n", blockHeight))
		k.SetReward(ctx, reward)
	}
}
