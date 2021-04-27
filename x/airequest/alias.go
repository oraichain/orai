package airequest

import (
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
)

const (
	ModuleName                     = types.ModuleName
	RouterKey                      = types.RouterKey
	StoreKey                       = types.StoreKey
	QuerierRoute                   = types.QuerierRoute
	DefaultParamspace              = types.DefaultParamspace
	EventTypeSetAIRequest          = types.EventTypeSetAIRequest
	AttributeRequestValidator      = types.AttributeRequestValidator
	EventTypeRequestWithData       = types.EventTypeRequestWithData
	AttributeRequestDSources       = types.AttributeRequestDSources
	AttributeRequestTCases         = types.AttributeRequestTCases
	AttributeRequestID             = types.AttributeRequestID
	AttributeContract              = types.AttributeContract
	AttributeRequestCreator        = types.AttributeRequestCreator
	AttributeRequestValidatorCount = types.AttributeRequestValidatorCount
	AttributeRequestInput          = types.AttributeRequestInput
	Denom                          = types.Denom
)

var (
	NewKeeper             = keeper.NewKeeper
	NewQuerier            = keeper.NewQuerier
	NewMsgSetAIRequestReq = types.NewMsgSetAIRequestReq
	ModuleCdc             = types.ModuleCdc
	RegisterCodec         = types.RegisterCodec
	NewGenesisState       = types.NewGenesisState
	RequestKeyPrefix      = types.RequestKeyPrefix
	ErrRequestNotFound    = types.ErrRequestNotFound
	NewAIRequest          = types.NewAIRequest
	ResultKeyPrefix       = types.ResultKeyPrefix
	RewardKeyPrefix       = types.RewardKeyPrefix
	ReportKeyPrefix       = types.ReportKeyPrefix
)

type (
	Keeper             = keeper.Keeper
	MsgSetAIRequestReq = types.MsgSetAIRequestReq
	AIRequest          = types.AIRequest
	GenesisState       = types.GenesisState
	Report             = types.Report
)
