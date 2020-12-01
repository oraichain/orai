package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airesult/types"
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
	reports := k.webSocketKeeper.GetReportsBlockHeight(ctx, blockHeight)
	// if there's no report from any validators, we skip
	if len(reports) == 0 {
		return
	}
	reward := types.DefaultReward(blockHeight)

	// Collect all the reports in the current block to get all the information for the reward
	for _, report := range k.webSocketKeeper.GetReportsBlockHeight(ctx, blockHeight) {
		k.ResolveRequestsFromReports(ctx, report, &reward, blockHeight)
	}
	// check if the reward is empty or not
	if len(reward.Validators) > 0 {
		fmt.Println("block for reward: ", blockHeight)
		k.SetReward(ctx, blockHeight, reward)
	}
}
