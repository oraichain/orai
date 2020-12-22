package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	storeName = "airequest"
	restName  = storeName
)

// RegisterRoutes registers provider-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}



// RegisterRoutes registers provider-related REST handlers to a router
func RegisterWebsocketRoutes(cliCtx context.CLIContext, r *mux.Router, in <-chan bool, out chan<- bool) {
	registerWebsocketQueryRoutes(cliCtx, r, in, out)	
}