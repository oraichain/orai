package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/oraichain/orai/x/websocket/exported"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	// TODO: Register the modules msgs
	cdc.RegisterConcrete(MsgCreateReport{}, "websocket/AddReport", nil)
	cdc.RegisterConcrete(MsgAddReporter{}, "websocket/AddReporter", nil)
	cdc.RegisterConcrete(MsgRemoveReporter{}, "websocket/RemoveReporter", nil)	

	// cdc.RegisterConcrete(&Report{}, "websocket/Report", nil)
	// cdc.RegisterConcrete(&Reporter{}, "websocket/Reporter", nil)
	// cdc.RegisterConcrete(&DataSourceResult{}, "websocket/DataSourceResult", nil)
	// cdc.RegisterConcrete(&TestCaseResult{}, "websocket/TestCaseResult", nil)
	// cdc.RegisterConcrete(&Validator{}, "websocket/Validator", nil)
	// cdc.RegisterConcrete(&ValResult{}, "websocket/ValResult", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	// TODO: register msgs here to run
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateReport{},
		&MsgAddReporter{},
		&MsgRemoveReporter{},
	)

	// registry.RegisterImplementations(
	// 	(*govtypes.Content)(nil),
	// 	&Report{},
	// 	&Reporter{},
	// 	&DataSourceResult{},
	// 	&TestCaseResult{},
	// 	&Validator{},
	// 	&ValResult{},
	// )

	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
