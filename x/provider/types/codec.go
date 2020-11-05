package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// TODO: Register the modules msgs
	cdc.RegisterConcrete(MsgCreateOracleScript{}, "provider/CreateOracleScript", nil)
	cdc.RegisterConcrete(MsgEditOracleScript{}, "provider/EditOracleScript", nil)
	cdc.RegisterConcrete(MsgCreateAIDataSource{}, "provider/CreateAIDataSource", nil)
	cdc.RegisterConcrete(MsgEditAIDataSource{}, "provider/EditAIDataSource", nil)
	// cdc.RegisterConcrete(MsgSetKYCRequest{}, "provider/SetKYCRequest", nil)
	// cdc.RegisterConcrete(MsgSetPriceRequest{}, "provider/SetPriceRequest", nil)
	cdc.RegisterConcrete(MsgCreateTestCase{}, "provider/SetTestCase", nil)
	cdc.RegisterConcrete(MsgEditTestCase{}, "provider/EditTestCase", nil)
	// cdc.RegisterConcrete(MsgCreateReport{}, "provider/AddReport", nil)
	// cdc.RegisterConcrete(MsgAddReporter{}, "provider/AddReporter", nil)
	// cdc.RegisterConcrete(MsgRemoveReporter{}, "provider/RemoveReporter", nil)
	// cdc.RegisterConcrete(MsgCreateStrategy{}, "provider/CreateStrategy", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
