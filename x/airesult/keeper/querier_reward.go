package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airesult/types"
)

func queryReward(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "error")
	}
	// id of the request
	blockHeight := path[0]
	blockHeightInt, err := strconv.Atoi(blockHeight)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrBlockHeightInvalid, err.Error())
	}

	// get reward
	reward, err := k.GetReward(ctx, int64(blockHeightInt))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrRewardNotfound, err.Error())
	}

	res, err := codec.MarshalJSONIndent(k.cdc, types.NewQueryResReward(reward))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
