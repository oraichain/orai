package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/oraichain/orai/x/websocket/exported"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// TODO: Register the modules msgs
	cdc.RegisterConcrete(MsgCreateReport{}, "websocket/AddReport", nil)
	cdc.RegisterConcrete(MsgAddReporter{}, "websocket/AddReporter", nil)
	cdc.RegisterConcrete(MsgRemoveReporter{}, "websocket/RemoveReporter", nil)
	// When exporting interfaces for other modules to use, we need to register those interfaces as well as concrete structs so that the message can be encoded properly according to the Codec module
	cdc.RegisterInterface((*exported.ReportI)(nil), nil) // has to be pointer of interface
	cdc.RegisterInterface((*exported.ReporterI)(nil), nil)
	cdc.RegisterInterface((*exported.DataSourceResultI)(nil), nil) // has to be pointer of interface
	cdc.RegisterInterface((*exported.TestCaseResultI)(nil), nil)   // has to be pointer of interface
	cdc.RegisterInterface((*exported.ValidatorI)(nil), nil)        // has to be pointer of interface
	cdc.RegisterInterface((*exported.ValResultI)(nil), nil)        // has to be pointer of interface
	cdc.RegisterConcrete(&Report{}, "websocket/Report", nil)
	cdc.RegisterConcrete(&Reporter{}, "websocket/Reporter", nil)
	cdc.RegisterConcrete(&DataSourceResult{}, "websocket/DataSourceResult", nil)
	cdc.RegisterConcrete(&TestCaseResult{}, "websocket/TestCaseResult", nil)
	cdc.RegisterConcrete(&Validator{}, "websocket/Validator", nil)
	cdc.RegisterConcrete(&ValResult{}, "websocket/ValResult", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
