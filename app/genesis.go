package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	mint "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONMarshaler) GenesisState {

	encCfg := MakeEncodingConfig()
	genesisState := ModuleBasics.DefaultGenesis(encCfg.Marshaler)

	// override
	mintGenesis := mint.DefaultGenesisState()
	mintGenesis.Params.BlocksPerYear = 6311200 // target 5-second block time
	mintGenesis.Params.MintDenom = "orai"

	genesisState[mint.ModuleName] = cdc.MustMarshalJSON(mintGenesis)

	return genesisState
}
