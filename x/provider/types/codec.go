package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	// TODO: register msgs here to run
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateAIDataSource{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateOracleScript{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateTestCase{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgEditAIDataSource{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgEditOracleScript{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgEditTestCase{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
