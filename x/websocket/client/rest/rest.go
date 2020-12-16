package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	storeName = "websocket"
	restName  = storeName
)

// RegisterRoutes registers provider-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, in <-chan bool, out chan<- bool) {
	registerQueryRoutes(cliCtx, r, in, out)
	registerTxRoutes(cliCtx, r)
}
