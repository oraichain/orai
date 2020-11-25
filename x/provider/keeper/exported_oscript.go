package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/provider/exported"
)

// This file implements all the exported functions of the Data Source struct that are exported by the Keeper

// GetOracleScriptI returns the oracle script object given the id of the oracle script
func (k Keeper) GetOracleScriptI(ctx sdk.Context, oScript string) (exported.OracleScriptI, error) {
	oracleScript, err := k.GetOracleScript(ctx, oScript)
	if err != nil {
		return nil, err
	}
	return oracleScript, nil
}
