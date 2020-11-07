package provider

import (
	"github.com/oraichain/orai/x/ai-request/keeper"
	"github.com/oraichain/orai/x/ai-request/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	IPFSUrl           = types.IPFSUrl
	IPFSAdd           = types.IPFSAdd
	IPFSCat           = types.IPFSCat
	DefaultParamspace = types.DefaultParamspace
)

var (
	NewKeeper          = keeper.NewKeeper
	NewQuerier         = keeper.NewQuerier
	NewMsgSetAIRequest = types.NewMsgSetAIRequest
	ModuleCdc          = types.ModuleCdc
	RegisterCodec      = types.RegisterCodec
	NewGenesisState    = types.NewGenesisState
	RequestKeyPrefix   = types.RequestKeyPrefix
	ResultKeyPrefix    = types.ResultKeyPrefix
	ReportKeyPrefix    = types.ReportKeyPrefix
	ReporterKeyPrefix  = types.ReporterKeyPrefix
	RewardKeyPrefix    = types.RewardKeyPrefix
	StrategyKeyPrefix  = types.StrategyKeyPrefix
)

type (
	Keeper          = keeper.Keeper
	MsgSetAIRequest = types.MsgSetAIRequest
	AIRequest       = types.AIRequest
	GenesisState    = types.GenesisState
)
