package airesult

import (
	"github.com/oraichain/orai/x/airesult/keeper"
	"github.com/oraichain/orai/x/airesult/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
)

var (
	NewKeeper       = keeper.NewKeeper
	NewQuerier      = keeper.NewQuerier
	ModuleCdc       = types.ModuleCdc
	RegisterCodec   = types.RegisterCodec
	NewGenesisState = types.NewGenesisState
	ResultKeyPrefix = types.ResultKeyPrefix
	RewardKeyPrefix = types.RewardKeyPrefix
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
