package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/x/airequest/exported"
)

// AIRequestSet is a set of functions exported for other modules to use through its keeper
type AIRequestSet interface {
	GetAIRequest(ctx sdk.Context, id string) (exported.AIRequestI, error)
}
