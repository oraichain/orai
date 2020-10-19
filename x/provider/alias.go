package provider

import (
	"github.com/ducphamle2/dexai/x/provider/keeper"
	"github.com/ducphamle2/dexai/x/provider/types"
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
	NewKeeper                = keeper.NewKeeper
	NewQuerier               = keeper.NewQuerier
	NewMsgCreateOracleScript = types.NewMsgCreateOracleScript
	NewMsgCreateAIDataSource = types.NewMsgCreateAIDataSource
	NewMsgSetAIRequest       = types.NewMsgSetAIRequest
	ModuleCdc                = types.ModuleCdc
	RegisterCodec            = types.RegisterCodec
	NewGenesisState          = types.NewGenesisState
	OScriptKeyPrefix         = types.OScriptKeyPrefix
	DataSourceKeyPrefix      = types.DataSourceKeyPrefix
	RequestKeyPrefix         = types.RequestKeyPrefix
	ResultKeyPrefix          = types.ResultKeyPrefix
	TestCaseKeyPrefix        = types.TestCaseKeyPrefix
	ReportKeyPrefix          = types.ReportKeyPrefix
	ReporterKeyPrefix        = types.ReporterKeyPrefix
	RewardKeyPrefix          = types.RewardKeyPrefix
	StrategyKeyPrefix        = types.StrategyKeyPrefix
)

type (
	Keeper                = keeper.Keeper
	MsgCreateOracleScript = types.MsgCreateOracleScript
	MsgCreateAIDataSource = types.MsgCreateAIDataSource
	MsgSetAIRequest       = types.MsgSetAIRequest
	QueryResOracleScript  = types.QueryResOracleScript
	QueryResAIDataSource  = types.QueryResAIDataSource
	OracleScript          = types.OracleScript
	AIDataSource          = types.AIDataSource
	AIRequest             = types.AIRequest
	TestCase              = types.TestCase
	GenesisState          = types.GenesisState
)
