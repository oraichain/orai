package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
var (
	ErrBlockHeightInvalid = sdkerrors.Register(ModuleName, 1, "The block height is invalid")
	ErrRewardNotfound     = sdkerrors.Register(ModuleName, 2, "The reward of the given block height is not found")
)
