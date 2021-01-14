package keeper

import (
	"github.com/CosmWasm/wasmd/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
