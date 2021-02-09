package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateOracleScript{}, "provider/CreateOracleScript", nil)
	cdc.RegisterConcrete(MsgEditOracleScript{}, "provider/EditOracleScript", nil)
	cdc.RegisterConcrete(MsgCreateAIDataSource{}, "provider/CreateAIDataSource", nil)
	cdc.RegisterConcrete(MsgEditAIDataSource{}, "provider/EditAIDataSource", nil)
	cdc.RegisterConcrete(MsgCreateTestCase{}, "provider/SetTestCase", nil)
	cdc.RegisterConcrete(MsgEditTestCase{}, "provider/EditTestCase", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	// TODO: register msgs here to run
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateAIDataSource{},
		&MsgCreateOracleScript{},
		&MsgCreateTestCase{},
		&MsgEditAIDataSource{},
		&MsgEditOracleScript{},
		&MsgEditTestCase{},
	)

	// registry.RegisterImplementations(
	// 	(*govtypes.Content)(nil),
	// 	&AIDataSource{},
	// 	&TestCase{},
	// )

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
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
