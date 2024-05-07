package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	appparams "github.com/oraichain/orai/app/params"
	enccodec "github.com/tharsis/ethermint/encoding/codec"
)

// MakeEncodingConfig creates an EncodingConfig for gravity.
func MakeEncodingConfig() appparams.EncodingConfig {
	encodingConfig := appparams.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	enccodec.RegisterLegacyAminoCodec(encodingConfig.Amino)
	enccodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
