package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	// TODO: Register the modules msgs
	// cdc.RegisterConcrete(MsgSetKYCRequest{}, "airequest/SetKYCRequest", nil)
	// cdc.RegisterConcrete(MsgSetPriceRequest{}, "airequest/SetPriceRequest", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
